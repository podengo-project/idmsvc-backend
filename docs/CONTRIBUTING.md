# Contributing to idmsvc-backend

The main way to contribute is providing changes to the repository, but not the
only one. If you see something wrong in the repository, create an issue for the
bug is another way you can help.

## Getting Started

The repository is using [Github flow][github_flow].

- [Fork][fork_repo] the repository in your namespace.
- Clone the repository locally.
- Create a branch: `git checkout -b my-cool-changes`
- Add changes:
  - If you change the api, run `make generate-api`.
  - If you add/update golang interfaces, run `make generate-mock`
  - If you add/update a kafka topic, run `make generate-events`
- Check everything build: `make build`
- Run locally by using: `make compose-clean clean build compose-up run`
- (optional)Check it deploys and works in ephemeral by: `make ephemeral-deploy`
- Add tests. See [TESTING.md](./dev/TESTING.md).
- Check everything is ok: `make compose-clean clean build lint compose-up test`
- Rebase, push changes to your forked repository, and create a PR.
- Update changes from the review until you get an ACK.
- Merge your changes :)

## Reporting bugs

- **Be sure you are using the last version of our service**.
- Please if you think this bug is a security issue, see [SECURITY.md](../SECURITY.md)
- Else, before create a new issue, please [search][search_issues]
into the current issues to check if it already exists.

Then it looks a new bug, so please create a new issue and fill the next template.

<!--
TODO When create an issue is enabled in GitHub,
     remove the template below and put in a
     github template.
-->

```markdown
### Description

Tell us a summary of the bug.

### Steps to replay

- I do action step 1.
- I do action step 2.
- I do action step 3.

What frequency you can replay the issue? (always / specific env / random)

### What is the observed wrong behavior

Tell us what is wrong.

### What is the expected behavior

Tell us what you were expecting instead of the wron behavior.

### Additional information

Please attach any further information such as:

- What commit did you observed this? `git rev-parse --short HEAD`
- What API version did you use? `(cd api; git rev-parse --short HEAD)`
- Copy & paste log blocks.
- Configuration that could impact.
- Any other additional information that could be useful to replay and
  analyse the issue.
```

Thank you for contributing to get a better software! we will study
the issue as soon as possible!

## Commit messages

Follow the [conventional commits guidelines][conventional_commits] to *make
reviews easier* and to make the VCS/git logs more valuable. The general
structure of a commit message is:

```
<type>[(optional scope)]: <description>

[optional body]

[optional footer(s)]
```

for instance

```raw
fix(HMS-9999): response 201 when a domain is registered

This change modified the status code for a success response when it
registers a domain, returning a 201 (Created) status code, instead
of 200 (Ok).

Signed-off-by: Alejandro Visiedo <avisiedo@redhat.com>
```

For a commit that has a github issue scope, the title would be:

```raw
fix(9999): response 201 when a domain is registered
```

Further information:

- Prefix the commit subject with one of these [_types_](prefix_types):
    - `build`, `ci`, `docs`, `feat`, `fix`, `perf`, `refactor`, `revert`,
      `test`, `style`, `chore`.
    - You can **ignore this for "fixup" commits** or any commits you expect to be 
      squashed.
- Append optional _scope_:
  - For a jira ticket for instance:
    ```raw
    fix(HMS-9999): response 201 when a domain is registered
    ```
  - For a github issue for instance:
    ```raw
    fix(9999): response 201 when a domain is registered
    ```
- _Commit title_ shouldn't start with a capital letter or end in a period.
- Use the _imperative voice_: "Fix bug" rather than "Fixed bug" or "Fixes bug."
- Try to keep the first line under 72 characters.
- A blank line must follow the subject.
- Breaking API changes must be indicated by
    1. "!" after the type/scope, and
    2. a "BREAKING CHANGE" footer describing the change.
       Example:
       ```
       refactor(HMS-9999)!: drop support for Python 2

       BREAKING CHANGE: refactor to use Python 3 features since Python 2
       is no longer supported.
       ```

### Automated builds (CI)

Each pull request must pass the automated builds, which execute linters, build
and tests (unit and smoke).

## Coding

[effective go][effective_go]

### Style

See `.editorconfig` for indentation and textwidth
Use `gofumpt` to format Go files!

[prefix_types]: https://github.com/commitizen/conventional-commit-types/blob/master/index.json
[github_flow]: https://docs.github.com/en/get-started/quickstart/github-flow
[fork_repo]: https://github.com/podengo-project/idmsvc-backend/fork
[conventional_commits]: https://www.conventionalcommits.org
[effective_go]: https://go.dev/doc/effective_go
[search_issues]: https://github.com/podengo-project/idmsvc-backend/pulls
