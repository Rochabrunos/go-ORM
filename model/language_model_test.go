package model

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	mocks "orm-golang/model/mock"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	var db = GetDBTestConnection()

	fmt.Print("Migrating the Model Language to the Test Database\n")
	if err := db.AutoMigrate(&Language{}); err != nil {
		fmt.Errorf("Fail to migrate the model Language: %v\n", err)
	}
	fmt.Printf("Migration has been successful\n")
}

func TestLanguageGetById(t *testing.T) {
	t.Run("Must return an error when called with an invalid ID", func(t *testing.T) {
		var model = &LanguageModel{}
		var db = GetDBTestConnection()
		var ctx = MockContext()
		ctx.AddParam("id", "A")

		err := model.GetById(ctx, db)

		assert.ErrorContains(t, err, "invalid id, make sure to pass a number")
	})

	t.Run("Must return an error if database fails", func(t *testing.T) {
		var model = &LanguageModel{}
		var ctx = MockContext()
		var db *gorm.DB
		var expectedError = errors.New(mocks.MockedErrorMessage)
		var mock sqlmock.Sqlmock
		var conn *sql.DB

		conn, mock, err := sqlmock.New()
		assert.Nil(t, err)
		defer conn.Close()

		dialector := postgres.New(postgres.Config{
			DSN:        "sqlmock_db_0",
			DriverName: "postgres",
			Conn:       conn,
		})

		db, err = gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)
		assert.NotNil(t, db)

		mock.
			ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "language" WHERE "language"."language_id" = $1 ORDER BY "language"."language_id" LIMIT 1`)).
			WithArgs(1).WillReturnError(expectedError)

		ctx.AddParam("id", "1")

		err = model.GetById(ctx, db)

		assert.ErrorIs(t, err, expectedError)
	})

	t.Run("Must return a language with the specified ID", func(t *testing.T) {
		var model = &LanguageModel{}
		var db = GetDBTestConnection()
		var ctx = MockContext()
		ctx.AddParam("id", "1")

		db.Exec("TRUNCATE TABLE language RESTART IDENTITY CASCADE")
		db.Create(&Language{Name: "English"})

		err := model.GetById(ctx, db)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), model.Languages[0].ID)

		db.Delete(&Language{ID: 1})
	})
}

func TestGetAllLanguages(t *testing.T) {

	t.Run("Must return an error if the database fails to retrieve objects", func(t *testing.T) {
		var model = &LanguageModel{}
		var expectedError = errors.New(mocks.MockedErrorMessage)
		var ctx = MockContext()
		var mock sqlmock.Sqlmock
		var conn *sql.DB

		conn, mock, err := sqlmock.New()
		assert.Nil(t, err)
		defer conn.Close()

		dialector := postgres.New(postgres.Config{
			DSN:        "sqlmock_db_0",
			DriverName: "postgres",
			Conn:       conn,
		})

		db, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)
		assert.NotNil(t, db)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"language\" LIMIT 10")).WithArgs().WillReturnError(expectedError)

		err = model.GetAll(ctx, db)

		assert.Error(t, err, expectedError)
	})

	t.Run("Must return the languages store in the database", func(t *testing.T) {
		var model = &LanguageModel{}
		var ctx = MockContext()
		var db = GetDBTestConnection()
		db.Exec("TRUNCATE TABLE language RESTART IDENTITY CASCADE;")
		db.Create([]Language{{Name: "English"}, {Name: "Spanish"}})

		err := model.GetAll(ctx, db)

		assert.NoError(t, err)
		assert.NotNil(t, model.Languages)
		assert.Equal(t, 2, len(model.Languages))

		db.Delete(&Language{ID: 1})
		db.Delete(&Language{ID: 2})
	})
}

func TestCreateNewLanguage(t *testing.T) {
	t.Run("Must return an error if called with invalid data", func(t *testing.T) {
		var model = &LanguageModel{}
		var ctx = MockContext()
		var db = GetDBTestConnection()

		ctx.Request = &http.Request{Header: http.Header{}}
		ctx.Request.Header.Set("content-type", "application/json")

		data := map[string]any{"Mock": "mocked name"}
		jsonBytes, _ := json.Marshal(data)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))

		err := model.CreateNew(ctx, db)

		assert.ErrorContains(t, err, "Error:Field validation for")
	})

	t.Run("Must return an error if it's not possible to create the new language", func(t *testing.T) {
		var model = &LanguageModel{}
		var expectedError = errors.New(mocks.MockedErrorMessage)
		var mock sqlmock.Sqlmock
		var conn *sql.DB
		conn, mock, err := sqlmock.New()
		assert.Nil(t, err)
		defer conn.Close()

		dialector := postgres.New(postgres.Config{
			DSN:        "sqlmock_db_0",
			DriverName: "postgres",
			Conn:       conn,
		})

		db, err := gorm.Open(dialector, &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		assert.NoError(t, err)
		assert.NotNil(t, db)

		ctx := MockContext()
		ctx.Request = &http.Request{Header: http.Header{}}
		ctx.Request.Header.Set("content-type", "application/json")

		data := map[string]any{"Name": "English"}
		jsonBytes, _ := json.Marshal(data)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "language" ("name","last_update") VALUES ($1,$2) RETURNING "language_id"`)).
			WithArgs(data["Name"], AnyTime{}).
			WillReturnError(expectedError)
		mock.ExpectRollback()

		err = model.CreateNew(ctx, db)

		assert.ErrorIs(t, err, expectedError)
	})

	t.Run("Must return a new language when called", func(t *testing.T) {
		var model = &LanguageModel{}
		var db = GetDBTestConnection()
		var ctx = MockContext()
		ctx.Request = &http.Request{Header: http.Header{}}
		ctx.Request.Header.Set("content-type", "application/json")

		data := map[string]any{"Name": "English"}
		jsonBytes, _ := json.Marshal(data)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))

		err := model.CreateNew(ctx, db)

		assert.NoError(t, err)
		assert.Equal(t, data["Name"], model.Languages[0].Name)
	})
}

func TestUpdateLanguageById(t *testing.T) {
	t.Run("Must return an error when the language doesn't exists", func(t *testing.T) {
		var model = &LanguageModel{}
		var expectedError = errors.New(mocks.MockedErrorMessage)
		var ctx = MockContext()
		var db *gorm.DB
		var mock sqlmock.Sqlmock
		var conn *sql.DB
		conn, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer conn.Close()

		dialect := postgres.New(postgres.Config{
			DSN:        "sqlmock_db_0",
			DriverName: "postgres",
			Conn:       conn,
		})

		db, err = gorm.Open(dialect, &gorm.Config{})
		assert.NoError(t, err)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "language" WHERE "language"."language_id" = $1 ORDER BY "language"."language_id" LIMIT 1`)).
			WithArgs(1).
			WillReturnError(expectedError)

		ctx.AddParam("id", "1")

		err = model.UpdateById(ctx, db)

		assert.ErrorIs(t, expectedError, err)
	})

	t.Run("Must return an error when called with invalid JSON body", func(t *testing.T) {
		var model = &LanguageModel{}
		var db = GetDBTestConnection()
		var ctx = MockContext()
		ctx.AddParam("id", "1")
		ctx.Request = &http.Request{}

		db.Exec(`TRUNCATE TABLE language RESTART IDENTITY CASCADE`)
		result := db.Create(&Language{Name: "English"})
		assert.NoError(t, result.Error)

		data := map[string]string{"English": "Name"}
		jsonBytes, err := json.Marshal(data)
		assert.NoError(t, err)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))

		err = model.UpdateById(ctx, db)

		assert.ErrorContains(t, err, "Error:Field validation for 'Name'")

		db.Delete(&Language{ID: 1})
	})

	t.Run("Must return an error if the database fail to update the language", func(t *testing.T) {
		var model = &LanguageModel{}
		var expectedError = errors.New(mocks.MockedErrorMessage)
		var ctx = MockContext()
		var db *gorm.DB
		var conn *sql.DB
		var mock sqlmock.Sqlmock

		conn, mock, err := sqlmock.New()
		assert.NoError(t, err)

		dialect := postgres.New(postgres.Config{
			DSN:        "sql_db_0",
			DriverName: "postgres",
			Conn:       conn,
		})

		db, err = gorm.Open(dialect, &gorm.Config{})
		assert.NoError(t, err)

		ctx.Request = &http.Request{}
		data := map[string]string{"Name": "English"}
		jsonBytes, err := json.Marshal(data)
		assert.NoError(t, err)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
		ctx.AddParam("id", "1")

		mockedLanguage := Language{ID: 1, Name: "Mocked Language Row", LastUpdate: time.Now()}
		mock.
			ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "language" WHERE "language"."language_id" = $1 ORDER BY "language"."language_id" LIMIT 1`)).
			WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"language_id", "name", "last_update"}).
			AddRow(mockedLanguage.ID, mockedLanguage.Name, mockedLanguage.LastUpdate))

		mock.ExpectBegin()
		mock.
			ExpectExec(regexp.QuoteMeta(`UPDATE "language" SET "name"=$1,"last_update"=$2 WHERE "language_id" = $3`)).
			WithArgs(data["Name"], AnyTime{}, mockedLanguage.ID).
			WillReturnError(expectedError)
		mock.ExpectRollback()

		err = model.UpdateById(ctx, db)

		assert.ErrorIs(t, expectedError, err)
	})

	t.Run("Should retrun a modified language", func(t *testing.T) {
		var model = &LanguageModel{}
		var db = GetDBTestConnection()
		var ctx = MockContext()
		ctx.AddParam("id", "1")
		ctx.Request = &http.Request{}

		db.Exec(`TRUNCATE TABLE language RESTART IDENTITY CASCADE`)
		db.Create(&Language{Name: "English"})

		data := map[string]string{"Name": "Spanish"}
		jsonData, err := json.Marshal(data)
		assert.NoError(t, err)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonData))

		err = model.UpdateById(ctx, db)

		assert.NoError(t, err)
		assert.Equal(t, data["Name"], model.Languages[0].Name)

		db.Delete(&Language{ID: 1})
	})
}

func TestDeleteLanguageById(t *testing.T) {

	t.Run("Must return an error if the Language with the id doesn't exist", func(t *testing.T) {
		var model = &LanguageModel{}
		var db = GetDBTestConnection()
		var ctx = MockContext()
		ctx.AddParam("id", "1")

		err := model.DeleteById(ctx, db)

		assert.ErrorContains(t, err, "record not found")
	})

	t.Run("Must return an error if the database fail to delete the language", func(t *testing.T) {
		var model = &LanguageModel{}
		var expectedError = errors.New(mocks.MockedErrorMessage)
		var ctx = MockContext()
		var db *gorm.DB
		var conn *sql.DB
		var mock sqlmock.Sqlmock

		conn, mock, err := sqlmock.New()
		assert.NoError(t, err)

		dialect := postgres.New(postgres.Config{
			DSN:        "sql_db_0",
			DriverName: "postgres",
			Conn:       conn,
		})

		db, err = gorm.Open(dialect, &gorm.Config{})
		assert.NoError(t, err)

		mockedLanguage := Language{ID: 1, Name: "Mocked Language Row", LastUpdate: time.Now()}
		mock.
			ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "language" WHERE "language"."language_id" = $1 ORDER BY "language"."language_id" LIMIT 1`)).
			WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"language_id", "name", "last_update"}).
			AddRow(mockedLanguage.ID, mockedLanguage.Name, mockedLanguage.LastUpdate))

		mock.ExpectBegin()
		mock.
			ExpectExec(regexp.QuoteMeta(`DELETE FROM "language" WHERE "language"."language_id" = $1`)).
			WithArgs(mockedLanguage.ID).
			WillReturnError(expectedError)
		mock.ExpectRollback()

		ctx.AddParam("id", "1")

		err = model.DeleteById(ctx, db)

		assert.ErrorIs(t, expectedError, err)
	})

	t.Run("Must delete a language when called", func(t *testing.T) {
		var model = &LanguageModel{}
		var db = GetDBTestConnection()
		var ctx = MockContext()
		ctx.AddParam("id", "1")

		db.Exec(`TRUNCATE TABLE language RESTART IDENTITY CASCADE`)
		db.Create(&Language{Name: "English"})
		err := model.DeleteById(ctx, db)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), model.Languages[0].ID)
		assert.Equal(t, "English", model.Languages[0].Name)

		db.Delete(&Language{ID: 1})
	})
}
