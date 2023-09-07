package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateXRhIdentity(t *testing.T) {
	var err error
	v := xrhiAlwaysTrue{}
	require.NotPanics(t, func() {
		err = v.ValidateXRhIdentity(nil)
	})
	assert.NoError(t, err)
}

func TestNewApiServiceValidator(t *testing.T) {
	v := NewApiServiceValidator(nil)
	assert.NotNil(t, v)
}

func TestRequestResponseValidator(t *testing.T) {
	v := RequestResponseValidator()
	assert.NotNil(t, v)
}

func TestCheckFormatIssuer(t *testing.T) {
	var err error
	// Empty string
	err = checkFormatIssuer("")
	assert.EqualError(t, err, "'issuer'='' format not matching")
	err = checkFormatIssuer("CN=Value?")
	assert.EqualError(t, err, "'issuer'='CN=Value?' format not matching")
	err = checkFormatIssuer("'")
	assert.EqualError(t, err, "'issuer'=''' format not matching")
	err = checkFormatIssuer("CN=John Doe, O=Example Corp, OU=Engineering, DC=example, DC=com")
	assert.NoError(t, err)
}

func TestCheckFormatSubject(t *testing.T) {
	var err error
	// Empty string
	err = checkFormatSubject("")
	assert.EqualError(t, err, "'subject'='' format not matching")
	err = checkFormatSubject("CN=Value?")
	assert.EqualError(t, err, "'subject'='CN=Value?' format not matching")
	err = checkFormatSubject("'")
	assert.EqualError(t, err, "'subject'=''' format not matching")
	err = checkFormatSubject("CN=John Doe, O=Example Corp, OU=Engineering, DC=example, DC=com")
	assert.NoError(t, err)
}

func TestCheckFormatRealmDomains(t *testing.T) {
	err := checkFormatRealmDomains("")
	assert.NoError(t, err)
}

func TestCheckCertificateFormat(t *testing.T) {
	cert := `-----BEGIN CERTIFICATE-----
MIIElzCCAv+gAwIBAgIBATANBgkqhkiG9w0BAQsFADA6MRgwFgYDVQQKDA9ITVNJ
RE0tREVWLlRFU1QxHjAcBgNVBAMMFUNlcnRpZmljYXRlIEF1dGhvcml0eTAeFw0y
MzA2MDkxNDA4MThaFw00MzA2MDkxNDA4MThaMDoxGDAWBgNVBAoMD0hNU0lETS1E
RVYuVEVTVDEeMBwGA1UEAwwVQ2VydGlmaWNhdGUgQXV0aG9yaXR5MIIBojANBgkq
hkiG9w0BAQEFAAOCAY8AMIIBigKCAYEAtNQ2nYiwLKUX7NJFQjwhxy4qjZvXa4uj
laPOfkQVqiOSodCsYUhQLwWOHlzn3JWG/kE9Q359tcZhKigA6sCfu3TR0MxNYkmz
O2XcqFz8jfoliX5xaI/WJ85O+R/uT/PQ6BLwBPS3yOGT5zRnmZZ0QuPfERyu15bY
Z7hGpp4NInK85PaGV+7Mjd2NqLBKcqUwbs/PXO1ag6lDuKfShteoh4e5+93Qoziz
nPSkbqnNe3/uFeXvTooxFr3G4GhcO40WrTC9rhLz0bM066RFqk9z/HgS/oVuZINg
bbpMY0klVzuEk2q/mmdgXNw8rgs1blhP2qp/6FQiwuAyupZ1lr4l/v9PZZSbguv5
dF5jEp327ywEF7x37u9qdb84ZAKc/GGJ6oPgwo/3RbPfNtB6n5PpA9pnFfdbBfhF
C85QECIQ1avEvhK+wPTcsrmjuf936Ca/qvhtqgR2UFDy1RtHpdzIG9G4zxcxG0b0
k/IgotoEOnj1tB8qKHkDPXvxKjEpw9Y5AgMBAAGjgacwgaQwHwYDVR0jBBgwFoAU
K/8Bu7M60qz8QwthqG1CSoqfig0wDwYDVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8E
BAMCAcYwHQYDVR0OBBYEFCv/AbuzOtKs/EMLYahtQkqKn4oNMEEGCCsGAQUFBwEB
BDUwMzAxBggrBgEFBQcwAYYlaHR0cDovL2lwYS1jYS5obXNpZG0tZGV2LnRlc3Qv
Y2Evb2NzcDANBgkqhkiG9w0BAQsFAAOCAYEAl4QHnja71OwMtFLouJp13sBgmi8B
vnR/9kPUvk4XWo+I7rouioYWzFfk6bD23WXaZtCA93IEGFXK2V6LrIwbiiqWMdBn
134+QpRqKY7avyV17Pb7W3aDBIm33oFIS8eQ0AVhSe//dz4LwdCNL6TifC+EjMdH
PLBa0d0iqTPp204kDBAPk7Nv4WWVyBTSzhhOzyJeNEYfbYcl+9ZVWtnEz+lEg5LF
kJz71jzd7+CEfiHL9sm1wZ7/+VdvAUmqSItqj/k6TaHh5HYL+f2Uvw2T5ggsNBCw
A+uTHDEe8RWwQJIYmGedn0KAAm9HLnlHdjBX8vSkqHlAMC4BxlK5q57BQK7+gDly
elHya1Shxe6shIT5k9bZkXoPZImRc/dnHGhQINMwtzQl+XQ0OkrUMNooFtzohR2y
An4B5m91EtWFjc1FgxepThC+aTTlcMhX8nwnbIiqgcwaN1k0YM78+CtgOu2r1kKd
nN7+faXYdZgwg8H1N2QeCm4XQ37q4m/Nm3P/
-----END CERTIFICATE-----
`
	err := checkCertificateFormat(cert)
	assert.NoError(t, err)

	cert = `-----BEGIN CERTIFICATE-----
MII..
-----END CERTIFICATE-----`
	err = checkCertificateFormat(cert)
	assert.EqualError(t, err, "Failed to decode CA certificate")
}

func TestInitOpenAPIFormats(t *testing.T) {
	InitOpenAPIFormats()
}
