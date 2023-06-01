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

func TestCheckFormatIssuerSubject(t *testing.T) {
	var err error
	// Empty string
	err = checkFormatIssuerSubject("")
	assert.EqualError(t, err, "'issuer' format not matching")
	err = checkFormatIssuerSubject("CN=Value?")
	assert.EqualError(t, err, "'issuer' format not matching")
	err = checkFormatIssuerSubject("'")
	assert.EqualError(t, err, "'issuer' format not matching")
	// err = checkFormatIssuerSubject("CN=John Doe, O=Example Corp, OU=Engineering, DC=example, DC=com")
	// assert.NoError(t, err)
}

func TestCheckFormatRealmDomains(t *testing.T) {
	err := checkFormatRealmDomains("")
	assert.NoError(t, err)
}

func TestCheckCertificateFormat(t *testing.T) {
	cert := `-----BEGIN CERTIFICATE-----
MIIF7zCCA9egAwIBAgIUIbYdqhrnlYUEb5g5kbvos+vRCNEwDQYJKoZIhvcNAQEL
BQAwgYYxCzAJBgNVBAYTAlhYMRIwEAYDVQQIDAlTdGF0ZU5hbWUxETAPBgNVBAcM
CENpdHlOYW1lMRQwEgYDVQQKDAtDb21wYW55TmFtZTEbMBkGA1UECwwSQ29tcGFu
eVNlY3Rpb25OYW1lMR0wGwYDVQQDDBRDb21tb25OYW1lT3JIb3N0bmFtZTAeFw0y
MzA2MDEyMjE2MjBaFw0zMzA1MjkyMjE2MjBaMIGGMQswCQYDVQQGEwJYWDESMBAG
A1UECAwJU3RhdGVOYW1lMREwDwYDVQQHDAhDaXR5TmFtZTEUMBIGA1UECgwLQ29t
cGFueU5hbWUxGzAZBgNVBAsMEkNvbXBhbnlTZWN0aW9uTmFtZTEdMBsGA1UEAwwU
Q29tbW9uTmFtZU9ySG9zdG5hbWUwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIK
AoICAQCgYQe0h96QJp2aquJYh2Cat2jd0eBeuwDbhEk4lhPlpP+9wlwS3T/LDka/
KKTeDe/L66m3kxCgReCk8kPX2FV+amHe3LVXo58jaV8uBKXvM7Ud/tD8ZJM5Caj8
fmpiosXavcgmytSqaqRkMAP6hkO0PDg1el/klU5I9ea2cNk13vtbyyNv/oDWztad
pEfcmEy7J55pHfpJSV2Smz0B1c1GPbfMbXiCMtm9LHscBYHIjHPs0S89jXw0664E
AYtQYQqwLSH9AbSdqXXBfq6CfFVhZ3UW7uBHFbhGRi0h9PBtb7d7seqacB5KAdyc
SCiL7XiVaDKnZxoD7ixmXEVMPKs0gWZiUxsz0bWR+yuN8a/tPjcpADEtU9gYbQlX
VZ9zXJWBYS4VHYUfklkxilmXMBrMl39M43uf1/NRjXGTp4TgiVvjCLt/HGivXbFn
KZtTwu4jRnC5paUj9gHqKAIwmxuuiTy953JunbnzOeUK1obeDF+djot2At2FaZh4
Q+Twglgh1CnNp8gamBfRAn0mvEufdccJL/bIxXRVybK5t9GLQfQPblvvWnJ7OqKH
f4Gu6oRziexjjdGQTUXfXIRRHS34R4FawEuoh0BMA/VrFj6XgSYqssJMoePrT1f4
o2ure+2K1/WEWJRP2YjECb+qs1FogSsjHji03qUX4z5Oruww2wIDAQABo1MwUTAd
BgNVHQ4EFgQUAZPj3mcQviqbw5pcFiLYUW3yXzAwHwYDVR0jBBgwFoAUAZPj3mcQ
viqbw5pcFiLYUW3yXzAwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOC
AgEANlScjH9CpOyK7ZGARUWsRncaUPXz3v/wC7sG+RKRY+Pmno/9AOrK0dcwUHeD
wt69H6V2NUVNtVNrJoapABs6KWkoCV2XWYeGh86JuIdxRL2hrXmmoZJ+ODQXDB6k
3CjtctnN2LR6diMNYn7RdM5AQG9xJYSk67tOjOXDDdyIHoYxXfhqiSw4F3quLESs
8kSzKuRWVZtg34wfnJDxb0kl2/Kho+kzafjwZS0jySTKZK5zLhOATr9eDXxSFECe
u4UftGPb9sjdJDFU4X8wLlE2flwucjf+zpRr9ixFyFqjiAIM0Y4ERvUUYWUTFrrG
hMheD0lDxXTnLgAEEgTO0TIxkxIFARe38O2EsYmIlg5zAMxG56o9pO5pIruF9mrS
QiZPBYOsgefGBhD1lAgpU7nGAqlt9k4MG+w77WkhJkECMRUmb3gymTpHCQ0T/PRg
1EyGke9atIlIZTNxtikBfo/eYMF+w0Ut8jJJO/T4VsZ7vgO16PwjaT9P3p5Wh6T3
m1bpFsKXUH3R5IFSXho6kI4BjIju7DauaUQhxItDNs/2FyS6AEX/qmwac1XlPQJr
wD/ueOPpj/mnOu+khYZ0+7Ai2Ay+NX0eiHaTaI9xBxG1ONeKjebnxgxT0dqmNcxn
qYnmleSM7EuIj75kFaeH0KM6+DS4Eo0k3dse/vPSygUaO0E=
-----END CERTIFICATE-----`
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
