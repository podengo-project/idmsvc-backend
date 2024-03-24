package rbac

import (
	"context"
)

type Rbac interface {
	IsAllowed(ctx context.Context, xrhid, permission string) (bool, error)
}
