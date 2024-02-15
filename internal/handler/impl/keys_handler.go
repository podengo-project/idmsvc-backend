package impl

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"gorm.io/gorm"
)

func (a *application) GetSigningKeys(ctx echo.Context, params public.GetSigningKeysParams) error {
	var (
		err    error
		tx     *gorm.DB
		keys   []string
		output *public.SigningKeysResponse
	)

	if tx = a.db.Begin(); tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	if keys, err = a.hostconfjwk.repository.GetPublicKeyArray(tx); err != nil {
		return err
	}

	if tx.Commit(); tx.Error != nil {
		return tx.Error
	}

	if output, err = a.hostconfjwk.presenter.PublicSigningKeys(keys); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, output)
}
