# bundleparse

A small Go library for parsing Antidote bundle lines and bundle files.

## Antidote DSL

- Each bundle line starts with a bundle name or `using:<bundle>`.
- Bundle lines can include annotations in `key:value` form that alter the bundle behavior.
- Annotations include `kind:<{zsh,defer,fpath,path,autoload,clone}>`, `branch:<branch>`, `path:<path>`, `pin:<SHA>`, `conditional:<func>`, `autoload:<path>`, `pre:<func>`, `post:<func>`, and `fpath-rule:<{append,prepend}>`.
- The `using:<bundle>` directive sets a shared bundle path and default annotations for subsequent bare bundle names.

## Usage

- Add as a dependency: `go get github.com/getantidote/bundleparse@latest`
- Run tests: `go test ./...`
- Run benchmarks: `just benchmark`
