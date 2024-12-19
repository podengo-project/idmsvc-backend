package sql

import (
	"database/sql/driver"
	"fmt"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
	"github.com/podengo-project/idmsvc-backend/internal/interface/interactor"
	"gorm.io/gorm"
)

func PrepSqlSelectDomainsByID(mock sqlmock.Sqlmock, withError bool, expectedErr error, domainID uint, data *model.Domain) {
	expectQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "domains" WHERE (org_id = $1 AND domain_uuid = $2) AND "domains"."deleted_at" IS NULL ORDER BY "domains"."id" LIMIT $3`)).
		WithArgs(
			data.OrgId,
			data.DomainUuid,
			1,
		)
	if withError {
		expectQuery.WillReturnError(expectedErr)
	} else {
		autoenrollment := false
		if data.AutoEnrollmentEnabled != nil {
			autoenrollment = *data.AutoEnrollmentEnabled
		}
		expectQuery.WillReturnRows(sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deleted_at",

			"org_id", "domain_uuid", "domain_name",
			"title", "description", "type",
			"auto_enrollment_enabled",
		}).
			AddRow(
				domainID,
				data.CreatedAt,
				data.UpdatedAt,
				nil,

				data.OrgId,
				data.DomainUuid,
				data.DomainName,
				data.Title,
				data.Description,
				data.Type,
				autoenrollment,
			))
	}
}

func FindByID(stage int, mock sqlmock.Sqlmock, expectedErr error, domainID uint, data *model.Domain) {
	for i := 1; i <= stage; i++ {
		switch i {
		case 1:
			PrepSqlSelectDomainsByID(mock, WithPredicateExpectedError(i, stage, expectedErr), expectedErr, domainID, data)
		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
}

func PrepSqlSelectCountDomainsByID(mock sqlmock.Sqlmock, withError bool, expectedErr error, data *model.Domain) {
	expectQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "domains" WHERE (org_id = $1 AND domain_uuid = $2) AND "domains"."deleted_at" IS NULL LIMIT $3`)).
		WithArgs(
			data.OrgId,
			data.DomainUuid,
			1,
		)
	if withError {
		if expectedErr == gorm.ErrRecordNotFound {
			expectQuery.WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(int64(0)))
		} else {
			expectQuery.WillReturnError(expectedErr)
		}
	} else {
		expectQuery.WillReturnRows(sqlmock.NewRows([]string{"count"}).
			AddRow(int64(1)))
	}
}

func PrepSqlDeleteDomainsByID(mock sqlmock.Sqlmock, withError bool, expectedErr error, data *model.Domain) {
	expectQuery := mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "domains" WHERE (org_id = $1 AND domain_uuid = $2) AND "domains"."id" = $3`)).
		WithArgs(
			data.OrgId,
			data.DomainUuid,
			data.ID,
		)
	if withError {
		expectQuery.WillReturnError(expectedErr)
	} else {
		expectQuery.WillReturnResult(driver.RowsAffected(1))
	}
}

func DeleteByID(stage int, mock sqlmock.Sqlmock, expectedErr error, data *model.Domain) {
	for i := 1; i <= stage; i++ {
		switch i {
		case 1:
			PrepSqlSelectDomainsByID(mock, WithPredicateExpectedError(i, stage, expectedErr), expectedErr, uint(1), data)
		case 2:
			PrepSqlSelectCountDomainsByID(mock, WithPredicateExpectedError(i, stage, expectedErr), expectedErr, data)
		case 3:
			PrepSqlDeleteDomainsByID(mock, WithPredicateExpectedError(i, stage, expectedErr), expectedErr, data)
		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
}

func PrepSqlUpdateDomainsForAgent(mock sqlmock.Sqlmock, withError bool, expectedErr error, domainID uint, data *model.Domain) {
	expectExec := mock.ExpectExec(regexp.QuoteMeta(`UPDATE "domains" SET "created_at"=$1,"updated_at"=$2,"org_id"=$3,"domain_uuid"=$4,"domain_name"=$5,"title"=$6,"description"=$7,"type"=$8,"auto_enrollment_enabled"=$9 WHERE (org_id = $10 AND domain_uuid = $11) AND "domains"."deleted_at" IS NULL AND "id" = $12`)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),

			data.OrgId,
			data.DomainUuid,
			data.DomainName,

			data.Title,
			data.Description,
			data.Type,
			data.AutoEnrollmentEnabled,

			data.OrgId,
			data.DomainUuid,
			data.ID,
		)
	if withError {
		expectExec.WillReturnError(expectedErr)
	} else {
		expectExec.WillReturnResult(
			driver.RowsAffected(1))
	}
}

func PrepSqlUpdateDomainsForUser(mock sqlmock.Sqlmock, withError bool, expectedErr error, domainID uint, data *model.Domain) {
	expectExec := mock.ExpectExec(regexp.QuoteMeta(`UPDATE "domains" SET "auto_enrollment_enabled"=$1,"description"=$2,"title"=$3 WHERE (org_id = $4 AND domain_uuid = $5) AND "domains"."deleted_at" IS NULL AND "id" = $6`)).
		WithArgs(
			data.AutoEnrollmentEnabled,
			data.Description,
			data.Title,

			data.OrgId,
			data.DomainUuid,
			domainID,
		)
	if withError {
		expectExec.WillReturnError(expectedErr)
	} else {
		expectExec.WillReturnResult(
			driver.RowsAffected(1))
	}
}

func UpdateUser(stage int, mock sqlmock.Sqlmock, expectedErr error, domainID uint, data *model.Domain) {
	if stage == 0 {
		return
	}
	if stage < 0 {
		panic("'stage' cannot be lower than 0")
	}
	if stage > 2 {
		panic("'stage' cannot be greater than 3")
	}

	mock.MatchExpectationsInOrder(true)
	for i := 1; i <= stage; i++ {
		switch i {
		case 1:
			if i == stage && expectedErr != nil {
				FindByID(1, mock, expectedErr, domainID, data)
			} else {
				FindByID(1, mock, nil, domainID, data)
				FindIpaByID(4, mock, nil, domainID, data)
			}
		case 2: // Update
			PrepSqlUpdateDomainsForUser(mock, WithPredicateExpectedError(i, stage, expectedErr), expectedErr, domainID, data)
		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
}

func PrepSqlInsertIntoDomains(mock sqlmock.Sqlmock, withError bool, expectedErr error, domainID uint, data *model.Domain) {
	expectQuery := mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "domains" ("created_at","updated_at","deleted_at","org_id","domain_uuid","domain_name","title","description","type","auto_enrollment_enabled","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "id"`)).
		WithArgs(
			data.Model.CreatedAt,
			data.Model.UpdatedAt,
			data.Model.DeletedAt,

			data.OrgId,
			data.DomainUuid,
			data.DomainName,
			data.Title,
			data.Description,
			data.Type,
			data.AutoEnrollmentEnabled,

			data.Model.ID,
		)
	if withError {
		expectQuery.WillReturnError(expectedErr)
	} else {
		expectQuery.WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(data.Model.ID))
	}
}

func Register(stage int, mock sqlmock.Sqlmock, expectedErr error, data *model.Domain) {
	for i := 1; i <= stage; i++ {
		switch i {
		case 1:
			PrepSqlInsertIntoDomains(mock, WithPredicateExpectedError(i, stage, expectedErr), expectedErr, uint(1), data)
		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
}

func PrepSqlSelectFromDomainsFilterMatchDomain(mock sqlmock.Sqlmock, withError bool, expectedErr error, options *interactor.HostConfOptions, domains []model.Domain) {
	expectQuery := mock.ExpectQuery(regexp.QuoteMeta(`SELECT "domains"."id","domains"."created_at","domains"."updated_at","domains"."deleted_at","domains"."org_id","domains"."domain_uuid","domains"."domain_name","domains"."title","domains"."description","domains"."type","domains"."auto_enrollment_enabled" FROM "domains" left join ipas on domains.id = ipas.id WHERE domains.org_id = $1 AND domains.domain_uuid = $2 AND domains.domain_name = $3 AND domains.type = $4 AND "domains"."deleted_at" IS NULL`)).
		WithArgs(
			options.OrgId,
			options.DomainId,
			options.DomainName,
			model.DomainTypeUint((string)(*options.DomainType)),
		)
	if withError {
		expectQuery.WillReturnError(expectedErr)
	} else {
		rows := sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deleted_at",

			"org_id", "domain_uuid", "domain_name",
			"title", "description", "type",
			"auto_enrollment_enabled",
		})
		for j := range domains {
			rows.AddRow(
				domains[j].ID,
				domains[j].CreatedAt,
				domains[j].UpdatedAt,
				domains[j].DeletedAt,

				domains[j].OrgId,
				domains[j].DomainUuid,
				domains[j].DomainName,
				domains[j].Title,
				domains[j].Description,
				domains[j].Type,
				domains[j].AutoEnrollmentEnabled,
			)
		}
		expectQuery = expectQuery.WillReturnRows(rows)
	}
}

func MatchDomain(stage int, mock sqlmock.Sqlmock, expectedErr error, options *interactor.HostConfOptions, domains []model.Domain) {
	for i := 1; i <= stage; i++ {
		switch i {
		case 1:
			PrepSqlSelectFromDomainsFilterMatchDomain(mock, WithPredicateExpectedError(i, stage, expectedErr), expectedErr, options, domains)
		case 2:
			if len(domains) == 0 {
				FindIpaByID(1, mock, expectedErr, domains[0].ID, &domains[0])
			}
			FindIpaByID(4, mock, expectedErr, domains[0].ID, &domains[0])
		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
}

func UpdateAgent(stage int, mock sqlmock.Sqlmock, expectedErr error, domainID uint, data *model.Domain) {
	for i := 1; i <= stage; i++ {
		switch i {
		case 1:
			if WithPredicateExpectedError(i, stage, expectedErr) {
				FindByID(1, mock, expectedErr, domainID, data)
			} else {
				FindByID(1, mock, nil, domainID, data)
				FindIpaByID(4, mock, nil, domainID, data)
			}
		case 2:
			PrepSqlUpdateDomainsForAgent(mock, WithPredicateExpectedError(i, stage, expectedErr), expectedErr, domainID, data)
		case 3:
			PrepSqlDeleteFromIpas(mock, WithPredicateExpectedError(i, stage, expectedErr), expectedErr, domainID)
		case 4:
			CreateIpaDomain(4, mock, expectedErr, domainID, data.IpaDomain)
		default:
			panic(fmt.Sprintf("scenario %d/%d is not supported", i, stage))
		}
	}
}
