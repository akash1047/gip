# GIP — GitHub Issues & Pull Requests

Converts GitHub Issues and Pull Requests into structured markdown files.

## Status

Go CLI with `fetch` and `doctor` subcommands.

## Build

Build from the repository root; the binary is written there and can be run as
`./gip`:

```bash
go build -o gip .
```

Go can cross-compile without extra tooling. Set `GOOS` and `GOARCH` for the
target; these examples build amd64 binaries for Linux, macOS, and Windows:

```bash
GOOS=linux GOARCH=amd64 go build -o gip-linux-amd64 .
```

```bash
GOOS=darwin GOARCH=amd64 go build -o gip-darwin-amd64 .
```

```bash
GOOS=windows GOARCH=amd64 go build -o gip-windows-amd64.exe .
```

## Usage

Fetch a public issue or PR without a token:

```bash
gip fetch owner/repo#123
```

Optionally set a fine-grained GitHub PAT for a higher API rate limit. Use
`doctor` to check that the token works:

```bash
export GIP_GITHUB_TOKEN=your-token
gip doctor
gip fetch owner/repo#123
gip fetch -o owner-repo-123.md owner/repo#123
```

`fetch` prints markdown (frontmatter + body) to stdout by default. Use
`-o` or `--output` with a path to write a file instead. `doctor` checks
that `GIP_GITHUB_TOKEN` is set and authenticates against the GitHub API.

## Idea

Given a GitHub repo (and optionally an issue/PR number), fetch the issue or
PR (title, body, comments, metadata) via the GitHub API and write it out as
a markdown file with consistent structure (frontmatter + body), suitable for
archiving, indexing, or feeding into other tools.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines,
branching workflow, and code of conduct.

To enable the optional warning for staged Markdown lines over 80 characters:

```bash
git config core.hooksPath .githooks
```

The hook prints `file:line` warnings and never blocks a commit. Check it with
a temporary staged file:

```bash
tmp=$(mktemp ./hook-check-XXXXXX.md)
printf '%081d\nshort\n' 0 > "$tmp"; git add "$tmp"
.githooks/pre-commit; echo "exit: $?"
git restore --staged "$tmp"; rm "$tmp"
```

## License

[MIT](LICENSE)
