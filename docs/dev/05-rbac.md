# RBAC integration TODO

- The current PoC does not include the integration with
  the RBAC service to provide console.dot platform.

Taking hmscontent implementation at: `https://github.com/content-services/content-sources-backend/`

- It implements a middleware which communicate with rbac service
  and authorize the request.
- It implements a permission mapping which allow to map the
  `method` and `path` of the request supported by the public
  api to the permissions allowed for the service, and afterward
  check if the current x-rh-identity is allowed for the mapped
  permission.

## Understandings

- The Roles and permissions are defined at: `https://github.com/RedHatInsights/rbac-config`
- One Role as a set of permissions.
- Roles are assigned to group of users.
- One user can belong to several group of users.

> The above let to define the set of permissions that can be
> assigned to one user.

- Into EE a keycloak instance is deployed for the namespace.
- We can create new users from keycloak.
- From IQE shell we can create random users.

## GAPS

- How can I create a user which belong to a user group
  with specific permissions (to check service authorization).
  See: https://docs.google.com/document/d/1bzt0jwx7upm0fI4XRp7fOykUvb_7SkfvDYxgKRwEUlc/edit#

## References

- https://github.com/RedHatInsights/rbac-config
- RBAC API documentation: https://console.redhat.com/docs/api/rbac/v1
- Note about RBAC: https://docs.google.com/document/d/1bzt0jwx7upm0fI4XRp7fOykUvb_7SkfvDYxgKRwEUlc/edit#
- Jira ticket about RBAC backend integration into hmscontent:
  https://issues.redhat.com/browse/HMS-470
