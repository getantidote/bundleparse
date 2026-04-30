package bundleparse

import (
	"fmt"
	"strings"
)

type Bundle struct {
	Name             string
	Kind             string
	Branch           string
	Path             string
	Pin              string
	Conditional      string
	Autoload         string
	Pre              string
	Post             string
	FpathRule        string
	Line             int
	ExtraAnnotations map[string]string
}

const (
	KeyKind        = "kind"
	KeyBranch      = "branch"
	KeyPath        = "path"
	KeyPin         = "pin"
	KeyConditional = "conditional"
	KeyAutoload    = "autoload"
	KeyPre         = "pre"
	KeyPost        = "post"
	KeyFpathRule   = "fpath-rule"

	KindZsh      = "zsh"
	KindPath     = "path"
	KindFpath    = "fpath"
	KindDefer    = "defer"
	KindClone    = "clone"
	KindAutoload = "autoload"

	FpathRuleAppend  = "append"
	FpathRulePrepend = "prepend"
)

var validKindValues = map[string]struct{}{
	KindZsh:      {},
	KindPath:     {},
	KindFpath:    {},
	KindDefer:    {},
	KindClone:    {},
	KindAutoload: {},
}

var validFpathRules = map[string]struct{}{
	FpathRuleAppend:  {},
	FpathRulePrepend: {},
}

type ParseError struct {
	Line int
	Err  error
}

func (e ParseError) Error() string {
	return fmt.Sprintf("line %d: %v", e.Line, e.Err)
}

func (e ParseError) Unwrap() error {
	return e.Err
}

// ParseBundles parses a multi-line bundle specification.
// Each non-empty line is passed through ParseLine, then validated and
// converted into a Bundle struct.
func ParseBundles(input string) ([]Bundle, error) {
	input = strings.ReplaceAll(input, "\r\n", "\n")
	lines := strings.Split(input, "\n")
	bundles := make([]Bundle, 0, len(lines))

	for idx, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parsed, err := ParseLine(line)
		if err != nil {
			return nil, ParseError{Line: idx + 1, Err: err}
		}
		if parsed.Directive != BundleDirective || parsed.Name == "" {
			continue
		}
		bundle, err := bundleFromParsed(parsed)
		if err != nil {
			return nil, ParseError{Line: idx + 1, Err: err}
		}
		bundle.Line = idx + 1
		bundles = append(bundles, bundle)
	}

	return bundles, nil
}

// ParseBundleLine parses a single bundle definition line.
func ParseBundleLine(line string) (Bundle, error) {
	parsed, err := ParseLine(line)
	if err != nil {
		return Bundle{}, err
	}
	if parsed.Directive != BundleDirective {
		return Bundle{}, nil
	}
	if parsed.Name == "" {
		return Bundle{}, nil
	}

	bundle, err := bundleFromParsed(parsed)
	if err != nil {
		return Bundle{}, err
	}
	return bundle, nil
}

func bundleFromParsed(parsed ParsedLine) (Bundle, error) {
	if parsed.Directive != BundleDirective {
		return Bundle{}, fmt.Errorf("parsed line is not a bundle")
	}
	if strings.TrimSpace(parsed.Name) == "" {
		return Bundle{}, fmt.Errorf("missing bundle name")
	}

	bundle := Bundle{
		Name:             parsed.Name,
		ExtraAnnotations: make(map[string]string, len(parsed.Annotations)),
	}

	for _, annotation := range parsed.Annotations {
		switch annotation.Key {
		case KeyKind:
			if err := validateKind(annotation.Value); err != nil {
				return Bundle{}, err
			}
			bundle.Kind = annotation.Value
		case KeyBranch:
			bundle.Branch = annotation.Value
		case KeyPath:
			bundle.Path = annotation.Value
		case KeyPin:
			bundle.Pin = annotation.Value
		case KeyConditional:
			bundle.Conditional = annotation.Value
		case KeyAutoload:
			bundle.Autoload = annotation.Value
		case KeyPre:
			bundle.Pre = annotation.Value
		case KeyPost:
			bundle.Post = annotation.Value
		case KeyFpathRule:
			if err := validateFpathRule(annotation.Value); err != nil {
				return Bundle{}, err
			}
			bundle.FpathRule = annotation.Value
		default:
			bundle.ExtraAnnotations[annotation.Key] = annotation.Value
		}
	}

	return bundle, nil
}

func validateKind(value string) error {
	if _, ok := validKindValues[value]; !ok {
		return fmt.Errorf("invalid kind %q", value)
	}
	return nil
}

func validateFpathRule(value string) error {
	if _, ok := validFpathRules[value]; !ok {
		return fmt.Errorf("invalid fpath-rule %q", value)
	}
	return nil
}
