# GIP — GitHub Issues & Pull Requests

Converts GitHub Issues and Pull Requests into structured markdown files.

## Status

Early stage — Go CLI in progress. No subcommands yet.

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
