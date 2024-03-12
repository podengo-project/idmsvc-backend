package rbac

import (
	"context"
)

type Rbac interface {
	IsAllowed(ctx context.Context, permission string) (bool, error)
}
