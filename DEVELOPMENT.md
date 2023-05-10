# How to use the development environment

Pre-requisites:

* Follow Getting Started steps at `README.md` file.

## Accounts

In order to use the ephemeral development environment you need several accounts.

* Github
* Quay
* RedHat Registry

See also:
<https://consoledot.pages.redhat.com/docs/dev/getting-started/ephemeral/getting-started-with-ees.html#_join_redhatinsights>

### Github

To get access to the ephemeral development environment you need to be a member
of the RedHatInsights organization at Github.

To get your account added to the RedHatInsights organization use the "add user
to github" form at:

<https://source.redhat.com/groups/public/consoledot/consoledot_wiki/console_dot_requests_submission_process#submission-process-for-operational-requests>

### Quay.io

Go to <https://quay.io/> and login with your Red Hat user account.

* Create a new repository.
* Go to <https://quay.io/repository/<rh> user>/<repository>?tab=settings
* Create a new robot account and give the account write permissions to your brand new image repository.
* Create the secrets file: `cp ./scripts/mk/private.example.mk secrets/private.mk`
* Fill out the login and token in `secrets/private.mk`

### Install python dependencies

    make .venv
    source .venv/bin/activate
    pip install -r requirements-dev.txt

## Deploying to ephemeral

    make ephemeral-login
    make ephemeral-namespace-create
    make ephemeral-deploy

    # If correct image was already built and pushed,
    # set EPHEMERAL_NO_BUILD=y to skip these steps.
    make ephemeral-deploy EPHEMERAL_NO_BUILD=y

## How to login to the ephemeral deployment

    make ephemeral-namespace-describe

It should give you a URLs to login to openshift and console.ephemeral:

    Console url: https://console-openshift-console.apps.c-rh-c-eph.8p0c.p1.openshiftapps.com/k8s/cluster/projects/ephemeral-hmup5p
    Frontend route: https://env-ephemeral-hmup5p-z0usxjmg.apps.c-rh-c-eph.8p0c.p1.openshiftapps.com
    Keycloak login: jdoe | Oq8ad9PxpOeTiQLp

The keycloak login is the login for the frontend route.

The keycloak administrator credentials for the keycloak instance deployed for the namespace
can be retrieved by the code snippet below:

    ```sh
    KEYCLOAK_ADMIN="$(base64 -d <<< "$( oc get secret/env-$(oc project -q)-keycloak -o jsonpath='{.data.username}' )" )"
    KEYCLOAK_PASSWORD="$(base64 -d <<< "$( oc get secret/env-$(oc project -q)-keycloak -o jsonpath='{.data.password}' )" )"
    printf "Admin user: %s\nAdmin password: %s\n" "$KEYCLOAK_ADMIN" "$KEYCLOAK_PASSWORD"
    ```

And we could access the keycloak console by the following command:

    ```sh
    xdg-open "https://$(oc get route "$(oc get routes | grep env-$(oc project -q)-auth | awk '{ print $1 }')" -o jsonpath='{.spec.host}')"
    ```

## Launching request against the API

From VSCode or any IDE which support the .http files, we could open
the scripts at `scripts/http/public.http` and launch the http commands
found there.

> You will need yq tool installed to read values in Makefile from `configs/config.yaml` file.
> See: <https://github.com/mikefarah/yq#install>

### Locally

To test our API locally, we can start the service by `make compose-up run` and launching
curl command against it as the below request:

    ```sh
    curl -X GET -H "$( ./scripts/x-rh-identity.sh 12345 jdoe )" http://localhost:8000/api/hmsidm/v1/todo
    ```

### Ephemeral

* Quick in-cluster request:

      ```sh
      oc exec -it "$( oc get pods -l pod=hmsidm-backend-service -o jsonpath='{.items[0].metadata.name}' )" -- curl -H "$( ./scripts/x-rh-identity.sh 12345 jdoe )" "http://localhost:8000/api/hmsidm/v1/todo"
      ```

* Quick out-cluster request:

      ```sh
      USER=jdoe
      PASSWORD="$( base64 -d <<< "$(oc get "secrets/env-$( oc project -q )-keycloak" -o jsonpath='{.data.defaultPassword}' )" )"
      curl -u "$USER:$PASSWORD" "https://$( oc get routes -l app=hmsidm-backend -o jsonpath='{.items[0].spec.host}' )/api/hmsidm/v1/todo"
      ```
