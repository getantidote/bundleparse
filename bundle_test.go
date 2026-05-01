package bundleparse

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseBundleLine(t *testing.T) {
	bundle, err := ParseBundleLine(`foo/bar kind:zsh branch:main path:plugins/git`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := Bundle{
		Name:             "foo/bar",
		Kind:             "zsh",
		Branch:           "main",
		Path:             "plugins/git",
		ExtraAnnotations: map[string]string{},
	}

	if !reflect.DeepEqual(bundle, want) {
		t.Fatalf("mismatch:\n got: %#v\nwant: %#v", bundle, want)
	}
}

func TestParseBundleLine_AdditionalAnnotations(t *testing.T) {
	bundle, err := ParseBundleLine(`foo/bar kind:zsh pin:v1 branch:main conditional:if-true autoload:yes pre:"echo hi" post:'echo bye' fpath-rule:prepend unknown:yes`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := Bundle{
		Name:        "foo/bar",
		Kind:        "zsh",
		Branch:      "main",
		Pin:         "v1",
		Conditional: "if-true",
		Autoload:    "yes",
		Pre:         "echo hi",
		Post:        "echo bye",
		FpathRule:   "prepend",
		ExtraAnnotations: map[string]string{
			"unknown": "yes",
		},
	}

	if !reflect.DeepEqual(bundle, want) {
		t.Fatalf("mismatch:\n got: %#v\nwant: %#v", bundle, want)
	}
}

func TestParseBundles(t *testing.T) {
	input := strings.Join([]string{
		`# comment line`,
		`foo/bar kind:zsh branch:main path:plugins/git`,
		``,
		`ohmyzsh/ohmyzsh kind:zsh path:plugins/git`,
	}, "\n")

	bundles, err := ParseBundles(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []Bundle{
		{
			Name:             "foo/bar",
			Kind:             "zsh",
			Branch:           "main",
			Path:             "plugins/git",
			Line:             2,
			ExtraAnnotations: map[string]string{},
		},
		{
			Name:             "ohmyzsh/ohmyzsh",
			Kind:             "zsh",
			Path:             "plugins/git",
			Line:             4,
			ExtraAnnotations: map[string]string{},
		},
	}

	if !reflect.DeepEqual(bundles, want) {
		t.Fatalf("mismatch:\n got: %#v\nwant: %#v", bundles, want)
	}
}

func TestParseBundles_UsingDirectiveAppliesToBareBundle(t *testing.T) {
	input := strings.Join([]string{
		`using:ohmyzsh/ohmyzsh kind:zsh`,
		`foo`,
	}, "\n")

	bundles, err := ParseBundles(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []Bundle{
		{
			Name:             "ohmyzsh/ohmyzsh",
			Kind:             "zsh",
			Path:             "ohmyzsh/ohmyzsh/foo",
			Line:             2,
			ExtraAnnotations: map[string]string{},
		},
	}

	if !reflect.DeepEqual(bundles, want) {
		t.Fatalf("mismatch:\n got: %#v\nwant: %#v", bundles, want)
	}
}

func TestParseBundles_UsingDirectivePreservesExplicitPath(t *testing.T) {
	input := strings.Join([]string{
		`using:ohmyzsh/ohmyzsh kind:zsh`,
		`foo path:plugins/foo`,
	}, "\n")

	bundles, err := ParseBundles(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []Bundle{
		{
			Name:             "ohmyzsh/ohmyzsh",
			Kind:             "zsh",
			Path:             "plugins/foo",
			Line:             2,
			ExtraAnnotations: map[string]string{},
		},
	}

	if !reflect.DeepEqual(bundles, want) {
		t.Fatalf("mismatch:\n got: %#v\nwant: %#v", bundles, want)
	}
}

func TestParseBundles_UsingDirectiveAppliesToMultipleBareBundles(t *testing.T) {
	input := strings.Join([]string{
		`using:ohmyzsh/ohmyzsh path:plugins pin:abcdef`,
		`git`,
		`extract`,
	}, "\n")

	bundles, err := ParseBundles(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []Bundle{
		{
			Name:             "ohmyzsh/ohmyzsh",
			Kind:             "zsh",
			Path:             "plugins/git",
			Pin:              "abcdef",
			Line:             2,
			ExtraAnnotations: map[string]string{},
		},
		{
			Name:             "ohmyzsh/ohmyzsh",
			Kind:             "zsh",
			Path:             "plugins/extract",
			Pin:              "abcdef",
			Line:             3,
			ExtraAnnotations: map[string]string{},
		},
	}

	if !reflect.DeepEqual(bundles, want) {
		t.Fatalf("mismatch:\n got: %#v\nwant: %#v", bundles, want)
	}
}

func TestParseBundleLine_DefaultsKindToZsh(t *testing.T) {
	bundle, err := ParseBundleLine(`foo`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bundle.Kind != KindZsh {
		t.Fatalf("expected default kind %q, got %q", KindZsh, bundle.Kind)
	}
}

func TestParseBundleLine_InvalidAnnotation(t *testing.T) {
	bundle, err := ParseBundleLine(`foo kind:zsh unknown:yes`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := bundle.ExtraAnnotations["unknown"]; got != "yes" {
		t.Fatalf("expected unknown annotation preserved, got %q", got)
	}
}

func TestParseBundleLine_InvalidKind(t *testing.T) {
	_, err := ParseBundleLine(`foo kind:invalid`)
	if err == nil {
		t.Fatal("expected error for invalid kind")
	}
}

func TestParseBundleLine_AllowedKindValues(t *testing.T) {
	kinds := []string{KindZsh, KindPath, KindFpath, KindDefer, KindClone, KindAutoload}
	for _, kind := range kinds {
		t.Run(kind, func(t *testing.T) {
			bundle, err := ParseBundleLine(`foo kind:` + kind)
			if err != nil {
				t.Fatalf("unexpected error for kind %q: %v", kind, err)
			}
			if bundle.Kind != kind {
				t.Fatalf("expected kind %q, got %q", kind, bundle.Kind)
			}
		})
	}
}

func TestParseBundleLine_InvalidFpathRule(t *testing.T) {
	_, err := ParseBundleLine(`foo fpath-rule:invalid`)
	if err == nil {
		t.Fatal("expected error for invalid fpath-rule")
	}
}

func TestParseBundleLine_AllowedFpathRules(t *testing.T) {
	rules := []string{FpathRuleAppend, FpathRulePrepend}
	for _, rule := range rules {
		t.Run(rule, func(t *testing.T) {
			bundle, err := ParseBundleLine(`foo fpath-rule:` + rule)
			if err != nil {
				t.Fatalf("unexpected error for fpath-rule %q: %v", rule, err)
			}
			if bundle.FpathRule != rule {
				t.Fatalf("expected fpath-rule %q, got %q", rule, bundle.FpathRule)
			}
		})
	}
}

func TestParseBundles_InvalidLine(t *testing.T) {
	_, err := ParseBundles("foo kind\nfoo kind:zsh")
	if err == nil {
		t.Fatal("expected error for malformed line")
	}
	if !strings.Contains(err.Error(), "line 1") {
		t.Fatalf("expected line number in error, got %q", err.Error())
	}
}

func TestParseBundles_LaterLineError(t *testing.T) {
	_, err := ParseBundles("foo/bar kind:zsh\n# comment\nfoo kind:invalid")
	if err == nil {
		t.Fatal("expected error for invalid kind on later line")
	}
	if !strings.Contains(err.Error(), "line 3") {
		t.Fatalf("expected line number 3 in error, got %q", err.Error())
	}
}

func TestParseBundles_EmptyAndCommentOnly(t *testing.T) {
	bundles, err := ParseBundles("# only comment\n   \n")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bundles) != 0 {
		t.Fatalf("expected no bundles, got %d", len(bundles))
	}
}
