# AGENTS.md

Instructions for coding agents working in this repo.

## Project

GIP converts GitHub Issues and PRs into structured markdown files. See
README.md for the product description.

## State

Go CLI (`go build ./...`), single `main` package, with the official Go MCP SDK
for the `mcp` subcommand. Plus
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
  not performed by an agent. It is a direct local fast-forward push, **not**
  a GitHub PR — GitHub's PR merge UI only offers merge commit, squash, or
  rebase, none of which fast-forward, so using it here would either add a
  merge commit or diverge `main`/`develop` SHAs and break the `--ff-only`
  sync below. If branch protection blocks direct pushes to `main`, the owner
  needs an exception for this step, or must fall back to a "Merge commit" PR
  (accepting the extra merge commit).

  ```bash
  git fetch origin
  git checkout develop
  git merge --ff-only origin/main
  git checkout main
  git merge --ff-only origin/develop
  git push origin main
  ```

- To sync `develop` after a `main` merge:

  ```bash
  git fetch origin
  git checkout develop
  git merge --ff-only origin/main
  git push origin develop
  ```

- To cut a release, after `develop` has been merged into `main`: tag `main`
  at the merge commit with the next `vX.Y.Z` (see README's Versioning
  section for the version-bump rules) and push the tag — this is what
  triggers the release binaries workflow. Not automated, not performed by
  an agent.

  ```bash
  git checkout main
  git pull origin main
  git tag -a vX.Y.Z -m "vX.Y.Z"
  git push origin vX.Y.Z
  ```

## Issues and PRs

- Every change starts as a GitHub issue. A PR must link at least one issue;
  arbitrary PRs with no linked issue are discouraged.
- A PR can solve one or more issues. For each linked issue, use a closing
  keyword (`Closes #N`) only if the PR fully resolves that issue — `develop`
  is this repo's default branch, so the keyword auto-closes it on merge. For
  an issue the PR only partially resolves, use `Addresses #N (partial)`
  instead; never a closing keyword for an issue that isn't fully solved.
  These can mix in one PR, e.g. `Closes #12, Addresses #14 (partial)`.
- Use the templates in `.github/ISSUE_TEMPLATE/` (bug report, feature
  request) and `.github/pull_request_template.md` when opening issues/PRs.

## Commit messages

- Prefer short [Conventional Commits](https://www.conventionalcommits.org/)
  format: `type(scope): summary` (e.g. `fix(parser): handle empty issue body`).
- If an agent creates a commit, GitHub issue, or PR, do not add any
  attribution, co-author, or session/tool info — no "Generated with", no
  `Co-Authored-By` lines. Plain conventional-commit message only.
