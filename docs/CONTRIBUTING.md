# Contributing to idmsvc-backend

## Getting Started

The repository is using [Github flow](https://docs.github.com/en/get-started/quickstart/github-flow).

- [Fork](https://github.com/podengo-project/idmsvc-backend/fork) the repository in your namespace.
- Clone the repository locally.
- Create a branch: `git checkout -b my-cool-changes`
- Add changes:
  - If you change the api, run `make generate-api`.
  - If you add/update golang interfaces, run `make generate-mock`
  - If you add/update a kafka topic, run `make generate-events`
- Check everything build: `make build`
- Check locally by using: `make compose-clean compose-up run`
- Check it deploys and works in ephemeral by: `make ephemeral-deploy`
- Add tests. See [TESTING.md](./dev/TESTING.md).
- Check everything is ok: `make compose-clean clean build lint vet compose-up test`
- Rebase and push your changes, and create a MR or PR.
- Update changes from the review until you get an ACK.
- Merge your changes :)

## Reporting bugs

TODO


## Guidelines

### Commit messages

Follow the [conventional commits guidelines][conventional_commits] to *make
reviews easier* and to make the VCS/git logs more valuable. The general
structure of a commit message is:

```
<type>([optional scope]): <description>

[optional body]

[optional footer(s)]
```

- Prefix the commit subject with one of these [_types_](https://github.com/commitizen/conventional-commit-types/blob/master/index.json):
    - `build`, `ci`, `docs`, `feat`, `fix`, `perf`, `refactor`, `revert`,
      `test`, `style`, `chore`.
    - You can **ignore this for "fixup" commits** or any commits you expect to be squashed.
- Append optional scope to _type_ such as `(db)`, `(inventory)`, `(client)`, …:
  - `(interactor)`: Scope the change to the input adapter components.
  - `(model)`: Scope the changes to `internal/module` and `scripts/db`.
  - `(repository)`: Scope the change to the repository components.
  - `(presenter)`: Scope the change to the output adapter components.
  - `(client)`: Scope the change to the client components.
  - `(middleware)`: Scope the changes to the middleware components.
  - `(route)`: Scope the changes to the router components.
- _Description_ shouldn't start with a capital letter or end in a period.
- Use the _imperative voice_: "Fix bug" rather than "Fixed bug" or "Fixes bug."
- Try to keep the first line under 72 characters.
- A blank line must follow the subject.
- Breaking API changes must be indicated by
    1. "!" after the type/scope, and
    2. a "BREAKING CHANGE" footer describing the change.
       Example:
       ```
       refactor(provider)!: drop support for Python 2

       BREAKING CHANGE: refactor to use Python 3 features since Python 2 is no longer supported.
       ```

### Automated builds (CI)

Each pull request must pass the automated builds.


## Coding

### Style

See `.editorconfig` for indentation and textwidth
Use `gofumpt` to format Go files!

This repository is

[conventional_commits]: https://www.conventionalcommits.org
