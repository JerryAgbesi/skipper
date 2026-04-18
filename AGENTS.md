# AGENTS.md

Guidance for AI coding agents contributing to Skipper. Derived from `CONTRIBUTING.md` — follow both.

## Environment

- Requires Go 1.25+.
- Build and run locally with `make run`.
- Run `make hooks` once after cloning to enable the pre-commit lint hook (`.githooks/pre-commit`).

## Before you push

Always run these and make sure they pass:

```bash
make fmt
make lint
go test ./...
```

## Making changes

- Branch off `main` using prefixes like `feat/`, `fix/`, or `docs/` (e.g. `feat/search-improvements`).
- Keep each PR focused on one logical change.
- Follow existing code style — do not reformat unrelated code.
- Every feature must ship with tests covering it; no feature lands without test coverage.
- Add or update tests whenever behavior changes, and keep `go test ./...` green.
- Update `README.md` when user-facing behavior or flags change.

## Commit messages

Short, imperative mood, lowercase. Examples:

```
add fuzzy match scoring
fix crash when SSH config is missing
```

## Pull requests

- Open against `main`.
- Describe *what* changed and *why*; link any related issue.
- Ensure CI is green before requesting review.

## Bug reports / feature requests

Use the issue templates. Include expected vs. actual behavior, reproduction steps, and OS + `skipper --version`.

## Conduct and licensing

- Abide by the [Code of Conduct](./CODE_OF_CONDUCT.md).
- Contributions are licensed under the [Apache License 2.0](./LICENSE).
