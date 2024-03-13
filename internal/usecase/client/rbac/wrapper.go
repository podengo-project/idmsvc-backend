package rbac

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/podengo-project/idmsvc-backend/internal/interface/client/rbac"
)

type rbacWrapper struct {
	application string
	client      ClientInterface
}

const (
	wildcard   = "*"
	separator  = ":"
	contextKey = "xrhid"
)

// New create a rbac client to check if the required
// permission is allowed.
func New(application string, rbacClient ClientInterface) rbac.Rbac {
	if application == "" {
		panic("application is an empty string")
	}
	if rbacClient == nil {
		panic("rbacClient is nil")
	}
	return &rbacWrapper{
		application: application,
		client:      rbacClient,
	}
}

func (c *rbacWrapper) IsAllowed(ctx context.Context, permission string) (bool, error) {
	var (
		service  string
		resource string
		verb     string
		err      error
		listACL  []string
	)
	if ctx == nil {
		return false, fmt.Errorf("ctx is nil")
	}
	if permission == "" {
		return false, fmt.Errorf("permission to check is an empty string")
	}
	service, resource, verb = c.decomposePermission(permission)
	if listACL, err = c.retrieveACL(ctx); err != nil {
		return false, err
	}
	for i := range listACL {
		if !c.matchPermission(service, resource, verb, listACL[i]) {
			continue
		}
		return true, nil
	}
	return false, fmt.Errorf("permission '%s' not allowed", permission)
}

func (c *rbacWrapper) matchPermission(service, resource, verb, aclItem string) bool {
	sACL, rACL, vACL := c.decomposePermission(aclItem)
	if !c.matchPermissionLabel(service, sACL) {
		return false
	}
	if !c.matchPermissionLabel(resource, rACL) {
		return false
	}
	if !c.matchPermissionLabel(verb, vACL) {
		return false
	}
	return true
}

func (c *rbacWrapper) matchPermissionLabel(label, aclLabel string) bool {
	if label == aclLabel || label == wildcard {
		return true
	}
	return false
}

func (c *rbacWrapper) decomposePermission(permission string) (s, r, v string) {
	tuple := strings.Split(string(permission), separator)
	if len(tuple) != 3 {
		panic("wrong permission tuple")
	}
	s = tuple[0]
	r = tuple[1]
	v = tuple[2]
	return
}

func (c *rbacWrapper) addXRHID(ctx context.Context, req *http.Request) error {
	xrhid := XRHIDRawFromCtx(ctx)
	req.Header.Set("X-Rh-Identity", xrhid)
	return nil
}

func (c *rbacWrapper) retrieveACL(ctx context.Context) ([]string, error) {
	// Credits on RHEnvision: https://github.com/RHEnVision/provisioning-backend/blob/main/internal/clients/http/rbac/rbac_client.go#L83
	// Credits on hmscontent-service:
	var (
		limit  int
		offset int
	)

	const limitDefault = 100

	limit = limitDefault
	offset = 0
	for {
		response, err := c.client.GetPrincipalAccess(
			ctx,
			&GetPrincipalAccessParams{
				Application: c.application,
				Limit:       &limit,
				Offset:      &offset,
			},
			c.addXRHID,
		)
		if err != nil {
			return []string{}, err
		}
		var dataBody []byte
		if _, err = response.Body.Read(dataBody); err != nil {
			return []string{}, err
		}
		var dataACL AccessPagination
		if err = json.Unmarshal(dataBody, &dataACL); err != nil {
			return []string{}, err
		}
		permission := make([]string, offset+len(dataACL.Data))
		for i := range dataACL.Data {
			permission[offset+i] = dataACL.Data[i].Permission
		}
		if *dataACL.Meta.Count == limitDefault {
			offset += limitDefault
			continue
		}
		return permission, nil
	}
}

// XRHIDRawFromCtx read the contextKey entry from
// a previoys created context by ContextWithXRHID
// Return the string with the raw string xrhid or
// a panic happen.
func XRHIDRawFromCtx(ctx context.Context) string {
	data := ctx.Value(contextKey)
	if dataString, ok := data.(string); ok {
		return dataString
	}
	panic("xrhid value is not a string")
}

// ContextWithXRHID create a new context
// Return a new context with the entry contextKey
func ContextWithXRHID(ctx context.Context, xrhidRaw string) context.Context {
	return context.WithValue(ctx, contextKey, xrhidRaw)
}
