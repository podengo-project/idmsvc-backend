package impl

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
)

type queryParams struct {
	Offset      int    `query:"offset"`
	Limit       int    `query:"limit"`
	Application string `query:"application"`
}

func (m *mockRbac) authenticate(c echo.Context) error {
	var (
		page   Page
		params queryParams
		err    error
	)
	if err = c.Bind(&params); err != nil {
		slog.Error(err.Error())
		return err
	}

	limit := params.Limit
	offset := params.Offset
	count := len(m.data)

	if limit == 0 {
		limit = 100
	}

	size := offset + limit
	if size > count {
		size = count - offset
	}
	permissions := make([]Permission, size)

	// Fill meta
	page.Meta = map[string]any{
		"count":  count,
		"limit":  limit,
		"offset": offset,
	}

	// Fill links
	var u *url.URL
	var q url.Values
	if u, err = url.Parse(m.GetBaseURL() + "/access/"); err != nil {
		return err
	}
	q = u.Query()
	q.Add("application", "idmsvc")
	q.Add("offset", strconv.Itoa(0))
	q.Add("limit", strconv.Itoa(limit))
	u.RawQuery = q.Encode()
	linkFirst := u.String()

	u.Query().Set("offset", strconv.Itoa(((count+limit-1)/limit)*limit))
	u.Query().Set("limit", strconv.Itoa(limit))
	u.RawQuery = q.Encode()
	linkLast := u.String()

	u.Query().Set("offset", strconv.Itoa(offset-limit))
	u.Query().Set("limit", strconv.Itoa(limit))
	u.RawQuery = q.Encode()
	linkPrevious := u.String()

	u.Query().Set("offset", strconv.Itoa(offset-limit))
	u.Query().Set("limit", strconv.Itoa(limit))
	u.RawQuery = q.Encode()
	linkNext := u.String()

	page.Links = map[string]string{
		"first":    linkFirst,
		"last":     linkLast,
		"previous": linkPrevious,
		"next":     linkNext,
	}
	// Fill data
	if offset < count {
		for i := 0; i < size; i++ {
			if offset+i > count {
				break
			}
			permissions[i] = m.data[i+offset]
		}
	}
	page.Data = permissions
	return c.JSON(http.StatusOK, &page)
}
