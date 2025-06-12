package helper

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pioz/faker"
	"go.openly.dev/pointy"
)

// GenRandNum generate a random number >= min and <= max interval.
// min is the lower boundary interval including the value.
// max is the higher boundary interval including the value.
func GenRandNum(min, max int64) int64 {
	return faker.Int64InRange(min, max+1)
}

// GenPastNearTime generate a past timestamp not further than delta.
// delta is the duration that set the threshold for the time.
// Return a time.Time for the expected interval.
func GenPastNearTime(delta time.Duration) time.Time {
	var value int64 = GenRandNum(0, int64(delta))
	return time.Now().UTC().Add(time.Duration(value) * -1)
}

// GenFutureNearTimeUTC generate a past timestamp not further than delta.
// delta is the duration that set the threshold for the time.
// Return a time.Time for the expected interval.
func GenFutureNearTimeUTC(delta time.Duration) time.Time {
	var value int64 = GenRandNum(0, int64(delta))
	return time.Now().UTC().Add(time.Duration(value))
}

// GenBetweenTimeUTC generate a timestamp between the given parameters
// as earlier as 'begin' and before 'end'.
// delta is the duration that set the threshold for the time.
// Return a time.Time for the expected interval.
func GenBetweenTimeUTC(begin time.Time, end time.Time) time.Time {
	return faker.TimeInRange(begin, end).UTC()
}

// GenRandString generate a random string from the letters set
// with length n.
// letters the set of letters to use.
// n the length of the string.
// Return a random string.
func GenRandString(letters []rune, n int) string {
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[GenRandNum(0, int64(len(letters)-1))]
	}
	return string(s)
}

// GenRandDomainLabel generate a random label according to RFC1035
// Return the string representing a valid domain label.
func GenRandDomainLabel() string {
	// See: https://www.rfc-editor.org/rfc/rfc1035
	//
	// <domain> ::= <subdomain> | " "
	// <subdomain> ::= <label> | <subdomain> "." <label>
	// <label> ::= <letter> [ [ <ldh-str> ] <let-dig> ]
	// <ldh-str> ::= <let-dig-hyp> | <let-dig-hyp> <ldh-str>
	// <let-dig-hyp> ::= <let-dig> | "-"
	// <let-dig> ::= <letter> | <digit>
	// <letter> ::= any one of the 52 alphabetic characters A through Z in
	// upper case and a through z in lower case
	// <digit> ::= any one of the ten digits 0 through 9
	//
	letDigHyp := []rune("abcdefghijklmnopqrstuvwxyz0123456789-")
	letDig := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	letter := []rune("abcdefghijklmnopqrstuvwxyz")
	len := GenRandNum(1, 63)
	if len == 1 {
		return GenRandString(letter, 1)
	}
	if len == 2 {
		return GenRandString(letter, 1) + GenRandString(letDig, 1)
	}
	return GenRandString(letter, 1) + GenRandString(letDigHyp, int(len-2)) + GenRandString(letDig, 1)
}

// GenRandDomainName generate a random domain name for testing
// using .test as defined at RFC 6761
// level is >= 2 and <= 4.
// Return a domain name
func GenRandDomainName(level int) string {
	if level < 2 || level > 4 {
		panic(fmt.Errorf("'level' must be in [2..4] range"))
	}
	labels := make([]string, level)
	for i := 0; i < level-1; i++ {
		label := GenRandDomainLabel()
		labels[i] = label
	}
	labels[level-1] = "test"
	return strings.Join(labels, ".")
}

// GenRandFQDNWithDomain Generate a random FQDN for the given domain.
// domain is the domain that belong the returned FQDN
// Return a FQDN string representation.
func GenRandFQDNWithDomain(domain string) string {
	return strings.Join([]string{GenRandDomainLabel(), domain}, ".")
}

// GenRandFQDN Generate a random FQDN using a random 3 level domain name.
// Return a string value.
func GenRandFQDN() string {
	return GenRandDomainName(3)
}

// GenRandBool generate a random boolean.
// Return a bool value.
func GenRandBool() bool {
	return faker.Bool()
}

// GenRandUserID Generate a random UUID.
// Return a string.
func GenRandUserID() string {
	return uuid.NewString()
}

// GenRandUsername generate a random username.
// Return a string.
func GenRandUsername() string {
	return strings.Join([]string{
		strings.ToLower(GenRandFirstName()),
		strings.ToLower(GenRandLastName()),
	}, ".")
}

// GenRandFirstName generate a random first name.
// Return a string.
func GenRandFirstName() string {
	return faker.FirstName()
}

// GenRandLastName generate a random last name.
// Return a string.
func GenRandLastName() string {
	return faker.LastName()
}

// GenRandEmail generate a random email.
// Return a string.
func GenRandEmail() string {
	return strings.Join([]string{GenRandUsername(), faker.Domain()}, "@")
}

// GenRandPointyBool generate a random bool pointer.
// Return nil, or a bool pointer.
func GenRandPointyBool() *bool {
	if faker.Bool() {
		return pointy.Bool(faker.Bool())
	}
	return nil
}

// GenRandParagraph generate a random paragraph.
// Return a multiline string.
func GenRandParagraph(n int) string {
	if n == 0 {
		n = int(GenRandNum(3, 10))
	}
	return faker.ArticleWithParagraphCount(n)
}

// GenPemCertificate currently return a static PEM certificate.
func GenPemCertificate() string {
	// TODO Implement this generator
	return `-----BEGIN CERTIFICATE-----
MIIOOjCCDSKgAwIBAgIRAIs+mCLMW+1cCl7dwvqQLe4wDQYJKoZIhvcNAQELBQAw
RjELMAkGA1UEBhMCVVMxIjAgBgNVBAoTGUdvb2dsZSBUcnVzdCBTZXJ2aWNlcyBM
TEMxEzARBgNVBAMTCkdUUyBDQSAxQzMwHhcNMjMxMDIzMTExODI0WhcNMjQwMTE1
MTExODIzWjAXMRUwEwYDVQQDDAwqLmdvb2dsZS5jb20wWTATBgcqhkjOPQIBBggq
hkjOPQMBBwNCAARFkJP4ajCGcFu/dhDICgGR1D8gXKRbXr+9snEjrW23dyYy8ECV
0dSw/xRprfep7Tstl4B58Yv/idte6cxHVTKVo4IMGzCCDBcwDgYDVR0PAQH/BAQD
AgeAMBMGA1UdJQQMMAoGCCsGAQUFBwMBMAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYE
FAwF+qpbtYBZSyha62gYCxXYB/HZMB8GA1UdIwQYMBaAFIp0f6+Fze6VzT2c0OJG
FPNxNR0nMGoGCCsGAQUFBwEBBF4wXDAnBggrBgEFBQcwAYYbaHR0cDovL29jc3Au
cGtpLmdvb2cvZ3RzMWMzMDEGCCsGAQUFBzAChiVodHRwOi8vcGtpLmdvb2cvcmVw
by9jZXJ0cy9ndHMxYzMuZGVyMIIJzQYDVR0RBIIJxDCCCcCCDCouZ29vZ2xlLmNv
bYIWKi5hcHBlbmdpbmUuZ29vZ2xlLmNvbYIJKi5iZG4uZGV2ghUqLm9yaWdpbi10
ZXN0LmJkbi5kZXaCEiouY2xvdWQuZ29vZ2xlLmNvbYIYKi5jcm93ZHNvdXJjZS5n
b29nbGUuY29tghgqLmRhdGFjb21wdXRlLmdvb2dsZS5jb22CCyouZ29vZ2xlLmNh
ggsqLmdvb2dsZS5jbIIOKi5nb29nbGUuY28uaW6CDiouZ29vZ2xlLmNvLmpwgg4q
Lmdvb2dsZS5jby51a4IPKi5nb29nbGUuY29tLmFygg8qLmdvb2dsZS5jb20uYXWC
DyouZ29vZ2xlLmNvbS5icoIPKi5nb29nbGUuY29tLmNvgg8qLmdvb2dsZS5jb20u
bXiCDyouZ29vZ2xlLmNvbS50coIPKi5nb29nbGUuY29tLnZuggsqLmdvb2dsZS5k
ZYILKi5nb29nbGUuZXOCCyouZ29vZ2xlLmZyggsqLmdvb2dsZS5odYILKi5nb29n
bGUuaXSCCyouZ29vZ2xlLm5sggsqLmdvb2dsZS5wbIILKi5nb29nbGUucHSCEiou
Z29vZ2xlYWRhcGlzLmNvbYIPKi5nb29nbGVhcGlzLmNughEqLmdvb2dsZXZpZGVv
LmNvbYIMKi5nc3RhdGljLmNughAqLmdzdGF0aWMtY24uY29tgg9nb29nbGVjbmFw
cHMuY26CESouZ29vZ2xlY25hcHBzLmNughFnb29nbGVhcHBzLWNuLmNvbYITKi5n
b29nbGVhcHBzLWNuLmNvbYIMZ2tlY25hcHBzLmNugg4qLmdrZWNuYXBwcy5jboIS
Z29vZ2xlZG93bmxvYWRzLmNughQqLmdvb2dsZWRvd25sb2Fkcy5jboIQcmVjYXB0
Y2hhLm5ldC5jboISKi5yZWNhcHRjaGEubmV0LmNughByZWNhcHRjaGEtY24ubmV0
ghIqLnJlY2FwdGNoYS1jbi5uZXSCC3dpZGV2aW5lLmNugg0qLndpZGV2aW5lLmNu
ghFhbXBwcm9qZWN0Lm9yZy5jboITKi5hbXBwcm9qZWN0Lm9yZy5jboIRYW1wcHJv
amVjdC5uZXQuY26CEyouYW1wcHJvamVjdC5uZXQuY26CF2dvb2dsZS1hbmFseXRp
Y3MtY24uY29tghkqLmdvb2dsZS1hbmFseXRpY3MtY24uY29tghdnb29nbGVhZHNl
cnZpY2VzLWNuLmNvbYIZKi5nb29nbGVhZHNlcnZpY2VzLWNuLmNvbYIRZ29vZ2xl
dmFkcy1jbi5jb22CEyouZ29vZ2xldmFkcy1jbi5jb22CEWdvb2dsZWFwaXMtY24u
Y29tghMqLmdvb2dsZWFwaXMtY24uY29tghVnb29nbGVvcHRpbWl6ZS1jbi5jb22C
FyouZ29vZ2xlb3B0aW1pemUtY24uY29tghJkb3VibGVjbGljay1jbi5uZXSCFCou
ZG91YmxlY2xpY2stY24ubmV0ghgqLmZscy5kb3VibGVjbGljay1jbi5uZXSCFiou
Zy5kb3VibGVjbGljay1jbi5uZXSCDmRvdWJsZWNsaWNrLmNughAqLmRvdWJsZWNs
aWNrLmNughQqLmZscy5kb3VibGVjbGljay5jboISKi5nLmRvdWJsZWNsaWNrLmNu
ghFkYXJ0c2VhcmNoLWNuLm5ldIITKi5kYXJ0c2VhcmNoLWNuLm5ldIIdZ29vZ2xl
dHJhdmVsYWRzZXJ2aWNlcy1jbi5jb22CHyouZ29vZ2xldHJhdmVsYWRzZXJ2aWNl
cy1jbi5jb22CGGdvb2dsZXRhZ3NlcnZpY2VzLWNuLmNvbYIaKi5nb29nbGV0YWdz
ZXJ2aWNlcy1jbi5jb22CF2dvb2dsZXRhZ21hbmFnZXItY24uY29tghkqLmdvb2ds
ZXRhZ21hbmFnZXItY24uY29tghhnb29nbGVzeW5kaWNhdGlvbi1jbi5jb22CGiou
Z29vZ2xlc3luZGljYXRpb24tY24uY29tgiQqLnNhZmVmcmFtZS5nb29nbGVzeW5k
aWNhdGlvbi1jbi5jb22CFmFwcC1tZWFzdXJlbWVudC1jbi5jb22CGCouYXBwLW1l
YXN1cmVtZW50LWNuLmNvbYILZ3Z0MS1jbi5jb22CDSouZ3Z0MS1jbi5jb22CC2d2
dDItY24uY29tgg0qLmd2dDItY24uY29tggsybWRuLWNuLm5ldIINKi4ybWRuLWNu
Lm5ldIIUZ29vZ2xlZmxpZ2h0cy1jbi5uZXSCFiouZ29vZ2xlZmxpZ2h0cy1jbi5u
ZXSCDGFkbW9iLWNuLmNvbYIOKi5hZG1vYi1jbi5jb22CFGdvb2dsZXNhbmRib3gt
Y24uY29tghYqLmdvb2dsZXNhbmRib3gtY24uY29tgh4qLnNhZmVudXAuZ29vZ2xl
c2FuZGJveC1jbi5jb22CDSouZ3N0YXRpYy5jb22CFCoubWV0cmljLmdzdGF0aWMu
Y29tggoqLmd2dDEuY29tghEqLmdjcGNkbi5ndnQxLmNvbYIKKi5ndnQyLmNvbYIO
Ki5nY3AuZ3Z0Mi5jb22CECoudXJsLmdvb2dsZS5jb22CFioueW91dHViZS1ub2Nv
b2tpZS5jb22CCyoueXRpbWcuY29tggthbmRyb2lkLmNvbYINKi5hbmRyb2lkLmNv
bYITKi5mbGFzaC5hbmRyb2lkLmNvbYIEZy5jboIGKi5nLmNuggRnLmNvggYqLmcu
Y2+CBmdvby5nbIIKd3d3Lmdvby5nbIIUZ29vZ2xlLWFuYWx5dGljcy5jb22CFiou
Z29vZ2xlLWFuYWx5dGljcy5jb22CCmdvb2dsZS5jb22CEmdvb2dsZWNvbW1lcmNl
LmNvbYIUKi5nb29nbGVjb21tZXJjZS5jb22CCGdncGh0LmNuggoqLmdncGh0LmNu
ggp1cmNoaW4uY29tggwqLnVyY2hpbi5jb22CCHlvdXR1LmJlggt5b3V0dWJlLmNv
bYINKi55b3V0dWJlLmNvbYIUeW91dHViZWVkdWNhdGlvbi5jb22CFioueW91dHVi
ZWVkdWNhdGlvbi5jb22CD3lvdXR1YmVraWRzLmNvbYIRKi55b3V0dWJla2lkcy5j
b22CBXl0LmJlggcqLnl0LmJlghphbmRyb2lkLmNsaWVudHMuZ29vZ2xlLmNvbYIb
ZGV2ZWxvcGVyLmFuZHJvaWQuZ29vZ2xlLmNughxkZXZlbG9wZXJzLmFuZHJvaWQu
Z29vZ2xlLmNughhzb3VyY2UuYW5kcm9pZC5nb29nbGUuY24wIQYDVR0gBBowGDAI
BgZngQwBAgEwDAYKKwYBBAHWeQIFAzA8BgNVHR8ENTAzMDGgL6AthitodHRwOi8v
Y3Jscy5wa2kuZ29vZy9ndHMxYzMvbW9WRGZJU2lhMmsuY3JsMIIBAgYKKwYBBAHW
eQIEAgSB8wSB8ADuAHUASLDja9qmRzQP5WoC+p0w6xxSActW3SyB2bu/qznYhHMA
AAGLXHjeOgAABAMARjBEAiA5xkLry/jLvTgkNitDIuFmwVuL3SVofIIa9S9BX0aA
lQIgWnF3vOB4zocZKrj1Ou4RgpTEZslWMbTi36QHOmI2rc8AdQDuzdBk1dsazsVc
t520zROiModGfLzs3sNRSFlGcR+1mwAAAYtceN4RAAAEAwBGMEQCIGNaR2SEb6Eg
AXBGxhriVs6xEmMSfNl+gxHXsAT1D26qAiB/nSwVTXKptImly3Nj1RlIBQ49kG6p
LVjwM7pFTFFHgDANBgkqhkiG9w0BAQsFAAOCAQEA2zTHxu4FM4XJYuFVpYaFR+Rt
4j27U4UZ/ju4CyY1h8PozUgpap6bKvsknfpltdaoHPqMBg7XdSX1hXpVJctAxPp/
QgC6mYTlLGIk095lqyWYYH1HR/kf7VSBJDy8KYx44WliqK3kleYxbnz6BqEa8z9P
0bILW8KvHyZt5tPYBvV2R42gueVgnJT9Kht3Pv0ZblajN0Ium+S03sDdMYT4nsKU
5agIfeiO/nIgF14Zn+zKg0ilSZzum8pB0UgENi5CsowoNyOkdXOcIGMfOvy6hhXc
kE32lJu4gVL1W2fSeii8K9y7pMGNUMbV+h2sF24EGxP5zlhruE2lJGRgONjFpA==
-----END CERTIFICATE-----`
}

// GenRandLocationLabel generate a random location label.
// Return a string with the label value.
func GenRandLocationLabel() string {
	// FIXME HMS-3152 If '-' is accepted add here and update middleware
	// regular expression
	set := []string{
		"europe", "boston", "france", "australia", "india",
		"brasil", "africa", "china",
	}
	return faker.Pick(set...)
}

// GenRandLocationDescription generate a random description
// by adding " location" to a random location label.
// Return a string with the description.
func GenRandLocationDescription() string {
	return fmt.Sprintf("%s location", GenRandLocationLabel())
}

// GenIssuerWithRealm generate a certificate issuer for the
// given issuer and realm.
// issuer is a string describing the certificate issuer, for instance "Verisign".
// realm is the domain realm that the certificate is issuer belongs to.
// Return a string with the issuer string.
func GenIssuerWithRealm(issuer string, realm string) string {
	// TODO To be reviewed
	result := fmt.Sprintf("CN=%s", issuer)
	for _, item := range strings.Split(realm, ".") {
		result = fmt.Sprintf("%s, DC=%s", result, item)
	}
	return result
}

// GenSubjectWithRealm generate a certificate subject for the
// given subject and realm.
// subject is a string describing the certificate subject, for instance "My Site".
// realm is the domain realm that the certificate is issued for.
// Return a string with the subject string.
func GenSubjectWithRealm(subject string, realm string) string {
	// TODO To be reviewed
	return GenIssuerWithRealm(subject, realm)
}
