# About infrastructure

## Local infrastructure

- Some local infrastructure is provided by using a docker-compose
  file located at: `deployments/docker-compose.yaml`.

- Currently it supports `podman-compose` :) and `docker-compose`.

- Build your service: `make build`
- Prepare the local infrastructure by: `make compose-build`
- Start local infrastructure by: `make compose-up`
- Run the service by: `make build run`
- Run some http scripts from vscode for interacting with your api.
- Stop the service.
- Stop the local infrastructure by: `make compose-down`

Currently it supports:

- Database postresql 15
- Prometheus.
- Grafana

## Ephemeral environments

NOTE: It requires some configuration, look at: `scripts/mk/private.example.mk` to copy at `secrets/private.mk` and update with your data.

NOTE: It requires a robot token with write permissions to your personal quay.io repository.

- The repository provides initial support to deploy in Ephemeral Environment
  (aka EE).
- First of all we need to log in to the Ephemeral Environment by:
  `make ephemeral-login`
- Reserve an ephemeral namespace by: `make ephemeral-namespace-create`.
- Build and deploy the image to be used (it uses the `.rhcicd/build_deploy.sh`
  script that use the jenkins job, so we can detect situations when running this),
  by the command: `make ephemeral-deploy`
- You can check the output generated when the template at
  `deployments/clowder.yaml` is processed by: `make ephemeral-process`

> All this makefile rules are built on top of `bonfire`,
> but automate actions to get a customized experience and
> letting to deploy the current repository local state.

- You can quickly deploy a dnsutil pod into the
  EE namespace by: `make ephemeral-run-dnsutil`

## Pipelines supporting scripts

- The PoC includes the scripts to integrate with the platform jenkins jobs.
  - `.rhcicd/build_deploy.sh` build and deploy the container image for our service.
  - `.rhcicd/pr_check.sh` launch the necessary checks, that will be used to check
    every pull request in out service repository.
- The PoC includes the service descriptor that manage `clowder` operator;
  the descriptor is found at: `deployments/clowder.yaml`

## Others

### Generate dependency diagram

- We can get a high level view of the dependencies by
  running: `make generate-deps`. Thanks to @rverdile.

### Generate Entity-Relationship model

- We will be able to generate the entity-relationship from the database
  by: `make generate-data-model`

See hmscontent repository at: 
`https://github.com/content-services/content-sources-backend/`
as it implements actually a rule to generate that.

> It requires the database container up and running.
>
> It requires to add a dependency that will be installed by:
> `make install-tools`

