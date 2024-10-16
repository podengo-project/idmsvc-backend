package sql

import (
	"database/sql/driver"
	"fmt"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/podengo-project/idmsvc-backend/internal/domain/model"
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
			"id", "created_at", "updated_at", "deletet_at",

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
