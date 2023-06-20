package model

import (
	"bytes"
	"database/sql"
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

func init() {
	var db = GetDBTestConnection()
	fmt.Print("Migrating the Model Category to the Test Database\n")
	if err := db.AutoMigrate(&Category{}); err != nil {
		fmt.Errorf("Fail to migrate the model Category: %v\n", err)
	}
	db.Exec("TRUNCATE TABLE category RESTART IDENTITY CASCADE;")
	fmt.Printf("Migration has been successful\n")
}

func TestGetAll(t *testing.T) {

	//Following the logic of the function, the only way to return an error is if the db.Find function fails.
	//This test mocks a failure in order to verify the correctness of the entire function in the event of an error
	t.Run("Must return an error if the database fails to retrieve objects", func(t *testing.T) {
		var model = &CategoryModel{}
		var want = errors.New(mocks.MockedErrorMessage)
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

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"category\" LIMIT 10")).WithArgs().WillReturnError(want)

		err = model.GetAll(ctx, db)

		assert.Error(t, err, want)
	})

	t.Run("Must return the categories store in the database", func(t *testing.T) {
		var model = &CategoryModel{}
		var db = GetDBTestConnection()
		db.Exec("TRUNCATE TABLE category RESTART IDENTITY CASCADE;")
		db.Create([]Category{{Name: "Action"}, {Name: "Drama"}})

		ctx := MockContext()

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
		model := &CategoryModel{}
		db := GetDBTestConnection()

		ctx := MockContext()
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
		var want = errors.New(mocks.MockedErrorMessage)
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
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "category" ("name","last_update") VALUES ($1,$2) RETURNING "category_id"`)).WithArgs("Action", time.Now().Local()).WillReturnError(want)
		mock.ExpectRollback()

		err = model.CreateNew(ctx, db)

		assert.ErrorIs(t, err, want)
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
