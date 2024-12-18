package sql

import (
	"database/sql/driver"
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
			"id", "created_at", "updated_at", "deleted_at",

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
			"id", "created_at", "updated_at", "deleted_at",

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
			"id", "created_at", "updated_at", "deleted_at",

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
			"id", "created_at", "updated_at", "deleted_at",

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

func PrepSqlInsertIntoIpas(mock sqlmock.Sqlmock, withError bool, expectedErr error, domainID uint, data *model.Ipa) {
	expectQuery := mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipas" ("created_at","updated_at","deleted_at","realm_name","realm_domains","id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
		WithArgs(
			data.Model.CreatedAt,
			data.Model.UpdatedAt,
			nil,

			data.RealmName,
			data.RealmDomains,
			domainID,
		)
	if withError {
		expectQuery.WillReturnError(expectedErr)
	} else {
		expectQuery.WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(domainID))
	}
}

func PrepSqlInsertIntoIpaCerts(mock sqlmock.Sqlmock, withError bool, expectedErr error, domainID uint, data *model.Ipa) {
	for j := range data.CaCerts {
		expectQuery := mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipa_certs" ("created_at","updated_at","deleted_at","ipa_id","issuer","nickname","not_after","not_before","pem","serial_number","subject") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "id"`)).
			WithArgs(
				data.CaCerts[j].CreatedAt,
				data.CaCerts[j].UpdatedAt,
				nil,

				domainID,
				data.CaCerts[j].Issuer,
				data.CaCerts[j].Nickname,
				data.CaCerts[j].NotAfter,
				data.CaCerts[j].NotBefore,
				data.CaCerts[j].Pem,
				data.CaCerts[j].SerialNumber,
				data.CaCerts[j].Subject,
			)
		if withError {
			expectQuery.WillReturnError(expectedErr)
			return
		}
		expectQuery.WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(100*data.CaCerts[0].ID + uint(j)))
	}
}

func PrepSqlInsertIntoIpaServers(mock sqlmock.Sqlmock, withError bool, expectedErr error, domainID uint, data *model.Ipa) {
	for j := range data.Servers {
		expectQuery := mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipa_servers" ("created_at","updated_at","deleted_at","ipa_id","fqdn","rhsm_id","location","ca_server","hcc_enrollment_server","hcc_update_server","pk_init_server") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "id"`)).
			WithArgs(
				data.Servers[j].CreatedAt,
				data.Servers[j].UpdatedAt,
				nil,

				domainID,
				data.Servers[j].FQDN,
				data.Servers[j].RHSMId,
				data.Servers[j].Location,
				data.Servers[j].CaServer,
				data.Servers[j].HCCEnrollmentServer,
				data.Servers[j].HCCUpdateServer,
				data.Servers[j].PKInitServer,
			)
		if withError {
			expectQuery.WillReturnError(expectedErr)
			return
		}
		expectQuery.WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(200*data.Servers[0].ID + uint(j)))
	}
}

func PrepSqlInsertIntoIpaLocations(mock sqlmock.Sqlmock, withError bool, expectedErr error, domainID uint, data *model.Ipa) {
	for j := range data.Locations {
		expectQuery := mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ipa_locations" ("created_at","updated_at","deleted_at","ipa_id","name","description") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
			WithArgs(
				data.Locations[j].CreatedAt,
				data.Locations[j].UpdatedAt,
				nil,

				domainID,
				data.Locations[j].Name,
				data.Locations[j].Description,
			)
		if withError {
			expectQuery.WillReturnError(expectedErr)
			return
		}
		expectQuery.WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(300*data.Locations[0].ID + uint(j)))
	}
}

func CreateIpaDomain(stage int, mock sqlmock.Sqlmock, expectedErr error, domainID uint, data *model.Ipa) {
	for i := 1; i <= stage; i++ {
		switch i {
		case 1:
			PrepSqlInsertIntoIpas(mock, WithPredicateExpectedError(i, stage, expectedErr), expectedErr, domainID, data)
		case 2:
			PrepSqlInsertIntoIpaCerts(mock, WithPredicateExpectedError(i, stage, expectedErr), expectedErr, domainID, data)
		case 3:
			PrepSqlInsertIntoIpaServers(mock, WithPredicateExpectedError(i, stage, expectedErr), expectedErr, domainID, data)
		case 4:
			PrepSqlInsertIntoIpaLocations(mock, WithPredicateExpectedError(i, stage, expectedErr), expectedErr, domainID, data)
		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
}

func PrepSqlDeleteFromIpas(mock sqlmock.Sqlmock, withError bool, expectedErr error, domainID uint) {
	expectExec := mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "ipas" WHERE "ipas"."id" = $1`)).
		WithArgs(
			domainID,
		)
	if withError {
		expectExec.WillReturnError(expectedErr)
	} else {
		expectExec.WillReturnResult(
			driver.RowsAffected(1),
		)
	}
}

func UpdateIpaDomain(stage int, mock sqlmock.Sqlmock, expectedErr error, domainID uint, data *model.Domain) {
	if stage == 0 {
		return
	}
	if stage < 0 {
		panic("'stage' cannot be lower than 0")
	}
	if stage > 5 {
		panic("'stage' cannot be greater than 5")
	}

	mock.MatchExpectationsInOrder(true)
	for i := 1; i <= stage; i++ {
		switch i {
		case 1:
			PrepSqlDeleteFromIpas(mock, WithPredicateExpectedError(i, stage, expectedErr), expectedErr, domainID)
		case 2:
			PrepSqlInsertIntoIpas(mock, WithPredicateExpectedError(i, stage, expectedErr), expectedErr, domainID, data.IpaDomain)
		case 3:
			PrepSqlInsertIntoIpaCerts(mock, WithPredicateExpectedError(i, stage, expectedErr), expectedErr, domainID, data.IpaDomain)
		case 4:
			PrepSqlInsertIntoIpaServers(mock, WithPredicateExpectedError(i, stage, expectedErr), expectedErr, domainID, data.IpaDomain)
		case 5:
			PrepSqlInsertIntoIpaLocations(mock, WithPredicateExpectedError(i, stage, expectedErr), expectedErr, domainID, data.IpaDomain)
		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
}
