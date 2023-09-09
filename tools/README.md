# go tools dependencies

> This is valid only for the tools that are installed as `go install <tool>`.

This directory keep the go tools dependencies so that we can manage them in a
centralized way, providing a better version control of the tools leveraged in
the repository.

Upsides:

- Simplify how the dependencies are installed as we only declare them in the
  `tools/tools.go` and `tools/go.mod` files.
- Save time, as now the dependencies use the same cache and not a clean
  environment, taking advantage of already downloaded packages and built parts.
- Simplify how to add new go tools to the repository, removing code duplication
  into the makefiles, and the need of creating a new makefile rule for each new
  tool to be installed; indeed no changes to the makefiles at all will be
  needed.
- As the `tools/tools.go` file has the build tag `tools`, this dependencies does
  not impact at all into the dependencies of our generated binaries by our
  repository. This reduce conflicts of dependencies between our code base and
  the go tools installed.

Downsides:

- Not detected yet.

## Installing the go tools

Use `make install-go-tools` which is a new rule that use the new system; old
rule `make install-tools` is kept as deprecated until everything is cleaned-up
and it is seen the new rule works as expected.

## Adding a new tool

- Add the dependency to `tools/go.mod` file.
- From the `tools/` directory, run `go mod tidy`.
- From the base repository directory now check that your tool install correctly
  by `make install-go-tools`.
- Now you will see the tool in your `bin/` directory of your repository.

For a final check, completely remove the `bin/` directory, and run:
`make install-go-tools build`

## References

- https://www.jvt.me/posts/2022/06/15/go-tools-dependency-management/
- https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
