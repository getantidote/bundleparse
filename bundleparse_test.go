package bundleparse

import (
	"reflect"
	"testing"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ParsedLine
		wantErr bool
	}{
		{
			name:  "basic bundle",
			input: `zsh-users/zsh-autosuggestions kind:zsh`,
			want: ParsedLine{
				Bundle: "zsh-users/zsh-autosuggestions",
				Annotations: []Annotation{{
					Key:   "kind",
					Value: "zsh",
				}},
			},
		},
		{
			name:  "only bundle",
			input: `zsh-users/zsh-autosuggestions`,
			want: ParsedLine{
				Bundle: "zsh-users/zsh-autosuggestions",
			},
		},
		{
			name:  "double quoted value",
			input: `foo pre:"echo hello world"`,
			want: ParsedLine{
				Bundle: "foo",
				Annotations: []Annotation{{
					Key:   "pre",
					Value: "echo hello world",
				}},
			},
		},
		{
			name:  "single quoted value",
			input: `foo post:'echo goodbye world'`,
			want: ParsedLine{
				Bundle: "foo",
				Annotations: []Annotation{{
					Key:   "post",
					Value: "echo goodbye world",
				}},
			},
		},
		{
			name:  "mixed quotes",
			input: `rupa/z pre:"echo 'hello' world" post:'echo "goodbye" world'`,
			want: ParsedLine{
				Bundle: "rupa/z",
				Annotations: []Annotation{
					{Key: "pre", Value: "echo 'hello' world"},
					{Key: "post", Value: `echo "goodbye" world`},
				},
			},
		},
		{
			name:  "escaped characters in double quotes",
			input: `foo pre:"echo \"hello\" world"`,
			want: ParsedLine{
				Bundle: "foo",
				Annotations: []Annotation{{
					Key:   "pre",
					Value: `echo "hello" world`,
				}},
			},
		},
		{
			name:  "comment ignored",
			input: `foo kind:zsh # this is a comment`,
			want: ParsedLine{
				Bundle: "foo",
				Annotations: []Annotation{{
					Key:   "kind",
					Value: "zsh",
				}},
			},
		},
		{
			name:  "trailing whitespace",
			input: "   foo    kind:zsh   ",
			want: ParsedLine{
				Bundle: "foo",
				Annotations: []Annotation{{
					Key:   "kind",
					Value: "zsh",
				}},
			},
		},
		{
			name:  "backslash escape",
			input: "foo/bar pre:echo\\ here",
			want: ParsedLine{
				Bundle: "foo/bar",
				Annotations: []Annotation{{
					Key:   "pre",
					Value: "echo here",
				}},
			},
		},
		{
			name:  "path with slashes",
			input: `ohmyzsh/ohmyzsh path:plugins/git`,
			want: ParsedLine{
				Bundle: "ohmyzsh/ohmyzsh",
				Annotations: []Annotation{{
					Key:   "path",
					Value: "plugins/git",
				}},
			},
		},
		{
			name:  "empty line",
			input: ``,
			want:  ParsedLine{},
		},
		{
			name:    "missing key value colon",
			input:   `foo kind`,
			wantErr: true,
		},
		{
			name:    "missing value",
			input:   `foo kind:`,
			wantErr: true,
		},
		{
			name:    "unterminated double quote",
			input:   `foo pre:"echo hello`,
			wantErr: true,
		},
		{
			name:    "unterminated single quote",
			input:   `foo pre:'echo hello`,
			wantErr: true,
		},
		{
			name:    "unterminated escape",
			input:   `foo pre:"echo hello\`,
			wantErr: true,
		},
		{
			name:  "multiple annotations",
			input: `foo/bar kind:zsh branch:main path:plugins/git`,
			want: ParsedLine{
				Bundle: "foo/bar",
				Annotations: []Annotation{
					{Key: "kind", Value: "zsh"},
					{Key: "branch", Value: "main"},
					{Key: "path", Value: "plugins/git"},
				},
			},
		},
		{
			name:  "bundle-like token",
			input: `kind:zsh`,
			want: ParsedLine{
				Bundle: "kind:zsh",
			},
		},
		{
			name:  "Full SSH URL",
			input: `git@github.com:zsh-users/zsh-autosuggestions kind:zsh post:a:b:c  # comment`,
			want: ParsedLine{
				Bundle: "git@github.com:zsh-users/zsh-autosuggestions",
				Annotations: []Annotation{
					{Key: "kind", Value: "zsh"},
					{Key: "post", Value: "a:b:c"},
				},
			},
		},
		{
			name:  "Full URL",
			input: `https://github.com/zsh-users/zsh-autosuggestions`,
			want: ParsedLine{
				Bundle: "https://github.com/zsh-users/zsh-autosuggestions",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLine(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil (result=%v)", got)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("mismatch:\n got: %#v\nwant: %#v", got, tt.want)
			}
		})
	}
}
