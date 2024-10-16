package sql

import (
	"fmt"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
)

func PrepSqlSelectIpas(mock sqlmock.Sqlmock, withError bool, expectedErr error, domainID uint, data *model.Domain) {
	expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ipas" WHERE id = $1 AND "ipas"."deleted_at" IS NULL ORDER BY "ipas"."id" LIMIT $2`)).
		WithArgs(
			domainID,
			1,
		)
	if withError {
		expectedQuery.WillReturnError(expectedErr)
	} else {
		expectedQuery.WillReturnRows(sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deletet_at",

			"realm_name", "realm_domains",
		}).AddRow(
			domainID,
			data.Model.CreatedAt,
			data.Model.UpdatedAt,
			data.Model.DeletedAt,

			data.IpaDomain.RealmName,
			data.IpaDomain.RealmDomains,
		))
	}
}

func PrepSqlSelectIpaCerts(mock sqlmock.Sqlmock, withError bool, expectedErr error, domainID uint, data *model.Domain) {
	expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ipa_certs" WHERE "ipa_certs"."ipa_id" = $1 AND "ipa_certs"."deleted_at" IS NULL`)).
		WithArgs(domainID)
	if withError {
		expectedQuery.WillReturnError(expectedErr)
	} else {
		rows := sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deletet_at",

			"ipa_id", "issuer", "nickname",
			"not_after", "not_before", "serial_number",
			"subject", "pem",
		})
		for j := range data.IpaDomain.CaCerts {
			rows.AddRow(
				domainID+uint(j)+1,
				data.IpaDomain.CaCerts[j].Model.CreatedAt,
				data.IpaDomain.CaCerts[j].Model.UpdatedAt,
				data.IpaDomain.CaCerts[j].Model.DeletedAt,

				domainID,
				data.IpaDomain.CaCerts[j].Issuer,
				data.IpaDomain.CaCerts[j].Nickname,
				data.IpaDomain.CaCerts[j].NotAfter,
				data.IpaDomain.CaCerts[j].NotBefore,
				data.IpaDomain.CaCerts[j].SerialNumber,
				data.IpaDomain.CaCerts[j].Subject,
				data.IpaDomain.CaCerts[j].Pem,
			)
		}
		expectedQuery.WillReturnRows(rows)
	}
}

func PrepSqlSelectIpaLocations(mock sqlmock.Sqlmock, withError bool, expectedErr error, domainID uint, data *model.Domain) {
	expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ipa_locations" WHERE "ipa_locations"."ipa_id" = $1 AND "ipa_locations"."deleted_at" IS NULL`)).
		WithArgs(domainID)
	if withError {
		expectedQuery.WillReturnError(expectedErr)
	} else {
		rows := sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deletet_at",

			"ipa_id",
			"name", "description",
		})
		for j := range data.IpaDomain.Locations {
			rows.AddRow(
				domainID+uint(j)+1,
				data.IpaDomain.Locations[j].Model.CreatedAt,
				data.IpaDomain.Locations[j].Model.UpdatedAt,
				data.IpaDomain.Locations[j].Model.DeletedAt,

				domainID,
				data.IpaDomain.Locations[j].Name,
				data.IpaDomain.Locations[j].Description,
			)
		}
		expectedQuery.WillReturnRows(rows)
	}
}

func PrepSqlSelectIpaServers(mock sqlmock.Sqlmock, withError bool, expectedErr error, domainID uint, data *model.Domain) {
	expectedQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ipa_servers" WHERE "ipa_servers"."ipa_id" = $1 AND "ipa_servers"."deleted_at" IS NULL`)).
		WithArgs(domainID)
	if withError {
		expectedQuery.WillReturnError(expectedErr)
	} else {
		rows := sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deletet_at",

			"ipa_id", "fqdn", "rhsm_id", "location",
			"ca_server", "hcc_enrollment_server", "hcc_update_server",
			"pk_init_server",
		})
		for j := range data.IpaDomain.Servers {
			rows.AddRow(
				domainID+uint(j)+1,
				data.IpaDomain.Servers[j].Model.CreatedAt,
				data.IpaDomain.Servers[j].Model.UpdatedAt,
				data.IpaDomain.Servers[j].Model.DeletedAt,

				domainID,
				data.IpaDomain.Servers[j].FQDN,
				data.IpaDomain.Servers[j].RHSMId,
				data.IpaDomain.Servers[j].Location,
				data.IpaDomain.Servers[j].CaServer,
				data.IpaDomain.Servers[j].HCCEnrollmentServer,
				data.IpaDomain.Servers[j].HCCUpdateServer,
				data.IpaDomain.Servers[j].PKInitServer,
			)
		}
		expectedQuery.WillReturnRows(rows)
	}
}

func FindIpaByID(stage int, mock sqlmock.Sqlmock, expectedErr error, domainID uint, data *model.Domain) {
	for i := 1; i <= stage; i++ {
		switch i {
		case 1:
			PrepSqlSelectIpas(mock, WithPredicateExpectedError(i, stage, expectedErr), expectedErr, domainID, data)
		case 2:
			PrepSqlSelectIpaCerts(mock, WithPredicateExpectedError(i, stage, expectedErr), expectedErr, domainID, data)
		case 3:
			PrepSqlSelectIpaLocations(mock, WithPredicateExpectedError(i, stage, expectedErr), expectedErr, domainID, data)
		case 4:
			PrepSqlSelectIpaServers(mock, WithPredicateExpectedError(i, stage, expectedErr), expectedErr, domainID, data)
		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
}
