# RBAC

This section explain how to use the RBAC mock integration
to let us check the operations locally. This component exists
for allowing automation on the smoke tests and future integration
tests, so we can check different scenarios as we need.

The rbac mock load 4 different profiles, so we can choose when
the service is started which profile to use by setting the
`APP_CLIENTS_RBAC_PROFILE` variable. The profiles that we can
use are:
- `super-admin`: only contain `*:*:*`
- `domain-admin`: contain the default Domain Admin permissions.
- `domain-readonly`: contain the default Domain Read permissions.
- `domain-no-perms`: does not contain any permission at all.
- `custom`: this is provided for development reasons, if we
  want to customize the permissions. We need to edit the
  `internal/infrastructure/service/impl/mock/rbac/impl/custom.yaml`
  file, build and restart the service again.

## Using the rbac mock

The RBAC mock is started when all the below is true:
- ENV_NAME=local
- APP_ENABLE_RBAC=true or it is enabled at the `configs/config.yaml` file.
- CLIENTS_RBAC_BASE_URL is not empty or it is defined at the `configs/config.yaml` file.

By default, the profile that is used by the rbac mock
is `domain-admin`.

## Using on the smoke and integration tests

- It is automatically started with the suite test.
- By default the profile used is the Domain Administrator.
- The suite has RbacMock attribute that we can use to
  set the permissions the mock will return, so we can
  simulate a different set of permissions for our
  service, when required. We can override the profile by
  the below:

```golang
	s.RbacMock.SetPermissions(mock_rbac.Profiles[mock_rbac.ProfileDomainReadOnly])
```

## Debugging from vscode

- Edit your `.vscode/launch.json` file, and add the below into the
  `"env"` section:

  ```
  "APP_CLIENTS_RBAC_PROFILE": "domain-readonly",
  ```

- You can use the predefined profiles `domain-admin`, `domain-readonly`,
  `domain-none`, `super-admin` and `custom` profiles.
- When you are using the `custom` profile, you could want
  to modify the list of permissions at the
  `internal/infrastructure/service/impl/mock/rbac/impl/custom.yaml`
  file.
- Once you are ready, launch the debugger, set your breakpoints
  at `internal/infrastructure/middleware/rbac.go` file, and
  debug your code step by step.
- To raise the breakpoint, some of the `./test/scripts/local-*`
  files could be useful.
