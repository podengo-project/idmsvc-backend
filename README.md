# idmsvc-backend (FreeIPA Domain Join)

**This project is currently inactive.  Issues and pull requests will
not be attended to.**

## Getting started

**Pre-requisites**:

- golang 1.24 (not `gcc-go`)
- docker or podman (>4.0.0)
- docker-compose or podman-compose
- python3
- openshift client [Installing OpenShift Client](https://docs.openshift.com/container-platform/4.12/cli_reference/openshift_cli/getting-started-cli.html#installing-openshift-cli).

Packages for fedora 41:

```sh
$ sudo dnf upgrade
$ sudo dnf install git golang podman podman-compose delve
$ sudo dnf remove gcc-go
```

(Optional) Installing VSCode by repository on fedora 41:

```sh
$ sudo rpm --import https://packages.microsoft.com/keys/microsoft.asc
$ cat <<EOF | sudo tee /etc/yum.repos.d/vscode.repo
[code]
name=Visual Studio Code
baseurl=https://packages.microsoft.com/yumrepos/vscode
enabled=1
gpgcheck=1
gpgkey=https://packages.microsoft.com/keys/microsoft.asc
EOF
$ sudo dnf check-update
$ sudo dnf install code
```

(Optional) Recomended extensions for vscode:

- Go (highly recommended)
- REST Client

Once tasks:

- Install used tools: `make install-tools`
- Create your `configs/config.yaml` file: `cp -vf configs/config.example.yaml configs/config.yaml`
- Create your `secrets/private.mk` file:
  ```sh
  $ mkdir secrets
  $ cp -vf scripts/mk/private.example.mk secrets/private.mk
  ```
- Follow the instructions in `secrets/private.mk` to set up Quay repository
  and Red Hat Repository access.
- Create your `configs/bonfire.yaml` file: `cp -vf configs/bonfire.example.yaml configs/bonfire.yaml`
- Add a `GITHUB_TOKEN` to `~/.config/bonfire/env` if deploying frequently at ephemeral environment,
  to avoid API access rate limit.

> Directory `secrets/` is set to be ignored by git and docker.

----

- Build by: `make build`
- Lint code by: `make lint`
- Start local infrastructure by: `make compose-up mock-rbac-up`
- Launch tests by: `make test`
- Run by: `make run`
- Try locally by running scripts at: `./test/scripts/local-*.sh`
  You can override the used xrhid by: `XRHID_AS="service-account" ./test/scripts/local-*.sh`
- Stop local infrastructure by: `make compose-down mock-rbac-down`
- Clean local infrastructure by: `make compose-clean`
- Print out useful rules by: `make help`

----

- Run with specific rbac profile and use local rbac mock:

  ```sh
  # If we use APP_CLIENTS_RBAC_PROFILE=custom
  # and we are checking custom changes, we would need
  # to restart mock-rbac, to simplify we can do:
  $ make compose-clean compose-build clean build compose-up

  # And finally start the service by:
  $ make run APP_CLIENTS_RBAC_PROFILE=domain-readonly
  $ curl "http://localhost:8020/api/rbac/v1/access?application=idmsvc"
  $ ./test/scripts/local-domains-list.sh   # Will success
  $ ./test/scripts/local-domains-token.sh  # Will fail as unauthorized
  ```

> - Bear in mind the rbac mock is started with the local infra.
> - Be aware to update some change in the rbac mock (such as some
>   change on the custom.yaml profile), you will need to rebuild
>   the container by `make mock-rbac-build` and restarting the
>   local infra by `make compose-clean compose-up`.

----

For ephemeral environment look at: [DEVELOPMENT.md](DEVELOPMENT.md) file.

## Project layout

```raw
internal/   Define the internal application components
├── api
│   ├── header: Hold code related with the http headers.
│   ├── metrics: Hold the service interface and definition for the
│   │            /metrics endpoint.
│   ├── openapi: Hold service interface and definitions for the
│   │            /openapi.json endpoint.
│   ├── private: The code generated for the private api.
│   └── public: The code generated for the public api (types, http
│               framework server specific, spec).
├── config: Hold the configuration structure and functions to read it.
├── domain/model: Define the business data model of the application;
│                 in this scenario match the database, so it uses the
│                 model for gorm.
├── handler: Hold application and handler interfaces.
│   └── impl: Implementation for the application interface.
├── infrastructure: specific code coupled to the http framework.
│   ├── context: helpers to set/get data to/from the go context.
│   ├── datastore: helpers to initialize database connector, and
│   │              run migrations.
│   ├── event: (delete) infrastructure to deal with asynchronous
│   │          processors in a similar way as the http handlers.
│   ├── logger: helper to start the log infrastructure using slog.
│   ├── middleware: all the middleware components comes here.
│   ├── router: wire the route of the service composing the different
│   │           api groups, and adding the middlewares.
│   ├── secrets: logic related with secrets.
│   ├── service: define the Service interface, understanding each service
│   │   │        as a some listener at a port, or a client broker to
│   │   │        process asynchronous events.
│   │   └── impl: Implement a Service for the application and different
|   │             listeners (api, metrics, kafka consumer)
|   └── token: logic to deal with domain and hostconf tokens.
├── interface: Define the interfaces for `interactor`,
│              `repository` and `presenter` components.
├── test: All the helpers for tests are here
│   ├── assert: provide new asserts to simplify test expectations.
│   ├── builder: make easier to generate new filled business data
│   │            and API structures.
│   ├── client: some helpers to deal with client handler tests.
│   ├── mock: Store all the generated mocks for the interfaces,
│   │         keeping the same directory structure
│   ├── perf: Performance tests 
│   ├── smoke: Smoke tests for the API.
│   └── sql: Helpers for testing the database repository. Help on
│            preparing the expectation for the database.
└── usecase: specific implementation for the `interface` directory
             for interactor, presenter and repository components.

cmd/      Define the binaries generated

api/      Define the public and private openapi specification

deployments/   Hold descriptors to deploy with clowder and local
               infrastructure with {podman,docker}-compose.

scripts/  Store useful scripts for the repository
├── db
│   └── migrations: sql scripts for the migrations.
├── http:  Hold .http files to quickly check the API
└── mk:  Hold all the makefile scripts
```

See: https://github.com/golang-standards/project-layout

## Architecture and Design docs

See: [Architecture and design](docs/ARCHITECTURE.md).

## Design API

You can design your API importing the `public.openapi.yaml` file
at [apicurito](https://console.redhat.com/application-services/api-designer/designs).

When you have made your changes, then do click **Actions** > **Download Design**,
and copy the downloaded file as `api/public.openapi.yaml`.

Now from the base of the repository now update code generated by `make generate-api`.

Finally, edit your source code to adapt to the changes (if necessary)
and update your unit tests to cover new code and update to the changes.

## Contributing

See [CONTRIBUTING.md](docs/CONTRIBUTING.md) guide.

## Aknowledgements

Thanks to all mates from hmscontent (content-sources), without all of you
this would not be possible; every moment at hmscontent has been adding
a grain of sand to build the path on this repository.

Big thanks to everybody.

## Generating reference code documentation

TODO Ask if leverage godoc for this

## References

- https://github.com/content-services/content-sources-backend
- https://manakuro.medium.com/clean-architecture-with-go-bce409427d31
- https://articles.wesionary.team/environment-variable-configuration-in-your-golang-project-using-viper-4e8289ef664d
- https://github.com/RedHatInsights/playbook-dispatcher
- https://faun.pub/dependency-injection-in-go-the-better-way-833e2cc65fab

Tech stack:

- Echo Framework: https://echo.labstack.com
- Logs:
  - https://pkg.go.dev/log/slog
  - https://lukas.zapletalovi.com/posts/2023/about-structured-logging-in-go121/
- Database: https://gorm.io/docs/index.html
- Kafka Client Library: https://github.com/confluentinc/confluent-kafka-go
- Testing:
  - Testify: https://pkg.go.dev/github.com/stretchr/testify
  - Mockery: https://github.com/vektra/mockery
  - SqlMock: https://github.com/DATA-DOG/go-sqlmock
    - Example of use (not working): https://github.com/manakuro/golang-clean-architecture/blob/202ac7e826bbecb8a104dad24a6730db160c9ec8/interface/repository/user_repository_test.go#L18
    - Example that drive to the worked version: https://github.com/go-gorm/gorm/issues/3565
- Generators:
  - oapi-generator: https://github.com/deepmap/oapi-codegen   (Generate boilerplate for openapi)
  - go-jsonschema: https://git.sr.ht/~emersion/go-jsonschema  (Generate types for event schemas)
  - mockery: https://github.com/vektra/mockery                (Generate mocks from golang interfaces)

## References to investigate

- https://developers.redhat.com/articles/2021/06/02/simulating-cloudevents-asyncapi-and-microcks
- https://access.redhat.com/documentation/en-us/red_hat_openshift_api_designer/1/guide/f4d2c457-2d1a-4bc9-913f-522847405c45
- https://github.com/rookie-ninja/rk-boot

## Tools

Http clients
- https://httpie.io/docs/cli
- https://github.com/AnWeber/httpyac

Validate API
- https://quobix.com/vacuum/
