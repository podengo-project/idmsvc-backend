package assert

import (
	"testing"

	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertDomain asserts that the expected domain with the current domain value without
// consider the domain id field and without consider the gorm.Model data.
// t the current test in progress.
// expected is the expected domain.
// current is the current resulting domain that should match with the expectation.
func AssertDomain(t *testing.T, expected *public.Domain, current *public.Domain) {
	if expected == current {
		return
	}
	require.NotNil(t, expected)
	require.NotNil(t, current)

	assert.Equal(t, expected.DomainId, current.DomainId)
	assert.Equal(t, expected.Title, current.Title)
	assert.Equal(t, expected.Description, current.Description)
	assert.Equal(t, expected.AutoEnrollmentEnabled, current.AutoEnrollmentEnabled)
	require.Equal(t, expected.DomainType, current.DomainType)
	if current.DomainType == "" {
		return
	}
	switch current.DomainType {
	case public.RhelIdm:
		AssertDomainIpa(t, expected.RhelIdm, current.RhelIdm)
	default:
		t.Errorf("body.DomainType has an unexpected value")
	}
}

// AssertDomainIpa asserts the current domain ipa data match the expectation
// without consider the gorm.Model data.
// t represents the current test.
// expected represents the expected domain ipa values.
// current represents the resulting domain ipa values.
func AssertDomainIpa(t *testing.T, expected *public.DomainIpa, current *public.DomainIpa) {
	require.NotNil(t, expected)
	require.NotNil(t, current)

	assert.Equal(t, expected.RealmName, current.RealmName)
	assert.Equal(t, expected.RealmDomains, current.RealmDomains)
	AssertLocations(t, expected.Locations, current.Locations)
	AssertServers(t, expected.Servers, current.Servers)
	AssertCaCerts(t, expected.CaCerts, expected.CaCerts)
	require.NotNil(t, len(expected.Locations), len(current.Locations))
	for i := range expected.Locations {
		assert.Equal(t, expected.Locations[i], current.Locations[i])
	}
	AssertAutomountLocations(t, expected.AutomountLocations, current.AutomountLocations)
}

// AssertLocations asserts the current slice of Locations match
// the slice of expected ones.
// t represents the current test.
// expected represents the expected slice of Location values.
// current represents the resulting slice of Location values.
func AssertLocations(t *testing.T, expected []public.Location, current []public.Location) {
	if expected == nil && nil == current {
		return
	}
	require.Equal(t, len(expected), len(current))
	foundEquals := 0
	for i := range expected {
		for j := range current {
			if expected[i].Name == current[j].Name &&
				expected[i].Description != nil && current[j].Description != nil &&
				*expected[i].Description == *current[j].Description {
				foundEquals++
			}
		}
	}
	assert.Equal(t, foundEquals, len(current))
}

// AssertServers asserts the current slice of public.DomainIpaServer match
// the slice of expected ones.
// t represents the current test.
// expected represents the expected slice of public.DomainIpaServer values.
// current represents the resulting slice of public.DomainIpaServer values.
func AssertServers(t *testing.T, expected []public.DomainIpaServer, current []public.DomainIpaServer) {
	if expected == nil && current == nil {
		return
	}
	require.Equal(t, len(expected), len(current))
	for i := range expected {
		for j := range current {
			if expected[i].Fqdn == current[j].Fqdn {
				assert.Equal(t, expected[i].CaServer, current[j].CaServer)
				assert.Equal(t, expected[i].HccEnrollmentServer, current[j].HccEnrollmentServer)
				assert.Equal(t, expected[i].HccUpdateServer, current[j].HccUpdateServer)
				assert.Equal(t, expected[i].PkinitServer, current[j].PkinitServer)
				assert.Equal(t, expected[i].SubscriptionManagerId, current[j].SubscriptionManagerId)
				assert.Equal(t, expected[i].Location, current[j].Location)
			}
		}
	}
}

// AssertCaCerts asserts the current slice of public.Certificate match
// the slice of expected ones.
// t represents the current test.
// expected represents the expected slice of public.Certificate values.
// current represents the resulting slice of public.Certificate values.
func AssertCaCerts(t *testing.T, expected []public.Certificate, current []public.Certificate) {
	require.NotNil(t, expected)
	require.NotNil(t, current)
	require.Equal(t, len(expected), len(current))

	for i := range expected {
		for j := range current {
			if expected[i].Issuer == current[j].Issuer && expected[i].SerialNumber == current[j].SerialNumber {
				assert.Equal(t, expected[i].Issuer, current[j].Issuer)
				assert.Equal(t, expected[i].Nickname, current[j].Nickname)
				assert.Equal(t, expected[i].NotAfter, current[j].NotAfter)
				assert.Equal(t, expected[i].NotBefore, current[j].NotBefore)
				assert.Equal(t, expected[i].Pem, current[j].Pem)
				assert.Equal(t, expected[i].Subject, current[i].Subject)
			}
		}
	}
}

// AssertAutomountLocations asserts the current slice for automount locations with
// the slice of expected ones.
// t represents the current test.
// expected represents the expected slice of strings or nil.
// current represents the resulting slice of strings or nil.
func AssertAutomountLocations(t *testing.T, expected *[]string, current *[]string) {
	if expected == current {
		return
	}
	require.NotNil(t, expected, current)
	require.Equal(t, len(*expected), len(*current))
	countEquals := 0
	for i := range *expected {
		for j := range *current {
			if (*expected)[i] == (*current)[j] {
				countEquals++
			}
		}
	}
	require.Equal(t, countEquals, len(*current))
}
