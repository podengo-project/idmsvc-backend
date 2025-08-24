package impl

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getDatabase() (*gorm.DB, sqlmock.Sqlmock) {
	var (
		db     *sql.DB
		gormDb *gorm.DB
		mock   sqlmock.Sqlmock
		err    error
	)

	db, mock, _ = sqlmock.New()
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gormDb, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return gormDb, mock
}

func TestNewTodoCreatedEventHandler(t *testing.T) {
	var (
		result event.Eventable
		gormDb *gorm.DB
	)

	// When nil is passed it returns nil
	result = NewTodoCreatedEventHandler(nil)
	assert.Nil(t, result)

	// https://github.com/go-gorm/gorm/issues/3565
	// When a gormDb connector is passed
	gormDb, _ = getDatabase()

	result = NewTodoCreatedEventHandler(gormDb)
	assert.NotNil(t, result)
}

func TestTodoCreatedEventHandlerOnMessage(t *testing.T) {
	var err error

	type TestCase struct {
		Name     string
		Given    *kafka.Message
		Expected error
	}

	testCases := []TestCase{
		{
			Name: "Get Not implemented error",

			Given: &kafka.Message{
				Value: []byte(`{
					"id": 12345,
					"title": "todo title",
					"decription": "todo description"
				}`),
			},
			Expected: fmt.Errorf("Not implemented"),
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.Name)
		gormDb, _ := getDatabase()
		require.NotNil(t, gormDb)
		handler := NewTodoCreatedEventHandler(gormDb)
		require.NotNil(t, handler)
		err = handler.OnMessage(testCase.Given)
		if testCase.Expected != nil {
			require.Error(t, err)
			assert.Contains(t, err.Error(), testCase.Expected.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}
