package router

import (
	_ "embed"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/podengo-project/idmsvc-backend/internal/api/openapi"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/handler"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/middleware"
	rbac_data "github.com/podengo-project/idmsvc-backend/internal/infrastructure/middleware/rbac-data"
	"github.com/podengo-project/idmsvc-backend/internal/metrics"
	"github.com/podengo-project/idmsvc-backend/internal/usecase/client/rbac"
)

// skipperValidate is an alias to represent skipper for API validation middleware
// FIXME Once the openapi specification is propery defined, remove this skipper
// to validate every request
var skipperValidate echo_middleware.Skipper = skipperUserPredicate

type enforceRoute struct {
	Method string
	Path   string
}

var userEnforceRoutes = []enforceRoute{
	{"POST", "/api/idmsvc/v1/domains/token"},
	{"GET", "/api/idmsvc/v1/domains"},
	{"PATCH", "/api/idmsvc/v1/domains/:uuid"},
	{"DELETE", "/api/idmsvc/v1/domains/:uuid"},
}

var systemEnforceRoutes = []enforceRoute{
	{"POST", "/api/idmsvc/v1/domains"},
	{"PUT", "/api/idmsvc/v1/domains/:uuid"},
	{"POST", "/api/idmsvc/v1/host-conf/:inventory_id/:fqdn"},
	{"GET", "/api/idmsvc/v1/signing_keys"},
}

var mixedEnforceRoutes = []enforceRoute{
	{"GET", "/api/idmsvc/v1/domains/:uuid"},
}

//go:embed rbac.yaml
var rbacMapBytes []byte

func getOpenapiPaths(cfg *config.Config, version string) func() []string {
	if cfg == nil {
		panic("'cfg' is nil")
	}
	if version == "" {
		panic("'version' is an empty string")
	}
	majorVersion := strings.Split(version, ".")[0]
	fullVersion := version
	pathPrefix := trimVersionFromPathPrefix(cfg.Application.PathPrefix)
	cachedPaths := []string{
		fmt.Sprintf("%s/v%s/openapi.json", pathPrefix, fullVersion),
		fmt.Sprintf("%s/v%s/openapi.json", pathPrefix, majorVersion),
	}
	return func() []string {
		return cachedPaths
	}
}

func newRbacSkipper(service string) echo_middleware.Skipper {
	if service == "" {
		panic("service is an empty string")
	}
	return func(c echo.Context) bool {
		var (
			cc middleware.DomainContextInterface
			ok bool
		)
		openapiPath := "/api/" + service + "/v1/openapi.json"
		ctx := c.Request().Context()
		routePath := c.Path()
		// The access to the openapi specification is public
		if routePath == openapiPath {
			slog.DebugContext(ctx, "route is '"+openapiPath+"'")
			return true
		}
		if cc, ok = c.(middleware.DomainContextInterface); !ok {
			slog.WarnContext(ctx, "'c' is not a DomainContextInterface")
			return false
		}
		if cc.XRHID().Identity.Type == "System" {
			return true
		}
		return false
	}
}

func initRbacMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	if !cfg.Application.EnableRBAC {
		return middleware.DefaultNooperation
	}

	service := rbac_data.RBACService("idmsvc")
	rbac_data.SetRbacServiceValidator(rbac_data.NewRbacServiceValidator(service))
	rbac_data.SetRbacResourceValidator(
		rbac_data.NewRbacResourceValidator(
			rbac_data.RBACResource("token"),
			rbac_data.RBACResource("domains"),
			rbac_data.RBACResource("host_conf"),
			rbac_data.RBACResource("signing_keys"),
		),
	)
	prefix, rbacMap := rbac_data.RBACMapLoad(rbacMapBytes)
	base := strings.TrimSuffix(cfg.Clients.RbacBaseURL, "/")
	client, err := rbac.NewClientWithResponses(base)
	if err != nil {
		panic(fmt.Errorf("error creating rbac client: %w", err))
	}
	rbacMiddleware := middleware.RBACWithConfig(
		&middleware.RBACConfig{
			Skipper:       newRbacSkipper("idmsvc"),
			Prefix:        prefix,
			PermissionMap: rbacMap,
			Client:        rbac.New("idmsvc", client),
		},
	)
	return rbacMiddleware
}

func guardNewGroupPublic(e *echo.Group, cfg *config.Config, app handler.Application, metrics *metrics.Metrics) {
	if e == nil {
		panic("'e' is nil")
	}
	if cfg == nil {
		panic("'cfg' is nil")
	}
	if app == nil {
		panic("'app' is nil")
	}
	if metrics == nil {
		panic("'metrics' is nil")
	}
}

func newGroupPublic(e *echo.Group, cfg *config.Config, app handler.Application, metrics *metrics.Metrics) *echo.Group {
	guardNewGroupPublic(e, cfg, app, metrics)
	// Initialize middlewares
	fakeIdentityMiddleware := middleware.DefaultNooperation
	if cfg.Application.AcceptXRHFakeIdentity {
		fakeIdentityMiddleware = middleware.FakeIdentityWithConfig(
			&middleware.FakeIdentityConfig{
				Skipper: skipperSystemPredicate,
			},
		)
	}

	parseXRHIDMiddleware := middleware.ParseXRHIDMiddlewareWithConfig(
		&middleware.ParseXRHIDMiddlewareConfig{},
	)

	mixedIdentityMiddleware := middleware.EnforceIdentityWithConfig(
		&middleware.IdentityConfig{
			Skipper: skipperMixedPredicate,
			Predicates: []middleware.IdentityPredicateEntry{
				{
					Name: "mixed-identity",
					Predicate: middleware.NewEnforceOr(
						middleware.EnforceSystemPredicate,
						middleware.EnforceUserPredicate,
						middleware.EnforceServiceAccountPredicate,
					),
				},
			},
		},
	)
	systemIdentityMiddleware := middleware.EnforceIdentityWithConfig(
		&middleware.IdentityConfig{
			Skipper: skipperSystemPredicate,
			Predicates: []middleware.IdentityPredicateEntry{
				{
					Name:      "system-identity",
					Predicate: middleware.EnforceSystemPredicate,
				},
			},
		},
	)
	userAndSAIdentityMiddleware := middleware.EnforceIdentityWithConfig(
		&middleware.IdentityConfig{
			Skipper: skipperUserPredicate,
			Predicates: []middleware.IdentityPredicateEntry{
				{
					Name: "user-and-sa-identity",
					Predicate: middleware.NewEnforceOr(
						middleware.EnforceUserPredicate,
						middleware.EnforceServiceAccountPredicate,
					),
				},
			},
		},
	)

	// FIXME Refactor to inject the config.Config dependency
	rbacMiddleware := initRbacMiddleware(cfg)
	bodyLimit := echo_middleware.BodyLimit(strconv.Itoa(cfg.Application.SizeLimitRequestBody))

	metricsMiddleware := middleware.MetricsMiddlewareWithConfig(
		&middleware.MetricsConfig{
			Metrics: metrics,
		},
	)
	validateAPI := middleware.DefaultNooperation
	if cfg.Application.ValidateAPI {
		middleware.InitOpenAPIFormats()
		validateAPI = middleware.RequestResponseValidatorWithConfig(
			// FIXME Get the values from the application config
			&middleware.RequestResponseValidatorConfig{
				Skipper:          nil,
				ValidateRequest:  true,
				ValidateResponse: false,
			},
		)
	}

	// Wire the middlewares
	e.Use(
		middleware.CreateContext(),
		metricsMiddleware,
		fakeIdentityMiddleware,
		parseXRHIDMiddleware,
		mixedIdentityMiddleware,
		systemIdentityMiddleware,
		userAndSAIdentityMiddleware,
		rbacMiddleware,
		echo_middleware.Secure(),
		// TODO Check if this is made by 3scale
		// middleware.CORSWithConfig(middleware.CORSConfig{}),

		// bodyLimit should be before validateAPI, because validateAPI
		// will read the whole request body to validate the input.
		bodyLimit,
		validateAPI,
	)

	// Setup routes
	public.RegisterHandlersWithBaseURL(e, app, "")
	openapi.RegisterHandlersWithBaseURL(e, app, "")
	return e
}

// skipperUserPredicate applied when using EnforceUserPredicate.
// ctx is the request context.
// Return true if enforce identity is skipped, else false.
func skipperUserPredicate(ctx echo.Context) bool {
	var r enforceRoute
	path := ctx.Path()
	method := ctx.Request().Method
	// it is not expected a big number of routes, but if that were
	// the case into the future, it is more efficient to check
	// directly against a hashmap instead of traversing the slice
	for i := range userEnforceRoutes {
		r = userEnforceRoutes[i]
		if method == r.Method && path == r.Path {
			return false
		}

	}
	return true
}

// skipperSystemPredicate applied when using EnforceSystemPredicate.
// ctx is the request context.
// Return true if enforce identity is skipped, else false.
func skipperSystemPredicate(ctx echo.Context) bool {
	// Read the route path __pattern__ that matched this request
	var r enforceRoute
	path := ctx.Path()
	method := ctx.Request().Method
	// it is not expected a big number of routes, but if that were
	// the case into the future, it is more efficient to check
	// directly against a hashmap instead of traversing the slice
	for i := range systemEnforceRoutes {
		r = systemEnforceRoutes[i]
		if method == r.Method && path == r.Path {
			return false
		}
	}
	return true
}

// skipperMixedPredicate applied for specific routes.
// ctx is the request context.
// Return true if enforce identity is skipped, else false.
func skipperMixedPredicate(ctx echo.Context) bool {
	// Read the route path __pattern__ that matched this request
	var r enforceRoute
	path := ctx.Path()
	method := ctx.Request().Method
	// it is not expected a big number of routes, but if that were
	// the case into the future, it is more efficient to check
	// directly against a hashmap instead of traversing the slice
	for i := range mixedEnforceRoutes {
		r = mixedEnforceRoutes[i]
		if method == r.Method && path == r.Path {
			return false
		}
	}
	return true
}

// newSkipperOpenapi skip /api/idmsvc/v*/openapi.json path
func newSkipperOpenapi(cfg *config.Config, version string) echo_middleware.Skipper {
	paths := getOpenapiPaths(cfg, version)()
	return func(ctx echo.Context) bool {
		route := ctx.Path()
		for i := range paths {
			if paths[i] == route {
				return true
			}
		}
		return false
	}
}
