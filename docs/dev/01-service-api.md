# Service API

## Configurations

All the configuration is centralized at `internal/config/` package.

- Keep in mind that for using environment variables, you need to
  add the field into the set defaults function, otherwise it won't
  be mapped. This is important because the deployment will use the
  values in form of environment variables, however the local environment
  will map them from `configs/config.yaml` file.
- The mapping between the Config structure and the environment variables
  is made automatically by [viper](https://github.com/spf13/viper.git).

## How it is implemented into the POC

- Application has every service to run in the background
  (understanding by service some background process, or
  a http listener).
  - Each service to run in the background implements the
    `ApplicationService` interface (it is not restricted to
    http endpoints).
  - The shutdown is synced by using context.
  - For API services, it configures an echo router and start
    to service the http responses.
  - For the Kafka consumer, start a golang function which subscribe
    to the kafka broker and the necessary topics.
- The http services, configure routing by adding the router
  configuration at: `internal/infrastructure/router/router.go` file.
  - This is the component that set up the configuration for the
    service API.
  - Here we set up middlewares and the http handlers for a method
    and path.
  - echo framework support middleware at global, group and handler
    level:
    - **Global** middlewares impact to all the handlers. They can be
      set for pre-routing (the route of the request was not matched
      yet) or post-routing (the route was matched). On this category
      is important middlewares to trim the trailing slash '/' (pre-routing),
      lecho middleware to write a trace for every request received, or
      request-id middleware.
    - **Group** middlewares impact a group of handlers; private and
      public handlers are organized on this way, so a set of middlewares
      are applied for public handlers (such as enforce X-Rh-Identity).
    - **Handler** middleware impact only on the handler where it is
      indicated. They are not used currently into this Proof of Concept.

API Groups:

- The API groups are splited on **public**, **private** and  **metrics**.
- Each API group has its own openapi specification yaml file, but **metrics**
  which was extracted from the definition of them.

  - **Private** group hold the private api that is not part of our
    service api, and it is not exposed on the API gateway.
    - `/private/healthz` and `/private/livez` are used for the kubernetes
      pod lifecycle; defined at `api/internal.openapi.yaml` file (probably
      `private.openapi.yaml` is a better name for this). TODO in the future
      this would be better to have their custom plaze as the `/metrics`
      endpoint, and leave this place only for inter-cluster communication
      API.
  - **Metrics** `/metrics` will be used for exposing the metrics endpoint, it
      will be defined at `api/metrics.openapi.yaml`.
  - **Public** group hold the public API exposed at `/api/idmsvc/vX` and
      it is defined at `api/public.openapi.yaml`.

For each resource we will find the following different components:

- **Interactor**: the responsibility of this component is to translate
  from API types to model types.
- **Repository**: the responsibility of this component is to deal with
  the data model stored into the database (using gorm for the databse abstraction).
  This basically implements the CRUD operations related with a resource.
  - Every repository method has as first argument the database connecter, the
    decision to pass this as a parameter and not store as a field for the
    specific structure is because to manage transaction level with the database,
    it needs to be managed from the caller to this methods, so it will be the
    application handler which will manage the transation lifecycle.
  - It is made in this way to allow to combine several repositories if it were
    necessary for some specific handler.
- **Presenter**: the responsibility of this component is to translate from
  the resulting data model after execute the business logic to the API response.

> The granularity of the above three components would be per API resource,
> so that if we had `/api/idmsvc/v1/domains`, we will repeat the above 3 components
> for the `/domains` resource.

> TODO Potential refactor in the future to correct the above to be aligned
> with what could be seen at https://github.com/avisiedo/go-microservice-1
> - Presenter is an input/output adapter that depends on the http
>   framework (what we see now as the handler).
> - Interactor represent the business logic which have a free form package.
> - Repositories represent any input/output to/from 3rd party systems,
>   such as database, s3 storage, asynchronous even producer, distributed
>   cache.
>
> The middlewares would be split on presenter and business logic, letting
> the code being reused on different frameworks, by just only changing the
> presenter.

Said the above, we need a `handler` that combines to define the business logic
of the operations; the `handler` component will govern the interaction with the
above components and interactions with third parties as it could be.

### Ey wait, don't you generate the http handler code with the `oapi-generator` tool?

> Ok, let's speak about the generated code.

- The service api generated code is stored at `internal/api/` directory, spliting
  each openapi in a directory.
- For each openapi file it is generated the files below:
  - **types.gen.go**: It holds all the types (parameters and schemas) defined at
    the openapi schema.
  - **server.gen.go**: It holds the code for specific http framework, wrapping
    our handler once the input parameters have been bound to the type structs,
    and call `our application handler` which implement the Service interface
    (this is confusing and I have to change the naming of the services that the
    Application instance instantiate and the Service that generate the
    `oapi-generator` tool, open to suggested name for Application Service
    Interface).
  - **spec.gen.go**: It holds helper code to expose the API openapi specification.
    Expose this api at `/...` (TODO Check the endpoint at console.dot platform) is
    a platform requirement. TODO not implemented on the Proof Of Concept.

- Yes, we generate a handler boilerplate that deal with the input and output
  parameters, but it does not implement the business logic.

### Let's come back to our `<Resource>Service`

- oapi-generator generate a ServerInterface interface per openapi 
  specification.
- Our `ServerInterface` specific implementation is hold at
  `internal/handler/impl` directory.
- We create a file per openapi specification.
- We compose an application interface which is composed by all the 
  ServiceInterface interfaces.
- The code at `internal/handler/impl` currently implements this `big` 
  interface.
  > Having a look to it today, I am not happy with the initial implementation that
  > I did here; a better approach would be:
  > - Define the application interface as it is.
  > - Implement each handler interface in independently.
  > - Implement the `big application interface` as a wrapper
  >   on the specific handlers.

- For ServiceInterface interfaces, split in several files the
  implementation for the main application service api, so that each
  resource has its own file.
  - This will reduce the size of the specific implementation file and
    allow to reduce the conflicts when committing to the repository.

Implement the handler interface (<Resource>Service) at
the body, which basically do (for this poc, but not limited to the below):

- call the interactor
- make the repository or busines logic actions.
- return the reponse using the presenter.


## Design principle followed (it needs more work)

- **S**ingle **R**esponsibility **P**rincipal: This is complex to achieve
  as many times we tend to add additional functionality to the
  components. In essence, we are going to decouple the responsibility
  of the components as much as possible so the unit tests get
  simple to implement.
  When we have to implement the unit test, if it get too complex
  to implement, that use to be an indicator that we are violating
  this principle.
<!-- - **O**pen **C**lose **P**rinciple: Our components are open for
  extension and close for modifications; in golang this means, I can
  create a new version of a component by encapsulating the old one,
  overriding the updated methods, and calling the wrapped old component. 
  Once we have created the new component version, we flag the old one as 
  deprecated to inform the user it has to update the references. We are 
  leaving the old component with no changes. -->
- **L**iskov **S**ubstitution **P**rinciple: Given two components that
  implement the same interface, they can be exchanged with no impact.
  I see this aligned with the OCP, so when we have to create and updated
  version of a component, the new version can substitute the old version
  with no impact at all (in theory).
- **I**nterface **S**egregation **P**rinciple: We split the functionality 
  per interfaces, and we try to create small interface implementations so 
  that we do not push to implement interfaces that we don't need into the 
  specific components. If a component need to combine several interfaces, 
  we can create a new interface which compose several interfaces if it 
  were necessary; in the practice this allow to break down more complex
  interfaces in smaller ones; this is specifally useful when creating
  the unit tests, allowing to test more isolated the SUT (Source Under
  Test).
- **D**ependency **I**nversion **P**rinciple: We depends on interfaces,
  instead of specific implementations. This allow to define the boundaries
  between the components, allow to mock the dependencies so unit tests
  are more isolated, and go do all this by injecting the specific
  components that implements the required interface of the caller. On the
  constructor functions it means we fill a specific structure, but we
  return the interface that implement the data. For testing use to be
  necessary the private methods, to increase the coverage, to avoid
  duplications, we would do something like the below:
  ```golang
  type myType struct {}
  type MyInterface interface {}
  func NewMyInterface() MyInterface {
    return newMyType()
  }
  func newMyInterface() *MyType {
    return &MyType{}
  }
  ```
  This allow us the interface is always accomplished and detected in early
  stage; at the same time allow the private methods to be available
  for the unit tests, and the divergency is minimal.

## I need to add a new middleware

- Add your middleware at: `internal/infrastructure/middleware/`.
- Update the http routers at `internal/infrastructure/router/` according to
  where do you want to add your middleware (api group, global api).

## I need to expose a new port (eg for the metrics)

- Create your openapi specification.
- Add your router at: `internal/infrastructure/router/`. This configure
  the middlewares, and call the Register function that was generated
  by the `oapi-generator` tool for us.
- Implement the Service interface at: `internal/service/impl/`.
- Instantiate the service above at:
  `internal/infrastructure/service/impl/application.go`.
- The current implementations into the PoC could be a guidelines
  to add more endpoints.

> To synchronize and gracefully stop a Service implementation, the context
> is used (`context.WithCancel` is instantiated form the passed context).
>
> The application is setting the signal handler into the `main.go` file
> and it must not be added any other signal handler thant that one,
> because only one handler will receive the signal, which could evoke
> not wished situations when SIGTERM signal is sent by kuberneste to the pod.

### Interface and dependencies

One pain in the neck in golang could be the cycle dependencies. When we get
one of them could be hard to fix them; I have some thoughts about it that
will prevent them.

The basic goal is to depends on interfaces and avoid to depends on specific
implementations.

On scenarios where an interface has several implementations, it seems logic
to add each implementation in an inner package; I mean, if I define
the interface at `internal/myinterface/` directory, and I have 2 implementations
for the interface, it seems logic to add the implementations at:

- `internal/myinterface/implementation1`
- `internal/myinterface/implementation2`

## Adding a new `myres` resource to my public api

- Add the data model at: `internal/domain/model/myres.go`
- Add a builder for the API Rest involved at: `internal/test/builder`
- Generate the migration scripts by: `make build && ./bin/db-tool new myres`
- Fill the migration scripts at: `scripts/db/migrations/`

- Add the `/myres` resource, with its parameters, schemas and its operations
  by using [API Designer](https://console.redhat.com/application-services/api-designer)
  using the file at `api/public.openapi.yaml` file.
- Update the openapi specification by exporting the design above at: `api/`.
- Generate the boilerplate by: `make generate-api`.
- Define the interfaces for the components at:
  - **Interactor** at `internal/interface/interactor/myres_interactor.go`
  - **Repository** at `internal/interface/repository/myres_repository.go`
  - **Presenter** at `internal/interface/presenter/myres_presenter.go`
- Implement the components above at:
  - **Interactor** at `internal/usecase/interactor/myres_interactor.go`
  - **Repository** at `internal/usecase/repository/myres_repository.go`
  - **Presenter** at `internal/usecase/repository/myres_presenter.go`
- Add the components to the application implementation at:
  `internal/handler/impl/application.go`
  ```golang
  type myresComponent struct {
    interactor interactor.MyResInteractor
    repository interactor.MyResRepository
    presenter  interactor.MyResPresenter
  }
  type application struct {
    ...
    myres myresComponent
    ...
  }
  func NewHandler(config *config.Config, db *gorm.DB) handler.Application {
    ...
    iMyRes := usecase_interactor.NewMyResInteractor()
    rMyRes := usecase_repository.NewMyResRepository()
    pMyRes := usecase_presenter.NewMyResPresenter()
    ...
    return &application{
      ...
      myres: myresComponent {
        interactor: iMyRes,
        repository: rMyRes,
        presenter: pMyRes,
      },
      ...
    }
  }
  ```
- Implement the handler at: `internal/handler/impl/myres_handler.go`
  - Be aware it implement methods for the same common object (currently):
  ```golang
  func (a *application)ListMyRes(ctx echo.Context, ...) error {
    // TODO Fill the logic here
    a.myres.interactor.List(...); err != nil {
      return err
    }
    tx := a.db.Begin()
    a.myres.interactor.List(db, ...); err != nil {
      tx.Rollback()
      return err
    }
    tx.Commit()
    a.myres.presenter.List(..., &output); err != nil {
      return err
    }
    return ctx.JSON(http.StatusOK, output)
  }
  ```

## I need to communicate with third party services

- Retrieve the openapi specification.
- Add a new directory at `internal/usecase/client/<service>/`.
- Generate the client proxy from the openapi specification.
- Now just use the proxy client to communicate with the third party service.
- We will need to pass through the following headers:
  - **X-Rh-Identity**: The identity header generated by the api gateway
    by using the JWT (IIRC at cs_jwt header that comes from the external
    source).
  - **X-Insights-Request-Id**: This is necessary for distributing
    tracing; Keep an eye on Open Telemetry for future changes that
    could impact on this header (or additional headers).

## I need to add a new data model or update it

**NOTE** bear in mind that the update process is more complicated, sometimes
requiring two deployments (set all the versions with the new fields and
tables, a second one to remove deprecated tables and fields).

- Do small changes on the model, so the risk is lower, even if this
  require more deployments.
- Add your data model at `internal/infrastructure/domain/model/`.
- Add your up/down migrations `./bin/db-tool migrate new MIGRATION_NAME`.
- Edit your up/down migration scripts at: `scripts/db/migrations/` directory.
- Add your Repository interface for the CRUD operations at:
  `internal/interface/repository/` directory.
- Add your Repository implementation at:
  `internal/usecase/repository/` directory.

> Bear in mind that `compose-up` apply the migrations upto the last current version.

See: https://consoledot.pages.redhat.com/docs/dev/best-practices/db-migrations.html

## Synchronize schema changes with *ipa-hcc* project

The enrollment agent and registration service from *ipa-hcc* project are
sharing the OpenAPI schema with the backend. Shared schema components in
`api/public.openapi.yaml` are marked with the vendor extension `x-rh-ipa-hcc`.
Any time a component is modified, added, or removed, the changes must be
synchronized with `ipa-hcc`. The script `ipa-hcc/contrib/convert_schema.py`
reads the OpenAPI file and generates *ipa-hcc*'s JSON schema files.

The `x-rh-ipa-hcc` is an object with a `type` field and an optional `name`
field:

- `type` must be one of `request`, `response`, or `defs`. The `defs` schemas
  are shared definitions (`$defs`).
- `name` is an optional field. By default, the schema component name is used.

## TODOs

- Add linting for sql migration scripts, and automatic fixer
  (see how hmscontent has this integrated).

## References

- Public APIs: https://console.redhat.com/docs/api
