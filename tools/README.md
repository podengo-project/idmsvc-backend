# Go tools dependencies

> This is valid only for the tools that are installed as `go install <tool>`.

This directory keep the go tools dependencies so that we can manage them in a
centralized way, providing a better version control of the tools leveraged in
the repository.

Upsides:

- Simplify how the dependencies are installed as we only declare them in the
  `tools/tools.go` and top-level `go.mod` files.
- Save time, as now the dependencies use the same cache and not a clean
  environment, taking advantage of already downloaded packages and built parts.
- Simplify how to add new go tools to the repository, removing code duplication
  into the Makefile, and the need of creating a new makefile rule for each new
  tool to be installed; indeed no changes to the Makefile at all will be
  needed.
- As the `tools/tools.go` file has the build tag `tools`, this dependencies does
  not impact at all into the dependencies of our generated binaries by our
  repository. This reduce conflicts of dependencies between our code base and
  the go tools installed.

Downsides:

- Some tools do not support this approach, e.g. `golangci-lint`.

## Installing Go tools

There is no need to install any Go tools manually. All make targets declare
their dependencies and install necessary tools automatically. You can install
tools manually with `make install-go-tools` or `make install-tools`.

## Adding a new tool

- Add the dependency to `tools/tools.go` file.
- From the root director, run `go get "the-tool-url"`.
- From the root directory, run `go mod tidy`.
- Add a variable to `scripts/mk/variables.mk` and update the `TOOLS` variable.
  The binary name must match a substring of the Go import name.
- From the base repository directory now check that your tool install correctly
  by `make install-go-tools`.
- Now you will see the tool in your `tools/bin/` directory of your repository.

For a final check, remove all build artifacts with `make cleanall`, then run:
`make install-go-tools build`

## Special circumstances

- `oapi-codegen` is installed with `go.mod` from project root directory to
  ensure that the tool and the backend code always uses the same version of
  `oapi-codegen`.
- `golangci-lint` [does not support](https://golangci-lint.run/usage/install/#install-from-source)
   the tools approach. Instead it is installed via `go install`.

## References

- https://www.jvt.me/posts/2022/06/15/go-tools-dependency-management/
- https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
