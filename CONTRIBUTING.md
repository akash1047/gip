# Contributing to GIP

Thank you for your interest in contributing to GIP!

## Current Project State

GIP is in its early stages. There is currently no chosen language,
framework, or build system. Please avoid scaffolding speculative structure
(config files, CI, abstractions) before there is real code requiring it.
When introducing implementation work, keep the stack as minimal and
workable as possible.

## Code of Conduct

Everyone contributing to GIP is expected to adhere to our
[Code of Conduct](CODE_OF_CONDUCT.md).

## Workflow Overview

1. **Every change starts with an issue.**
   - Open a feature request or bug report using the templates in
     [`.github/ISSUE_TEMPLATE/`](.github/ISSUE_TEMPLATE/).
   - Discuss proposed changes in the issue before starting implementation.
2. **By default, every pull request should close exactly one issue.**
   - Link the issue in your PR description using `Closes #N`.
   - As an explicit exception, a PR may solve multiple issues, or partially
     solve one or more issues, as long as it is explicit about it (e.g.
     `Addresses #N, #M` or `Addresses #N (partial)`) — but such PRs must
     NOT use closing keywords (`Closes #N`).
   - Arbitrary pull requests without an associated issue are discouraged.

## Git Branching Model

- **Long-lived branches:** `main` and `develop`. Both maintain a strictly
  linear history (no merge commits).
- **PR branches:** All work should be done in short-lived branches created
  from and targeted against `develop`.
- **Merging:** Pull requests into `develop` are squash-merged.
- **Releases:** Merging `develop` into `main` is performed manually by the
  repository maintainer.

### Workflow Example

```bash
# Update local develop branch
git checkout develop
git pull origin develop

# Create a short-lived feature or fix branch
git checkout -b feat/your-feature-name

# Make changes and commit
git add .
git commit -m "feat(scope): short description"

# Push to your fork or branch and open a PR against develop
git push -u origin feat/your-feature-name
```

## Commit Message Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/)
specification:

- Format: `type(scope): summary` (e.g. `feat(parser): add markdown serializer`,
  `fix(cli): handle missing token error`).
- Use lowercase for types: `feat`, `fix`, `docs`, `chore`, `refactor`, `test`,
  etc.
- Keep summaries concise and imperative.

## Submitting Pull Requests

- Open pull requests against the `develop` branch.
- Fill out the [Pull Request Template](.github/pull_request_template.md)
  completely, including the issue reference (`Closes #N`, or
  `Addresses #N, #M` / `Addresses #N (partial)`) and a brief summary of
  verification.
