package errors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func httpErrorFromErr(t *testing.T, err error) (he *echo.HTTPError) {
	assert.Error(t, err)
	he, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	return he
}

func TestNewHTTPError(t *testing.T) {
	internal := errors.New("internal error")
	err := NewHTTPErrorWithInternal(internal, http.StatusForbidden, "forbidden %s!", "resource")
	he := httpErrorFromErr(t, err)
	assert.Equal(t, he.Code, http.StatusForbidden)
	assert.Equal(t, he.Message, "forbidden resource!")
	assert.ErrorIs(t, he.Internal, internal)

	err = NewHTTPErrorWithInternal(internal, http.StatusForbidden, "forbidden")
	he = httpErrorFromErr(t, err)
	assert.Equal(t, he.Code, http.StatusForbidden)
	assert.Equal(t, he.Message, "forbidden")
	assert.ErrorIs(t, he.Internal, internal)
}

func TestNilArgError(t *testing.T) {
	err := NilArgError("param")
	he := httpErrorFromErr(t, err)
	assert.Equal(t, he.Code, http.StatusInternalServerError)
	assert.Equal(t, he.Message, "'param' cannot be nil")
	assert.Nil(t, he.Internal)
}

func TestEmptyArgError(t *testing.T) {
	err := EmptyArgError("param")
	he := httpErrorFromErr(t, err)
	assert.Equal(t, he.Code, int(http.StatusBadRequest))
	require.IsType(t, "string", he.Message)
	msgString, ok := he.Message.(string)
	require.True(t, ok)
	assert.Equal(t, "'param' cannot be empty", msgString)
	assert.NoError(t, he.Internal)
}
