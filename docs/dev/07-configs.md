# Adding application configuration

## Updating deployment

If we want to govern the value on deployment time, we need to edit the
deployment file `deployment/clowder.yaml` file, and add the parameter `APP_FOO`
to the OpenShift template resource.

```yaml
parameters:
  - name: APP_FOO
    # Set the feault value, if any
    value: "100"
    description: |
      My Foo parameter
```

We need to inject this value in the form of environment variable into the pods
that will need this parameter, so in the same file, add the parameter as an
environment variable:

```yaml
---
apiVersion: template.openshift.io/v1
kind: Template
metadata:
  name: idmsvc

objects:
  - apiVersion: cloud.redhat.com/v1alpha1
    kind: ClowdApp
    metadata:
      name: ${APP_NAME}-backend
    spec:
      deployments:
        - name: service
          podSpec:
            env:
              - name: APP_FOO
                value: ${APP_FOO}
```

For string values:

```yaml
---
apiVersion: template.openshift.io/v1
kind: Template
metadata:
  name: hmsidm

objects:
  - apiVersion: cloud.redhat.com/v1alpha1
    kind: ClowdApp
    metadata:
      name: ${APP_NAME}-backend
    spec:
      deployments:
        - name: service
          podSpec:
            env:
              - name: APP_FOO
                value: "${{APP_FOO}}"
```

## Adding the new parameter to the service

You will need to edit the file `internal/config/config.go` and add to the
structure `Application` the values you need:

```golang
// Application hold specific application settings
type Application struct {
    // My foo parameter
    Foo int `mapstructure:"foo"`
}
```

and set default value at `setDefaults` function.

## Update config.example.yaml

Update the template for the local configuration, so that other developers are
aware of the new parameter, and they can customize when they require for its
local environment. Edit `configs/config.example.yaml`, and your current
`configs/config.yaml` files.

## Use from the configuration

Any setting is used from the `config.Config` structure, and a reference
to the structure is injected wherever it is required to be read the
values.
