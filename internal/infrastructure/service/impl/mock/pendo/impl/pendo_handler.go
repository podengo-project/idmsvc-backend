package impl

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
	"github.com/podengo-project/idmsvc-backend/internal/interface/client/pendo"
)

func (m *mockPendo) guardKindAndGroup(kind pendo.Kind, group pendo.Group) (pendo.Kind, pendo.Group, error) {
	if kind == "" {
		return "", "", echo.NewHTTPError(http.StatusBadRequest, "'kind' is an empty string")
	}
	if kind != pendo.KindVisitor && kind != pendo.KindAccount {
		return "", "", echo.NewHTTPError(http.StatusBadRequest, "'kind' must be 'visitor' or 'account' but it is '%s'", kind)
	}
	if group == "" {
		return "", "", echo.NewHTTPError(http.StatusBadRequest, "'group' is an empty string")
	}
	return kind, group, nil
}

func (m *mockPendo) guardCreateMetadataAccountCustomValue(kind pendo.Kind, group pendo.Group) (pendo.Kind, pendo.Group, error) {
	return m.guardKindAndGroup(kind, group)
}

// CreateMetadataAccountCustomValue implement POST /metadata/:kind/:group/value
// endpoint for the pendo mock.
// See: https://engageapi.pendo.io/?bash%23#bd14400e-0ff9-49e6-932e-61e7a65c6f3c
func (m *mockPendo) CreateMetadataAccountCustomValue(ctx echo.Context) error {
	var (
		kind    pendo.Kind
		group   pendo.Group
		err     error
		metrics pendo.SetMetadataRequest
	)
	if kind, group, err = m.guardCreateMetadataAccountCustomValue(
		pendo.Kind(ctx.Param("kind")),
		pendo.Group(ctx.Param("group")),
	); err != nil {
		kind = kind
		group = group
		return err
	}
	if err := ctx.Bind(metrics); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}
	// TODO We need more knowledge to know how this works
	m.lock.Lock()
	defer m.lock.Unlock()
	m.metrics = metrics
	return ctx.String(http.StatusOK, http.StatusText(http.StatusOK))
}

func (m *mockPendo) guardGetMetadataAccountCustomValue(
	kind pendo.Kind,
	group pendo.Group,
	ID pendo.ID,
	fieldName pendo.FieldName,
) (pendo.Kind, pendo.Group, pendo.ID, pendo.FieldName, error) {
	var err error
	if kind, group, err = m.guardKindAndGroup(kind, group); err != nil {
		return "", "", "", "", err
	}
	if ID == "" {
		return "", "", "", "", echo.NewHTTPError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}
	if fieldName == "" {
		return "", "", "", "", echo.NewHTTPError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}
	return kind, group, ID, fieldName, nil
}

// GetMetadataAccountCustomValue implement GET /metadata/:kind/:group/value/:id/:fieldName
// See: https://engageapi.pendo.io/?bash%23#cb0a05c4-709f-4644-bbd7-18f066b94d7e
func (m *mockPendo) GetMetadataAccountCustomValue(ctx echo.Context) error {
	var (
		kind      pendo.Kind
		group     pendo.Group
		ID        pendo.ID
		fieldName pendo.FieldName
		err       error
		metrics   pendo.SetMetadataRequest
	)
	logger := context.LogFromCtx(ctx.Request().Context())
	if kind, group, ID, fieldName, err = m.guardGetMetadataAccountCustomValue(
		pendo.Kind(ctx.Param("kind")),
		pendo.Group(ctx.Param("group")),
		pendo.ID(ctx.Param("id")),
		pendo.FieldName(ctx.Param("fieldName")),
	); err != nil {
		logger.Error("bad parameter")
		return err
	}
	if err := ctx.Bind(metrics); err != nil {
		logger.Error("bad payload",
			slog.String("kind", string(kind)),
			slog.String("group", string(group)),
			slog.String("id", string(ID)),
			slog.String("fieldName", string(fieldName)),
		)
		return echo.NewHTTPError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}
	// TODO We need more knowledge to know how this works
	m.lock.RLock()
	defer m.lock.RUnlock()
	return ctx.JSON(http.StatusOK, m.metrics)
}
