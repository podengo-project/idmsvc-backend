package logger

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSlogMetaHandler(t *testing.T) {
	var mh SlogMetaHandler

	assert.NotPanics(t, func() {
		mh = NewSlogMetaHandler()
	})

	assert.NotNil(t, mh)
}

func TestTestNewSlogMetaHandler_Add(t *testing.T) {
	mh := &slogMetaHandler{
		handlers: []slog.Handler{},
	}
	h := slog.NewTextHandler(
		os.Stderr,
		nil,
	)
	mh.Add(h)

	assert.Equal(t, 1, len(mh.handlers))
}

func TestTestNewSlogMetaHandler_Enabled(t *testing.T) {
	mh := &slogMetaHandler{
		handlers: []slog.Handler{},
	}
	h := slog.NewTextHandler(
		os.Stderr,
		nil,
	)
	mh.Add(h)

	ctx := context.Background()
	assert.NotNil(t, ctx)

	assert.True(t, mh.Enabled(ctx, slog.LevelError.Level()))
	assert.False(t, mh.Enabled(ctx, slog.LevelDebug.Level()))
}

type stringWriter struct {
	buffer *bytes.Buffer
}

func (s *stringWriter) Write(p []byte) (n int, err error) {
	n, err = s.buffer.Write(p)
	if err != nil {
		return n, err
	}
	return len(p), nil
}

func (s *stringWriter) String() string {
	return s.buffer.String()
}

func slogMetaHandlerCreate(mh *slogMetaHandler, strWriter ...*stringWriter) {
	opts := &slog.HandlerOptions{
		AddSource: false,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Do not print date/time or we can't verfiy
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue("NODATE")
			}

			return a
		},
	}

	for i := range strWriter {
		h1 := slog.NewTextHandler(
			strWriter[i],
			opts,
		)
		mh.Add(h1)
	}
}

func TestTestNewSlogMetaHandler_WithAttrs(t *testing.T) {
	mh := &slogMetaHandler{
		handlers: []slog.Handler{},
	}

	w := &stringWriter{buffer: &bytes.Buffer{}}
	slogMetaHandlerCreate(mh, w)

	logger := slog.New(mh)

	logger = logger.With("x", "y")
	logger.Error("z")

	assert.Equal(t, "time=NODATE level=ERROR msg=z x=y\n", w.String())
}

func TestTestNewSlogMetaHandler_WithGroup(t *testing.T) {
	mh := &slogMetaHandler{
		handlers: []slog.Handler{},
	}

	w := &stringWriter{buffer: &bytes.Buffer{}}
	slogMetaHandlerCreate(mh, w)

	logger := slog.New(mh)
	logger = logger.WithGroup("group")
	logger.Error("z", slog.String("foo", "bar"))

	assert.Equal(t, "time=NODATE level=ERROR msg=z group.foo=bar\n", w.String())
}

func TestTestNewSlogMetaHandler_Log(t *testing.T) {
	mh := &slogMetaHandler{
		handlers: []slog.Handler{},
	}

	w1 := &stringWriter{buffer: &bytes.Buffer{}}
	w2 := &stringWriter{buffer: &bytes.Buffer{}}
	slogMetaHandlerCreate(mh, w1, w2)

	logger := slog.New(mh)
	logger.Error("boom", slog.String("bang", "crash"))

	expected_log := "time=NODATE level=ERROR msg=boom bang=crash\n"
	assert.Equal(t, expected_log, w1.String())
	assert.Equal(t, expected_log, w2.String())
}
