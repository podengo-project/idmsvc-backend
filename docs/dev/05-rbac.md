# RBAC

This section explains how to use the RBAC mock integration to let us
check the operations locally. This component supports automation of
the smoke tests and future integration tests, so we can check
different scenarios as we need.

The rbac mock loads 4 different profiles.  The mock service at
start-up loads the profile specified in the
`APP_CLIENTS_RBAC_PROFILE` environment variable.  The profiles are:

- `super-admin`: only contains `*:*:*`;
  this is different from the `domain-admin` profile, and exists
  for checking that the wildcards are properly managed.
- `domain-admin` (default): the default Domain Admin permissions.
- `domain-readonly`: the default Domain Read permissions.
- `domain-no-perms`: does not contain any permission at all.
- `custom`: this is provided for development reasons, if we want to
  customize the permissions.  The custom permission set is
  hard-coded in
  `internal/infrastructure/service/impl/mock/rbac/impl/custom.yaml`.
  If you change them, rebuild and restart the mock service.

## Using the rbac mock

- Start the rbac mock by: `make mock-rbac-up`
- Stop the rbac mock by: `make mock-rbac-down`
- Start with a different profile:
  `make mock-rbac-up APP_CLIENTS_RBAC_PROFILE=domain-readonly`

## Updating the rbac mock image

Once it is created, no additional changes, but if that were the
case, you can run:

```sh
$ export MOCK_RBAC_CONTAINER="quay.io/podengo/mock-rbac:latest"
$ make compose-clean mock-rbac-build \
    QUAY_EXPIRATION="never"
$ podman tag "${MOCK_RBAC_CONTAINER}" YOURTAG
$ podman push YOURTAG
```

> Be aware to update the tag at `scripts/mk/mock-rbac.mk` file
> by setting `MOCK_RBAC_CONTAINER` value.

## Using on the smoke and integration tests

- It is automatically started with the suite test. To avoid
  a port conflict, the tests launch the mock rbac service
  listening on port 8021.
- By default the profile used is the Domain Administrator.
- The suite has the `As(...)` public method that can be used
  to set the RBAC profile and indicate the XRHID identity to
  used for the requests. It could be used as the below:
  ```golang
  // As Admin using a user identity
  s.As(RBACAdmin, XRHIDUser)

  // As Viewer role using a Service Account
  s.As(RBACViewer, XRHIDServiceAccount)
  ```
  The order above is exchangable, and dynamic as a type check is used
  internally; aditionally, if more profiles are indicated, the last
  in the variadic argument override the previous ones.
- The current infrastructure for the tests does not support running
  the tests in parallel, because no random port is choosen currently
  when the mock rbac service is started.

## Debugging from vscode

- Be sure your rbac mock is using the expected profile (see above).
- When you are using the `custom` profile, you could want
  to modify the list of permissions at the
  `internal/infrastructure/service/impl/mock/rbac/impl/custom.yaml`
  file; be sure you restart the rbac mock container after
  change this file by `make compose-clean mock-rbac-build compose-up`.
- When you are ready, launch the debugger, set your breakpoints
  at `internal/infrastructure/middleware/rbac.go` file, and
  debug your code step by step.
- To raise the breakpoint, some of the `./test/scripts/local-*`
  files could be useful; or open the [scripts/http/public.http](../../scripts/http/public.http)
  file on your IDE.
