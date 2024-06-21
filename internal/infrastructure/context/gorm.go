package context

import (
	"context"

	"gorm.io/gorm"
)

type keyDB string

// CtxWithDB create a context that contain the specified
// *gorm.DB value.
func CtxWithDB(ctx context.Context, db *gorm.DB) context.Context {
	if ctx == nil {
		panic("'ctx' is nil")
	}
	if db == nil {
		panic("'db' is nil")
	}
	key := keyDB("db")
	return context.WithValue(ctx, key, db)
}

// DBFromCtx get a db from a specified context.
func DBFromCtx(ctx context.Context) *gorm.DB {
	if ctx == nil {
		panic("'ctx' is nil")
	}
	key := keyDB("db")
	db, ok := ctx.Value(key).(*gorm.DB)
	if !ok {
		panic("'db' could not be read")
	}
	return db
}
