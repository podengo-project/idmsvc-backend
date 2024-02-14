package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBeforeCreateHook(t *testing.T) {
	item := &HostconfJwk{}
	assert.Empty(t, item.CreatedAt)
	assert.Empty(t, item.UpdatedAt)

	item.BeforeCreate(nil)
	assert.NotEmpty(t, item.CreatedAt)
	assert.NotEmpty(t, item.UpdatedAt)
}
