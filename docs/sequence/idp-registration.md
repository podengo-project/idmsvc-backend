# Registration process

![Registration Sequence Diagram](idp-registration.svg)

> Simplification: Between the external request and our service
> the X-Rh-Identity and X-Rh-Insights-Request-Id headers
> are aggregated to the request that the micro-service receive,
> authorizing the request.

* (1) The Administrator create a new domain.
  So a `POST /api/idmsvc/v1/domains` request is sent, using
  a user credentials with `hmsidm:domains:write` permission.
* (2) idm-domains-backend request ACL for the current x-rh-identity.
* (3) rbac service return the ACL list:
  * Check the `hmsidm:domains:write` permission exists into the response.
* (4) Create the domain entry for the organization and return the
  resource (here is the only moment the token is returned).
  If RBAC validation fails, return **403 Forbidden**.
* (5) Administrator run the `ipa-hcc register <domain_uuid> <token>`
  as indicated into the UI step indications.
* (6) A `PUT /api/idmsvc/v1/domains/<domain_uuid>/ipa` http
  request is sent to the service, using the RHSM certificate
  assigned to the host by `rhc`, and the `X-Rh-Idm-Registration-Token`
  returned when the domain was created. The `X-Rh-Identity` has the
  cn of the certificate.
* (7) idm-domains-backend request ACL for the current x-rh-identity.
* (8) rbac service return the ACL list.
  * Check permission `hmsidm:domains_ipa:write` exists into the response.
  * Check token and token expiration (for register).
* (9) TODO Request a host inventory request filtering with the cn content
  (check this behavior).
* (10) TODO Response with one item (success) or no items (host does not exist).
* (11) The sub-object
  for the ipa section is returned as response, and the token is
  set to null and its expiration date (remove token).
  If RBAC check fails a **403 Forbidden** response is returned.
  If the token validation failed then a **403 Forbidden** response
  is returned.
* (12) The administrator come back to the UI and request the
  information for the domain by `GET /api/idmsvc/v1/domains/<domain_uuid>`.
* (13) idm-domains-backend request ACL for the current x-rh-identity.
* (14) rbac service return the ACL list.
  * Check the `hmsidm:domain:read` permission exists into the list.
* (15) The domain resource is returned and it contains all the updated information.
  If the RBAC check fails then a **403 Forbidden** response is returned.
  If no domain information is found for the current organization
  a **404 Not Found** response is returned.

---

About permissions:

- Administrator Domain (role), assigned to the Administrator:
  - hmsidm:domains:write
  - hmsidm:domains:read
- Domain Server Agent (role), assigned to the RHSM certificate:
  - hmsidm:domains_ipa:write

## Manual requests - Stage

1. Generate a token at: https://access.stage.redhat.com/management/api
   ```
   OFFLINE_TOKEN="<your offline generated token>"
   ```
2. Get an access token by:
   ```
   ACCESS_TOKEN="$(curl "https://sso.stage.redhat.com/auth/realms/redhat-external/protocol/openid-connect/token" -d grant_type=refresh_token -d client_id=rhsm-api -d refresh_token="$OFFLINE_TOKEN" | jq -r '."access_token"')"
   ```
3. Check the API by:
   ```
   curl -H "Authorization: Bearer ${ACCESS_TOKEN}" "https://console.stage.redhat.com/api/inventory/v1/hosts"
   ```
   > Be aware your HTTPS_PROXY environment variable point out to the right place
4. From the ipa server host, retrieve the CN field by:
   ```
   kinit admin
   CN="$( ipa host-show server.hmsidm-dev.test | grep RHSM | awk -F ":" '{print $2}' | awk -F "," '{print $2}' | awk -F "=" '{print $2}' )"
   ```
5. Launch the request against the insights inventory by:
   ```
   curl -H "Authorization: Bearer ${ACCESS_TOKEN}" "https://console.stage.redhat.com/api/inventory/v1/hosts?filter\[system_profile\]\[owner_id\]=${CN}" | jq 
   ```

```
# Request by fqdn

The below could return more than one record.

curl -s -H "Authorization: Bearer ${ACCESS_TOKEN}" "https://console.stage.redhat.com/api/inventory/v1/hosts?fqdn=server.hmsidm-dev.test"

# Request by the CN of the x-rh-identity
curl -X 'GET' \
  "https://console.stage.redhat.com/api/inventory/v1/hosts?filter\[system_profile\]\[owner_id\]=${CN}" \
  -H 'accept: application/json' \
  -H "Authorization: Bearer ${ACCESS_TOKEN}"
```

> API Documentation at: https://console.stage.redhat.com/docs/api
> Inventory API Documentation at: https://console.stage.redhat.com/docs/api/inventory

### Inventory

**Server**

Once the host is registered with rhc you can find information at: <TODO json file>

**Host VM**

Once the host is registered with rhc you can find information at: <TODO json file>

----

Check the host by:

```
curl -H "Authorization: Bearer ${ACCESS_TOKEN}" "https://console.redhat.com/api/inventory/v1/hosts?registered_with=insights&"
```

## References

- [ipa-hcc repository](https://gitlab.cee.redhat.com/identity-management/idmocp/ipa-hcc).
- [Red Hat Insights API Cheat Sheet](https://developers.redhat.com/cheat-sheets/red-hat-insights-api-cheat-sheet).
- [Red Hat CRC Platform API Documentation](https://console.redhat.com/docs/api).
- [Inventory API Documentation](https://console.redhat.com/docs/api/inventory/v1).
