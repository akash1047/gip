# AGENTS.md

Instructions for coding agents working in this repo.

## Project

GIP converts GitHub Issues and PRs into structured markdown files. See
README.md for the product description.

## State

Go CLI (`go build ./...`), single `main` package, stdlib only so far. Plus
`.githooks/pre-commit`, a Bash warning hook for staged Markdown lines over 80
characters. Enable it with `git config core.hooksPath .githooks`.

## Guidelines

- Don't scaffold speculative structure (config files, CI, abstractions) before
  there's real code that needs it.
- When adding the first implementation, pick the smallest workable stack
  (single script/CLI is fine) rather than a framework, unless the user asks
  for one.
- Keep this file and README.md updated as the project takes shape (actual
  commands, entry points, structure) instead of documenting intentions.

## Git branching

- Two long-lived branches: `main` and `develop`. Both must have linear
  history — no merge commits into them.
- All changes land via short-lived PR branches, raised against `develop`.
- Default merge policy for `develop` ← PR branch: **squash merge**.
- `develop` → `main` is done manually by the human owner. Not automated,
  not performed by an agent.
- To sync `develop` after a `main` merge:

  ```bash
  git fetch origin
  git checkout develop
  git merge --ff-only origin/main
  git push origin develop
  ```

## Issues and PRs

- Every change starts as a GitHub issue. By default, every PR should close
  exactly one issue (`Closes #N`). As an explicit exception, a PR may solve
  multiple issues, or partially solve one or more issues, as long as it is
  explicit about it (e.g. `Addresses #N, #M` or `Addresses #N (partial)`) —
  but such PRs must NOT use closing keywords (`Closes #N`). Arbitrary PRs
  with no linked issue are discouraged.
- Use the templates in `.github/ISSUE_TEMPLATE/` (bug report, feature
  request) and `.github/pull_request_template.md` when opening issues/PRs.

## Commit messages

- Prefer short [Conventional Commits](https://www.conventionalcommits.org/)
  format: `type(scope): summary` (e.g. `fix(parser): handle empty issue body`).
- If an agent creates a commit, GitHub issue, or PR, do not add any
  attribution, co-author, or session/tool info — no "Generated with", no
  `Co-Authored-By` lines. Plain conventional-commit message only.
