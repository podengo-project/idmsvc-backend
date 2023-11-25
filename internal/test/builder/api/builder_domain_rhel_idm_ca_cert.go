package api

import (
	"strconv"
	"time"

	"github.com/podengo-project/idmsvc-backend/internal/api/public"
	builder_helper "github.com/podengo-project/idmsvc-backend/internal/test/builder/helper"
)

type Certificate interface {
	Build() public.Certificate
	WithIssuer(value string) Certificate
	WithNickname(value string) Certificate
	WithNotAfter(value time.Time) Certificate
	WithNotBefore(value time.Time) Certificate
	WithSerialNumber(value string) Certificate
	WithSubject(value string) Certificate
	WithPem(value string) Certificate
}

type certificate public.Certificate

func NewCertificate(realm string) Certificate {
	notBefore := builder_helper.GenPastNearTime(1 * time.Hour)
	notAfter := builder_helper.GenPastNearTime(1 * time.Hour)
	return &certificate{
		// TODO Add random function for the issuer
		Issuer: builder_helper.GenIssuerWithRealm(
			"An Issuer", realm),
		// TODO Fill correctly this field
		Nickname:  "My Test Certificate",
		NotAfter:  notAfter,
		NotBefore: notBefore,
		// TODO Is a UUID a valid serial number in a certificate?
		// SerialNumber: uuid.NewString(),
		SerialNumber: strconv.Itoa(int(builder_helper.GenRandNum(1, 99999999))),
		// TODO Add random function for the subject
		Subject: builder_helper.GenSubjectWithRealm(
			"A Subject", realm),
		Pem: builder_helper.GenPemCertificate(),
	}
}

func (b *certificate) Build() public.Certificate {
	return public.Certificate(*b)
}

func (b *certificate) SetIssuer(value string) Certificate {
	b.Issuer = value
	return b
}

func (b *certificate) SetNickname(value string) Certificate {
	return b
}

func (b *certificate) SetNotAfte(value time.Duration) Certificate {
	b.SetNotAfte(value)
	return b
}

func (b *certificate) WithIssuer(value string) Certificate {
	b.Issuer = value
	return b
}

func (b *certificate) WithNickname(value string) Certificate {
	b.Nickname = value
	return b
}

func (b *certificate) WithNotAfter(value time.Time) Certificate {
	b.NotAfter = value
	return b
}
func (b *certificate) WithNotBefore(value time.Time) Certificate {
	b.NotBefore = value
	return b
}

func (b *certificate) WithSerialNumber(value string) Certificate {
	b.SerialNumber = value
	return b
}

func (b *certificate) WithSubject(value string) Certificate {
	b.Subject = value
	return b
}

func (b *certificate) WithPem(value string) Certificate {
	b.Pem = value
	return b
}
