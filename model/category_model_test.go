package model

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"testing"
	"time"

	mocks "orm-golang/model/mock"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func init() {
	var db = GetDBTestConnection()

	fmt.Print("Migrating the Model Category to the Test Database\n")
	if err := db.AutoMigrate(&Category{}); err != nil {
		fmt.Errorf("Fail to migrate the model Category: %v\n", err)
	}
	db.Exec("TRUNCATE TABLE category RESTART IDENTITY CASCADE;")
	fmt.Printf("Migration has been successful\n")
}

func TestGetById(t *testing.T) {
	t.Run("Must return an error when called with an invalid ID", func(t *testing.T) {
		var model = &CategoryModel{}
		var db = GetDBTestConnection()
		var ctx = MockContext()
		ctx.AddParam("id", "A")

		err := model.GetById(ctx, db)

		assert.ErrorContains(t, err, "invalid id, make sure to pass a number")
	})

	t.Run("Must return an error if database fails", func(t *testing.T) {
		var model = &CategoryModel{}
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
			ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category" WHERE "category"."category_id" = $1 ORDER BY "category"."category_id" LIMIT 1`)).
			WithArgs(1).WillReturnError(expectedError)

		ctx.AddParam("id", "1")

		err = model.GetById(ctx, db)

		assert.ErrorIs(t, err, expectedError)
	})

	t.Run("Must return a category with the specified ID", func(t *testing.T) {
		var model = &CategoryModel{}
		var db = GetDBTestConnection()
		var ctx = MockContext()
		ctx.AddParam("id", "1")

		db.Exec("TRUNCATE TABLE category RESTART IDENTITY CASCADE")
		db.Create(&Category{Name: "Action"})

		err := model.GetById(ctx, db)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), model.Categories[0].ID)

		db.Delete(&Category{ID: 1})
	})
}

func TestGetAll(t *testing.T) {

	//Following the logic of the function, the only way to return an error is if the db.Find function fails.
	//This test mocks a failure in order to verify the correctness of the entire function in the event of an error
	t.Run("Must return an error if the database fails to retrieve objects", func(t *testing.T) {
		var model = &CategoryModel{}
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

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"category\" LIMIT 10")).WithArgs().WillReturnError(expectedError)

		err = model.GetAll(ctx, db)

		assert.Error(t, err, expectedError)
	})

	t.Run("Must return the categories store in the database", func(t *testing.T) {
		var model = &CategoryModel{}
		var ctx = MockContext()
		var db = GetDBTestConnection()
		db.Exec("TRUNCATE TABLE category RESTART IDENTITY CASCADE;")
		db.Create([]Category{{Name: "Action"}, {Name: "Drama"}})

		err := model.GetAll(ctx, db)

		assert.NoError(t, err)
		assert.NotNil(t, model.Categories)
		assert.Equal(t, 2, len(model.Categories))

		db.Delete(&Category{ID: 1})
		db.Delete(&Category{ID: 2})
	})
}

func TestCreateNew(t *testing.T) {
	t.Run("Must return an error if called with invalid data", func(t *testing.T) {
		var model = &CategoryModel{}
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

	t.Run("Must return an error if it's not possible to create the new category", func(t *testing.T) {
		var model = &CategoryModel{}
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

		data := map[string]any{"Name": "Action"}
		jsonBytes, _ := json.Marshal(data)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "category" ("name","last_update") VALUES ($1,$2) RETURNING "category_id"`)).WithArgs(data["Name"], AnyTime{}).WillReturnError(expectedError)
		mock.ExpectRollback()

		err = model.CreateNew(ctx, db)

		assert.ErrorIs(t, err, expectedError)
	})

	t.Run("Must return a new category when called", func(t *testing.T) {
		var model = &CategoryModel{}
		var db = GetDBTestConnection()
		var ctx = MockContext()
		ctx.Request = &http.Request{Header: http.Header{}}
		ctx.Request.Header.Set("content-type", "application/json")

		data := map[string]any{"Name": "Drama"}
		jsonBytes, _ := json.Marshal(data)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))

		err := model.CreateNew(ctx, db)

		assert.NoError(t, err)
		assert.Equal(t, "Drama", model.Categories[0].Name)
	})
}

func TestUpdateById(t *testing.T) {
	t.Run("Must return an error when the category doesn't exists", func(t *testing.T) {
		var model = &CategoryModel{}
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

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category" WHERE "category"."category_id" = $1 ORDER BY "category"."category_id" LIMIT 1`)).WithArgs(1).WillReturnError(expectedError)

		ctx.AddParam("id", "1")

		err = model.UpdateById(ctx, db)

		assert.ErrorIs(t, expectedError, err)
	})

	t.Run("Must return an error when called with invalid JSON body", func(t *testing.T) {
		var model = &CategoryModel{}
		var db = GetDBTestConnection()
		var ctx = MockContext()
		ctx.AddParam("id", "1")
		ctx.Request = &http.Request{}

		db.Exec(`TRUNCATE TABLE category RESTART IDENTITY CASCADE`)
		result := db.Create(&Category{Name: "Action"})
		assert.NoError(t, result.Error)

		data := map[string]string{"Drama": "Name"}
		jsonBytes, err := json.Marshal(data)
		assert.NoError(t, err)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))

		err = model.UpdateById(ctx, db)

		assert.ErrorContains(t, err, "Error:Field validation for 'Name'")

		db.Delete(&Category{ID: 1})
	})

	t.Run("Must return an error if the database fail to update the category", func(t *testing.T) {
		var model = &CategoryModel{}
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
		data := map[string]string{"Name": "Drama"}
		jsonBytes, err := json.Marshal(data)
		assert.NoError(t, err)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
		ctx.AddParam("id", "1")

		mockedCategory := Category{ID: 1, Name: "Mocked Category Row", LastUpdate: time.Now()}
		mock.
			ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category" WHERE "category"."category_id" = $1 ORDER BY "category"."category_id" LIMIT 1`)).
			WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"category_id", "name", "last_update"}).
			AddRow(mockedCategory.ID, mockedCategory.Name, mockedCategory.LastUpdate))

		mock.ExpectBegin()
		mock.
			ExpectExec(regexp.QuoteMeta(`UPDATE "category" SET "name"=$1,"last_update"=$2 WHERE "category_id" = $3`)).
			WithArgs(data["Name"], AnyTime{}, mockedCategory.ID).
			WillReturnError(expectedError)
		mock.ExpectRollback()

		err = model.UpdateById(ctx, db)

		assert.ErrorIs(t, expectedError, err)
	})

	t.Run("Should retrun a modified category", func(t *testing.T) {
		var model = &CategoryModel{}
		var db = GetDBTestConnection()
		var ctx = MockContext()
		ctx.AddParam("id", "1")
		ctx.Request = &http.Request{}

		db.Exec(`TRUNCATE TABLE category RESTART IDENTITY CASCADE`)
		db.Create(&Category{Name: "Drama"})

		data := map[string]string{"Name": "Action"}
		jsonData, err := json.Marshal(data)
		assert.NoError(t, err)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonData))

		err = model.UpdateById(ctx, db)

		assert.NoError(t, err)
		assert.Equal(t, data["Name"], model.Categories[0].Name)

		db.Delete(&Category{ID: 1})
	})
}

func TestDeleteById(t *testing.T) {

	t.Run("Must return an error if the Category with the id doesn't exist", func(t *testing.T) {
		var model = &CategoryModel{}
		var db = GetDBTestConnection()
		var ctx = MockContext()
		ctx.AddParam("id", "1")

		err := model.DeleteById(ctx, db)

		assert.ErrorContains(t, err, "record not found")
	})

	t.Run("Must return an error if the database fail to delete the category", func(t *testing.T) {
		var model = &CategoryModel{}
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

		mockedCategory := Category{ID: 1, Name: "Mocked Category Row", LastUpdate: time.Now()}
		mock.
			ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "category" WHERE "category"."category_id" = $1 ORDER BY "category"."category_id" LIMIT 1`)).
			WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"category_id", "name", "last_update"}).
			AddRow(mockedCategory.ID, mockedCategory.Name, mockedCategory.LastUpdate))

		mock.ExpectBegin()
		mock.
			ExpectExec(regexp.QuoteMeta(`DELETE FROM "category" WHERE "category"."category_id" = $1`)).
			WithArgs(mockedCategory.ID).
			WillReturnError(expectedError)
		mock.ExpectRollback()

		ctx.AddParam("id", "1")

		err = model.DeleteById(ctx, db)

		assert.ErrorIs(t, expectedError, err)
	})

	t.Run("Must delete a category when called", func(t *testing.T) {
		var model = &CategoryModel{}
		var db = GetDBTestConnection()
		var ctx = MockContext()
		ctx.AddParam("id", "1")

		db.Exec(`TRUNCATE TABLE category RESTART IDENTITY CASCADE`)
		db.Create(&Category{Name: "Action"})
		err := model.DeleteById(ctx, db)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), model.Categories[0].ID)
		assert.Equal(t, "Action", model.Categories[0].Name)

		db.Delete(&Category{ID: 1})
	})
}
