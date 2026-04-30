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

		bundle, err := ParseBundleLine(line)
		if err != nil {
			return nil, ParseError{Line: idx + 1, Err: err}
		}
		if bundle.Name == "" {
			continue
		}
		bundles = append(bundles, bundle)
	}

	return bundles, nil
}

// ParseBundleLine parses a single bundle definition line.
func ParseBundleLine(line string) (Bundle, error) {
	values, err := ParseLine(line)
	if err != nil {
		return Bundle{}, err
	}
	if len(values) == 0 {
		return Bundle{}, nil
	}

	bundle, err := bundleFromMap(values)
	if err != nil {
		return Bundle{}, err
	}
	return bundle, nil
}

func bundleFromMap(values map[string]string) (Bundle, error) {
	name, ok := values["bundle"]
	if !ok || strings.TrimSpace(name) == "" {
		return Bundle{}, fmt.Errorf("missing bundle name")
	}

	bundle := Bundle{
		Name:             name,
		ExtraAnnotations: make(map[string]string, len(values)-1),
	}

	for key, value := range values {
		if key == "bundle" {
			continue
		}

		switch key {
		case KeyKind:
			if err := validateKind(value); err != nil {
				return Bundle{}, err
			}
			bundle.Kind = value
		case KeyBranch:
			bundle.Branch = value
		case KeyPath:
			bundle.Path = value
		case KeyPin:
			bundle.Pin = value
		case KeyConditional:
			bundle.Conditional = value
		case KeyAutoload:
			bundle.Autoload = value
		case KeyPre:
			bundle.Pre = value
		case KeyPost:
			bundle.Post = value
		case KeyFpathRule:
			if err := validateFpathRule(value); err != nil {
				return Bundle{}, err
			}
			bundle.FpathRule = value
		default:
			bundle.ExtraAnnotations[key] = value
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
