package impl

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	app_context "github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
	"gorm.io/gorm"
)

func (a *application) GetSigningKeys(ctx echo.Context, params public.GetSigningKeysParams) error {
	var (
		err         error
		tx          *gorm.DB
		keys        []string
		revokedKids []string
		output      *public.SigningKeysResponse
	)
	handlerName := "RegisterDomain"
	logger := app_context.LogFromCtx(ctx.Request().Context())
	if tx = a.db.Begin(); tx.Error != nil {
		logger.Error(errDBTXCommit, slog.String("handler", handlerName))
		return tx.Error
	}
	defer tx.Rollback()

	c := app_context.CtxWithDB(ctx.Request().Context(), tx)
	if keys, revokedKids, err = a.hostconfjwk.repository.GetPublicKeyArray(c); err != nil {
		logger.Error(errDBGeneralError,
			slog.String("handler", handlerName),
		)
		return err
	}

	if tx.Commit(); tx.Error != nil {
		logger.Error(errDBTXCommit, slog.String("handler", handlerName))
		return tx.Error
	}

	if output, err = a.hostconfjwk.presenter.PublicSigningKeys(keys, revokedKids); err != nil {
		logger.Error(errOutputAdapter, slog.String("handler", handlerName))
		return err
	}

	return ctx.JSON(http.StatusOK, output)
}
