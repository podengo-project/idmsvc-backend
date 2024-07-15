package impl

import (
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"
)

type queryParams struct {
	Offset      int    `query:"offset"`
	Limit       int    `query:"limit"`
	Application string `query:"application"`
}

func (m *mockRbac) accessHandler(c echo.Context) error {
	const (
		metaCount  = "count"
		metaLimit  = "limit"
		metaOffset = "offset"

		queryOffset      = "offset"
		queryLimit       = "limit"
		queryApplication = "application"

		linksFirst    = "first"
		linksLast     = "last"
		linksPrevious = "previous"
		linksNext     = "next"
	)
	var (
		page   Page
		params queryParams
		err    error
		u      *url.URL
		q      url.Values
	)
	if err = c.Bind(&params); err != nil {
		slog.Error(err.Error())
		return err
	}

	limit := params.Limit
	offset := params.Offset
	count := len(m.data)

	appName := m.appName
	if appName == "" {
		panic("APP_NAME is empty or unset")
	}
	if u, err = url.Parse(m.GetBaseURL() + "/access/"); err != nil {
		return err
	}
	if params.Application != appName {
		q = u.Query()
		q.Set("application", params.Application)
		q.Set("offset", strconv.Itoa(0))
		q.Set("limit", strconv.Itoa(limit))
		u.RawQuery = q.Encode()
		return c.JSON(http.StatusOK, &Page{
			Meta: map[string]any{
				"count":  1,
				"limit":  10,
				"offset": 0,
			},
			Links: map[string]string{
				"First": u.String(),
				"Last":  u.String(),
			},
			Data: []Permission{
				{
					Permission:          "*:*:*",
					ResourceDefinitions: nil,
				},
			},
		})
	}

	if limit <= 0 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	size := min(count-offset, limit)
	if size < 0 {
		size = 0 // prevent underflow
	}

	permissions := make([]Permission, size)

	// Fill meta
	page.Meta = map[string]any{
		metaCount:  count,
		metaLimit:  limit,
		metaOffset: offset,
	}

	// Fill links
	q = u.Query()
	q.Add(queryApplication, "idmsvc")
	q.Add(queryOffset, strconv.Itoa(0))
	q.Add(queryLimit, strconv.Itoa(limit))
	u.RawQuery = q.Encode()
	linkFirst := u.RequestURI()

	q = u.Query()
	newOffset := (count / limit) * limit
	q.Set(queryOffset, strconv.Itoa(newOffset))
	q.Set(queryLimit, strconv.Itoa(limit))
	u.RawQuery = q.Encode()
	linkLast := u.RequestURI()

	// Create the links map.  Add "previous" and "next"
	// conditionally.
	page.Links = map[string]string{
		linksFirst: linkFirst,
		linksLast:  linkLast,
	}

	// "previous" link
	linkPrevious := ""
	if offset > 0 {
		previousOffset := max(offset-limit, 0)
		q = u.Query()
		q.Set(queryOffset, strconv.Itoa(previousOffset))
		q.Set(queryLimit, strconv.Itoa(limit))
		u.RawQuery = q.Encode()
		linkPrevious = u.RequestURI()
		page.Links[linksPrevious] = linkPrevious
	}

	// "next" link
	if offset+limit < count {
		q = u.Query()
		q.Set(queryOffset, strconv.Itoa(offset+limit))
		q.Set(queryLimit, strconv.Itoa(limit))
		u.RawQuery = q.Encode()
		linkNext := u.RequestURI()
		page.Links[linksNext] = linkNext
	}

	// Fill data
	for i := 0; i < size; i++ {
		if offset+i >= len(m.data) {
			break
		}
		permissions[i] = m.data[offset+i]
	}
	page.Data = permissions
	return c.JSON(http.StatusOK, &page)
}
