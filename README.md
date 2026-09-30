# GIP — GitHub Issues & Pull Requests

Converts GitHub Issues and Pull Requests into structured markdown files.

## Status

Go CLI with `fetch`, `doctor`, and `mcp` subcommands.

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
./gip fetch owner/repo#123
```

Optionally set a fine-grained GitHub PAT for a higher API rate limit. Use
`doctor` to check that the token works:

```bash
export GIP_GITHUB_TOKEN=your-token
./gip doctor
./gip fetch owner/repo#123
./gip fetch -o owner-repo-123.md owner/repo#123
./gip fetch --format=metadata,title,body owner/repo#123
```

`fetch` prints a title heading and the body (including comments) to stdout by
default. Use `--format` with a comma-separated list of `metadata` (YAML front
matter), `title`, and `body` to choose which sections to print and in what
order — they're concatenated in the order given, e.g. `--format=body,title`
prints the body before the title. Use `-o` or `--output` with a path to write
a file instead. `doctor` checks that a token is set and authenticates against
the GitHub API. `GIP_GITHUB_TOKEN` takes precedence over `GITHUB_TOKEN`.

Run the MCP server over stdio (the default) or Streamable HTTP at `/mcp`:

```bash
./gip mcp
./gip mcp --transport=http --addr=:8080
# equivalent: ./gip mcp --http=:8080
```

The server exposes `gip_fetch` with a required `ref` (`owner/repo#123`) and
optional `format` (default `title,body`). It returns the same Markdown as
`fetch`. Stdio stdout carries MCP messages only; diagnostics go to stderr.

## Versioning

Releases are tagged `vX.Y.Z` following [SemVer](https://semver.org). The
project is pre-1.0 (starting at `v0.1.0`): the API/CLI surface has no
stability guarantee yet, so a `MINOR` bump marks a breaking change and
`PATCH` marks a backward-compatible fix or docs change, while `MAJOR` stays
at 0 until the maintainer decides the surface is stable enough for `v1.0.0`.
Tags are cut on `main` after a `develop` → `main` merge.

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
