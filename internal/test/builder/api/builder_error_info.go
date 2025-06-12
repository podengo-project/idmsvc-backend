package api

import (
	"net/http"
	"strconv"

	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"go.openly.dev/pointy"
)

type ErrorInfo interface {
	Build() *public.ErrorInfo
	WithCode(value string) ErrorInfo
	WithDetail(value string) ErrorInfo
	WithId(value string) ErrorInfo
	WithStatus(value int) ErrorInfo
	WithTitle(value string) ErrorInfo
}

type errorInfo public.ErrorInfo

// NewErrorInfo create new ErrorInfo builder
func NewErrorInfo(statusCode int) ErrorInfo {
	return &errorInfo{
		Code:   nil,
		Detail: nil,
		Id:     "",
		Status: strconv.Itoa(statusCode),
		Title:  http.StatusText(statusCode),
	}
}

func (b *errorInfo) Build() *public.ErrorInfo {
	return (*public.ErrorInfo)(b)
}

func (b *errorInfo) WithCode(value string) ErrorInfo {
	if value == "" {
		b.Code = nil
	} else {
		b.Code = pointy.String(value)
	}
	return b
}

func (b *errorInfo) WithDetail(value string) ErrorInfo {
	if value == "" {
		b.Detail = nil
	} else {
		b.Detail = pointy.String(value)
	}
	return b
}

func (b *errorInfo) WithId(value string) ErrorInfo {
	b.Id = value
	return b
}

func (b *errorInfo) WithStatus(value int) ErrorInfo {
	b.Status = strconv.Itoa(value)
	return b
}

func (b *errorInfo) WithTitle(value string) ErrorInfo {
	b.Title = value
	return b
}
