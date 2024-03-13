package rbac_data

// Example
//
// const service = RBACService("myservice")
// const resource = RBACResource("domains")
// SetRbacPermissionValidators(service, []string{
//   "domains", "tokens", "host-conf",
// })
// rbacMap, err := NewRBACMap().
// 	Add("/domains", http.MethodGet, NewRbacPermission(service, resource, RbacVerbRead)).
// 	Build()
// if err != nil {
// 	panic(err)
// }

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

var (
	rbacServiceValidator  RbacServiceValidator  = DefaultRbacServiceValidator
	rbacResourceValidator RbacResourceValidator = DefaultRbacResourceValidator
	rbacVerbValidator     RbacVerbValidator     = DefaultRbacVerbValidator
)

// RBACService
type RBACService string

const (
	RBACServiceUndefined RBACService = ""
	RBACServiceAny       RBACService = "*"
)

type RbacServiceValidator func(s RBACService) RBACService

var serviceValidator = DefaultRbacServiceValidator

// SetRbacServiceValidator set the service validator function that will be
// called when loading the rbac mapping.
// v is the validator function.
func SetRbacServiceValidator(v RbacServiceValidator) {
	if v == nil {
		panic("service validator is nil")
	}
	serviceValidator = v
}

func NewRbacServiceValidator(myService RBACService) RbacServiceValidator {
	return func(s RBACService) RBACService {
		DefaultRbacServiceValidator(s)
		if s != myService {
			panic(fmt.Sprintf("Service must be '%s'", myService))
		}
		return s
	}
}

func DefaultRbacServiceValidator(s RBACService) RBACService {
	if s == RBACServiceUndefined {
		panic("service is undefined")
	}
	return s
}

// RBACResource
type RBACResource string

const (
	RBACResourceUndefined RBACResource = ""
	RBACResourceAny       RBACResource = "*"
)

// RbacResourceValidator validate the provided resource and panic
// if something is wrong, else return the provided value.
type RbacResourceValidator func(r RBACResource) RBACResource

// NewRbacResourceValidator return a handler to manage this scenario
func NewRbacResourceValidator(resources ...RBACResource) RbacResourceValidator {
	resourceMap := make(map[RBACResource]bool, len(resources))
	for i := range resources {
		resourceMap[resources[i]] = true
	}
	return func(resource RBACResource) RBACResource {
		if ok := resourceMap[resource]; !ok {
			panic(fmt.Sprintf("resource '%s' is not valid", string(resource)))
		}
		return resource
	}
}

func DefaultRbacResourceValidator(r RBACResource) RBACResource {
	if r == RBACResourceUndefined {
		panic("resource is undefined")
	}
	return r
}

func SetRbacResourceValidator(h RbacResourceValidator) {
	if h == nil {
		panic("resource validator is nil")
	}
	rbacResourceValidator = h
}

type RBACVerb string

// The following constants result from the schema below
// https://github.com/RedHatInsights/rbac-config/blob/master/schemas/permissions.schema
const (
	RbacVerbUndefined RBACVerb = ""
	RbacVerbAny       RBACVerb = "*"
	RbacVerbRead      RBACVerb = "read"
	RbacVerbWrite     RBACVerb = "write"
	RbacVerbCreate    RBACVerb = "create"
	RbacVerbUpdate    RBACVerb = "update"
	RbacVerbDelete    RBACVerb = "delete"
	RbacVerbLink      RBACVerb = "link"
	RbacVerbUnlink    RBACVerb = "unlink"
	RbacVerbOrder     RBACVerb = "order"
	RbacVerbExecute   RBACVerb = "execute"
)

type RbacVerbValidator func(v RBACVerb) RBACVerb

func DefaultRbacVerbValidator(v RBACVerb) RBACVerb {
	switch v {
	case RbacVerbUndefined:
		panic("Verb is undefined")
	case RbacVerbAny:
	case RbacVerbRead:
	case RbacVerbWrite:
	case RbacVerbCreate:
	case RbacVerbUpdate:
	case RbacVerbDelete:
	case RbacVerbLink:
	case RbacVerbUnlink:
	case RbacVerbOrder:
	case RbacVerbExecute:
	default:
		panic(fmt.Sprintf("Verb '%s' is not supported", v))
	}
	return v
}

// RBACPermission
type RBACPermission string

// RBACPermissionUndefined represent the undefined permission
const (
	RBACPermissionUndefined = RBACPermission("")
)

// RBACPermissionValidator
type RBACPermissionValidator func(p RBACPermission) RBACPermission

// NewRbacPermission
func NewRbacPermission(service RBACService, resource RBACResource, verb RBACVerb) RBACPermission {
	service = rbacServiceValidator(service)
	resource = rbacResourceValidator(resource)
	verb = rbacVerbValidator(verb)
	p := fmt.Sprintf("%s:%s:%s", service, resource, verb)
	return RBACPermission(p)
}

// RBACMap mapping type
type RBACMap map[Route]map[Method]RBACPermission

func (r *RBACMap) getPermissionGuards(prefix, route, method string) error {
	if prefix == "" {
		return errors.New("prefix is an empty string when it was expected for instance '/api/someservice/v1'")
	}
	if route == "" {
		return errors.New("route is an empty string")
	}
	if method == "" {
		return errors.New("method is an empty string")
	}
	return nil
}

// GetPermission
func (r *RBACMap) GetPermission(prefix, route, method string) (RBACPermission, error) {
	var (
		err        error
		methods    map[Method]RBACPermission
		permission RBACPermission
		ok         bool
	)
	if err = r.getPermissionGuards(prefix, route, method); err != nil {
		return "", echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	// FIXME HMS-3522 It can receive /api/idmsvc/v1 or /api/idmsvc/v1.0
	route = strings.TrimPrefix(route, prefix)
	if methods, ok = (*r)[Route(route)]; !ok {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Route not authorized")
	}
	if permission, ok = methods[Method(method)]; !ok {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Method not authorized")
	}
	return permission, nil
}
