package impl

import (
	"context"
	"testing"
	"time"

	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func helperConfig() *config.Config {
	cfg := &config.Config{}
	config.Load(cfg)

	// Override config values
	cfg.Application.EnableRBAC = true
	if cfg.Clients.RbacBaseURL == "" {
		panic("when rbac is enabled, set your 'clients.rbac_base_url' at your 'configs/config.yaml' file or CLIENTS_RBAC_BASE_URL variable to override")
	}
	return cfg
}

func TestNewMockRbacGuards(t *testing.T) {
	assert.PanicsWithValue(t, "ctx is nil", func() {
		newRbacMockGuards(nil, nil)
	})

	ctx := context.Background()
	assert.PanicsWithValue(t, "cfg is nil", func() {
		newRbacMockGuards(ctx, nil)
	})

	cfg := &config.Config{}
	assert.PanicsWithValue(t, "Config.Application.EnableRBAC is false", func() {
		newRbacMockGuards(ctx, cfg)
	})

	cfg.Application.EnableRBAC = true
	assert.PanicsWithValue(t, "Config.Clients.RbacBaseURL is an empty string", func() {
		newRbacMockGuards(ctx, cfg)
	})

}

func TestNewMockRbac(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cfg := &config.Config{}
	cfg.Application.EnableRBAC = true
	cfg.Clients.RbacBaseURL = "\n"
	assert.PanicsWithValue(t, "error parsing rbac client url: parse \"\\n\": net/url: invalid control character in URL", func() {
		NewRbacMock(ctx, cfg)
	})
	cancel()

	// Success scenario
	ctx, cancel = context.WithCancel(context.Background())
	cfg = helperConfig()
	assert.NotPanics(t, func() {
		NewRbacMock(ctx, cfg)
	})
	cancel()
}

func TestStartStop(t *testing.T) {
	var err error
	cfg := helperConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv, _ := NewRbacMock(ctx, cfg)
	assert.NotPanics(t, func() {
		err = srv.Start()
		if err == nil {
			defer srv.Stop()
		}
	})
	require.NoError(t, err)
}

func TestGetBaseURL(t *testing.T) {
	var err error
	cfg := helperConfig()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv, mock := NewRbacMock(ctx, cfg)
	require.NotNil(t, srv)
	require.NotNil(t, mock)
	err = srv.Start()
	defer srv.Stop()
	require.NoError(t, err)
	err = mock.WaitAddress(3 * time.Second)
	require.NoError(t, err)
	assert.NotEqual(t, "", mock.GetBaseURL())
}

func isStringInSlice(str string, slice []string) bool {
	if slice == nil {
		return false
	}
	for i := range slice {
		if str == slice[i] {
			return true
		}
	}
	return false
}

func EqualStringSlices(t *testing.T, expected, current []string) bool {
	for i := range expected {
		if !isStringInSlice(expected[i], current) {
			t.Logf("expected '%s' string not found into the current slice", expected[i])
			return false
		}
	}

	for i := range current {
		if !isStringInSlice(current[i], expected) {
			t.Logf("current '%s' string not found into the expected slice", current[i])
			return false
		}
	}

	return true
}

func TestLoadProfile(t *testing.T) {
	assert.PanicsWithValue(t, "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `--` into []string", func() {
		LoadProfile([]byte(`--`))
	}, "Panic on unmarshalling the yaml data")

	LoadProfile([]byte(`---
- idmsvc:resource:read
`))
}

func TestProfiles(t *testing.T) {
	profiles := []string{
		ProfileSuperAdmin,
		ProfileDomainAdmin,
		ProfileDomainReadOnly,
		ProfileDomainNoPerms,
		ProfileCustom,
	}

	// forward reference
	for _, profile := range profiles {
		_, ok := Profiles[profile]
		require.True(t, ok)
	}

	// reverse reference
	for profile := range Profiles {
		require.True(t, isStringInSlice(profile, profiles))
	}
}
