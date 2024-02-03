package model

import (
	"strconv"
	"time"

	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	"github.com/podengo-project/idmsvc-backend/internal/test/builder/helper"
	"gorm.io/gorm"
)

type IpaCert interface {
	Build() model.IpaCert
	WithIpaID(value uint) IpaCert
	WithModel(value gorm.Model) IpaCert
	WithIssuer(value string) IpaCert
	WithSubject(value string) IpaCert
	WithNickname(value string) IpaCert
	WithSerialNumber(value string) IpaCert
	WithPem(value string) IpaCert
	WithNotAfter(value time.Time) IpaCert
	WithNotBefore(value time.Time) IpaCert
}

type ipaCert struct {
	IpaCert model.IpaCert
}

func NewIpaCert(gormModel gorm.Model, realm string) IpaCert {
	return &ipaCert{
		IpaCert: model.IpaCert{
			Model:        gormModel,
			IpaID:        uint(helper.GenRandNum(1, 999999)),
			Issuer:       helper.GenIssuerWithRealm(helper.GenRandDomainLabel(), realm),
			Subject:      helper.GenSubjectWithRealm(helper.GenRandDomainLabel(), realm),
			Nickname:     helper.GenRandEmail(),
			SerialNumber: strconv.FormatInt(helper.GenRandNum(1, 99999999), 10),
			Pem:          helper.GenPemCertificate(),
			NotAfter:     helper.GenFutureNearTimeUTC(time.Hour),
			NotBefore:    helper.GenPastNearTime(time.Hour),
		},
	}
}

func (b *ipaCert) Build() model.IpaCert {
	return b.IpaCert
}

func (b *ipaCert) WithIpaID(value uint) IpaCert {
	b.IpaCert.IpaID = value
	return b
}
func (b *ipaCert) WithModel(value gorm.Model) IpaCert {
	b.IpaCert.Model = value
	return b
}
func (b *ipaCert) WithIssuer(value string) IpaCert {
	b.IpaCert.Issuer = value
	return b
}
func (b *ipaCert) WithSubject(value string) IpaCert {
	b.IpaCert.Subject = value
	return b
}
func (b *ipaCert) WithNickname(value string) IpaCert {
	b.IpaCert.Nickname = value
	return b
}
func (b *ipaCert) WithSerialNumber(value string) IpaCert {
	b.IpaCert.SerialNumber = value
	return b
}
func (b *ipaCert) WithPem(value string) IpaCert {
	b.IpaCert.Pem = value
	return b
}
func (b *ipaCert) WithNotAfter(value time.Time) IpaCert {
	b.IpaCert.NotAfter = value
	return b
}
func (b *ipaCert) WithNotBefore(value time.Time) IpaCert {
	b.IpaCert.NotBefore = value
	return b
}
