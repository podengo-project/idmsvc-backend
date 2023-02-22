# Contributing guide

The repository is using [Github flow](https://docs.github.com/en/get-started/quickstart/github-flow).

- Fork the repository in your namespace.
- Clone the repository locally.
- Create a branch.
- Add changes:
  - If you change the api, run `make generate-api`.
  - If you add/update golang interfaces, run `make generate-mock`
  - If you add/update a kafka topic, run `make generate-events`
- Check everything build: `make build`
- Check locally by using: `make compose-clean compose-up run`
- Check it deploys and works in ephemeral by: `make ephemeral-deploy`
- Add unit tests, if your change add a new interface, generate the mocks.
  by `make generate-mocks`; they will be generated at `internal/test/mock`
  package.
- Check unit tests and linters are pasing by `make test lint`
- Rebase and push your changes, and create a MR or PR.
- Update changes from the review until you get an ACK.
- Merge your changes :)
