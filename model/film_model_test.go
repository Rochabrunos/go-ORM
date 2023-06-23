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

var mockedFilms = []Film{
	{
		Title:           "Chamber Italian",
		Description:     "A Fateful Reflection of a Moose And a Husband who must Overcome a Monkey in Nigeria",
		ReleaseYear:     "2006",
		RentalDuration:  7,
		RentalRate:      4.99,
		Length:          117,
		ReplacementCost: 14.99,
		Rating:          "NC-17",
		SpecialFeatures: []string{"Trailers"},
		FullText:        `'chamber':1 'fate':4 'husband':11 'italian':2 'monkey':16 'moos':8 'must':13 'nigeria':18 'overcom':14 'reflect':5`,
	},
	{
		Title:           "Grosse Wonderful",
		Description:     "A Epic Drama of a Cat And a Explorer who must Redeem a Moose in Australia ",
		ReleaseYear:     "2006",
		RentalDuration:  5,
		RentalRate:      4.99,
		Length:          49,
		ReplacementCost: 19.99,
		Rating:          "R",
		SpecialFeatures: []string{"Behind the Scenes"},
		FullText:        `'australia':18 'cat':8 'drama':5 'epic':4 'explor':11 'gross':1 'moos':16 'must':13 'redeem':14 'wonder':2`,
	},
}

var mockedLanguage = Language{Name: "English"}

func init() {
	var db = GetDBTestConnection()

	db.AutoMigrate(&Language{})
	if err := db.AutoMigrate(&Film{}); err != nil {
		fmt.Errorf("Fail to migrate the model Film: %v\n", err)
		panic(err)
	}
	db.Migrator().CreateConstraint(&Film{}, "Language")

	db.Model(&Language{}).Exec("TRUNCATE TABLE language RESTART IDENTITY CASCADE")
	db.Create(&mockedLanguage)
	mockedFilms[0].Language = mockedLanguage
	mockedFilms[1].Language = mockedLanguage
	fmt.Printf("Migration has been successful\n")
}

func TestFilmGetById(t *testing.T) {
	t.Run("Must return an error when called with an invalid ID", func(t *testing.T) {
		var model = &FilmModel{}
		var db = GetDBTestConnection()
		var ctx = MockContext()
		ctx.AddParam("id", "A")

		err := model.GetById(ctx, db)

		assert.ErrorContains(t, err, "invalid id, make sure to pass a number")
	})

	t.Run("Must return an error if database fails", func(t *testing.T) {
		var model = &FilmModel{}
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
			ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "film" WHERE "film"."film_id" = $1 ORDER BY "film"."film_id" LIMIT 1`)).
			WithArgs(1).WillReturnError(expectedError)

		ctx.AddParam("id", "1")

		err = model.GetById(ctx, db)

		assert.ErrorIs(t, err, expectedError)
	})

	t.Run("Must return a film with the specified ID", func(t *testing.T) {
		var model = &FilmModel{}
		var film = mockedFilms[0]
		var db = GetDBTestConnection()
		var ctx = MockContext()
		ctx.AddParam("id", "1")

		db.Exec("TRUNCATE TABLE film RESTART IDENTITY CASCADE")
		db.Create(&film)

		err := model.GetById(ctx, db)

		assert.NoError(t, err)
		assert.Equal(t, film.ID, model.Films[0].ID)

		db.Delete(&Film{ID: 1})
	})
}

func TestFilmGetAll(t *testing.T) {

	//Following the logic of the function, the only way to return an error is if the db.Find function fails.
	//This test mocks a failure in order to verify the correctness of the entire function in the event of an error
	t.Run("Must return an error if the database fails to retrieve objects", func(t *testing.T) {
		var model = &FilmModel{}
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

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"film\" LIMIT 10")).WithArgs().WillReturnError(expectedError)

		err = model.GetAll(ctx, db)

		assert.Error(t, err, expectedError)
	})

	t.Run("Must return the film store in the database", func(t *testing.T) {
		var model = &FilmModel{}
		var ctx = MockContext()
		var db = GetDBTestConnection()
		db.Exec("TRUNCATE TABLE film RESTART IDENTITY CASCADE;")
		db.Create([]Film{mockedFilms[0], mockedFilms[1]})

		err := model.GetAll(ctx, db)

		assert.NoError(t, err)
		assert.NotNil(t, model.Films)
		assert.Equal(t, 2, len(model.Films))

		db.Delete(&Film{ID: 1})
		db.Delete(&Film{ID: 2})
	})
}

func TestFilmCreateNew(t *testing.T) {
	t.Run("Must return an error if called with invalid data", func(t *testing.T) {
		var model = &FilmModel{}
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

	t.Run("Must return an error if it's not possible to create the new film", func(t *testing.T) {
		var model = &FilmModel{}
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

		data := mockedFilms[0]
		jsonBytes, _ := json.Marshal(data)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
		mock.ExpectBegin()
		mock.
			ExpectQuery(regexp.QuoteMeta(`INSERT INTO "language" ("name","last_update","language_id") VALUES ($1,$2,$3) ON CONFLICT DO NOTHING RETURNING "language_id"`)).
			WithArgs(data.Language.Name, AnyTime{}, data.Language.ID).
			WillReturnRows(sqlmock.NewRows([]string{"name", "last_update", "language_id"}).AddRow("English", time.Time{}, 1))
		mock.ExpectQuery(
			regexp.QuoteMeta(`INSERT INTO "film" ("title","description","release_year","language_id","rental_duration","rental_rate","length","replacement_cost","rating","last_update","special_features","full_text") 
								VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING "film_id"`)).
			WithArgs(data.Title, data.Description, data.ReleaseYear, data.Language.ID,
				data.RentalDuration, data.RentalRate, data.Length, data.ReplacementCost,
				data.Rating, AnyTime{}, data.SpecialFeatures, data.FullText).
			WillReturnError(expectedError)
		mock.ExpectRollback()

		err = model.CreateNew(ctx, db)

		assert.ErrorIs(t, err, expectedError)
	})

	t.Run("Must return a new film when called", func(t *testing.T) {
		var model = &FilmModel{}
		var db = GetDBTestConnection()
		var ctx = MockContext()
		ctx.Request = &http.Request{Header: http.Header{}}
		ctx.Request.Header.Set("content-type", "application/json")

		db.Exec("TRUNCATE TABLE film RESTART IDENTITY CASCADE;")

		data := mockedFilms[0]
		jsonBytes, _ := json.Marshal(data)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))

		err := model.CreateNew(ctx, db)

		assert.NoError(t, err)
		assert.Equal(t, data.Title, model.Films[0].Title)
	})
}

func TestFilmUpdateById(t *testing.T) {
	t.Run("Must return an error when the film doesn't exists", func(t *testing.T) {
		var model = &FilmModel{}
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

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "film" WHERE "film"."film_id" = $1 ORDER BY "film"."film_id" LIMIT 1`)).
			WithArgs(1).
			WillReturnError(expectedError)

		ctx.AddParam("id", "1")

		err = model.UpdateById(ctx, db)

		assert.ErrorIs(t, expectedError, err)
	})

	t.Run("Must return an error when called with invalid JSON body", func(t *testing.T) {
		var model = &FilmModel{}
		var db = GetDBTestConnection()
		var ctx = MockContext()
		ctx.AddParam("id", "1")
		ctx.Request = &http.Request{}

		db.Exec(`TRUNCATE TABLE film RESTART IDENTITY CASCADE`)
		result := db.Create([]Film{mockedFilms[0]})
		assert.NoError(t, result.Error)

		data := map[string]string{"Name": "Title"}
		jsonBytes, err := json.Marshal(data)
		assert.NoError(t, err)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))

		err = model.UpdateById(ctx, db)

		assert.ErrorContains(t, err, "Error:Field validation for 'Title'")

		db.Delete(&Film{ID: 1})
	})

	t.Run("Must return an error if the database fail to update the film", func(t *testing.T) {
		var model = &FilmModel{}
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
		data := mockedFilms[0]
		jsonBytes, err := json.Marshal(data)
		assert.NoError(t, err)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
		ctx.AddParam("id", "1")
		mock.
			ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "film" WHERE "film"."film_id" = $1 ORDER BY "film"."film_id" LIMIT 1`)).
			WithArgs(1).WillReturnRows(sqlmock.NewRows(
			[]string{"film_id", "title", "description", "release_year", "language_id", "rental_duration", "rental_rate", "length", "replacement_cost", "rating", "last_update", "special_features", "full_text"}).
			AddRow(1, mockedFilms[1].Title, mockedFilms[1].Description, mockedFilms[1].ReleaseYear, mockedFilms[1].LanguageID,
				mockedFilms[1].RentalDuration, mockedFilms[1].RentalRate, mockedFilms[1].Length, mockedFilms[1].ReplacementCost,
				mockedFilms[1].Rating, time.Time{}, mockedFilms[1].SpecialFeatures, mockedFilms[1].FullText))

		mock.ExpectBegin()
		mock.
			ExpectQuery(regexp.QuoteMeta(`INSERT INTO "language" ("name","last_update","language_id") VALUES ($1,$2,$3) ON CONFLICT DO NOTHING RETURNING "language_id"`)).
			WithArgs(data.Language.Name, AnyTime{}, data.Language.ID).
			WillReturnRows(sqlmock.NewRows([]string{"name", "last_update", "language_id"}).AddRow("English", time.Time{}, 1))
		mock.
			ExpectExec(regexp.QuoteMeta(`UPDATE "film" 
				SET "title"=$1,"description"=$2,"release_year"=$3,"language_id"=$4,"rental_duration"=$5,"rental_rate"=$6,"length"=$7,"replacement_cost"=$8,"rating"=$9,"last_update"=$10,"special_features"=$11,"full_text"=$12
				WHERE "film_id" = $13`)).
			WithArgs(data.Title, data.Description, data.ReleaseYear, data.Language.ID,
				data.RentalDuration, data.RentalRate, data.Length, data.ReplacementCost,
				data.Rating, AnyTime{}, data.SpecialFeatures, data.FullText, 1).
			WillReturnError(expectedError)
		mock.ExpectRollback()

		err = model.UpdateById(ctx, db)

		assert.ErrorIs(t, expectedError, err)
	})

	t.Run("Should retrun a modified film", func(t *testing.T) {
		var model = &FilmModel{}
		var db = GetDBTestConnection()
		var ctx = MockContext()
		ctx.AddParam("id", "1")
		ctx.Request = &http.Request{}

		db.Exec(`TRUNCATE TABLE film RESTART IDENTITY CASCADE`)
		db.Create(&mockedFilms[0])

		data := mockedFilms[1]
		jsonData, err := json.Marshal(data)
		assert.NoError(t, err)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonData))

		err = model.UpdateById(ctx, db)

		assert.NoError(t, err)
		assert.Equal(t, data.Title, model.Films[0].Title)

		db.Delete(&Film{ID: 1})
	})
}

func TestFilmDeleteById(t *testing.T) {

	t.Run("Must return an error if the Film with the id doesn't exist", func(t *testing.T) {
		var model = &FilmModel{}
		var db = GetDBTestConnection()
		var ctx = MockContext()
		ctx.AddParam("id", "1")

		err := model.DeleteById(ctx, db)

		assert.ErrorContains(t, err, "record not found")
	})

	t.Run("Must return an error if the database fail to delete the film", func(t *testing.T) {
		var model = &FilmModel{}
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

		mock.
			ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "film" WHERE "film"."film_id" = $1 ORDER BY "film"."film_id" LIMIT 1`)).
			WithArgs(1).WillReturnRows(sqlmock.NewRows(
			[]string{"film_id", "title", "description", "release_year", "language_id", "rental_duration", "rental_rate", "length", "replacement_cost", "rating", "last_update", "special_features", "full_text"}).
			AddRow(1, mockedFilms[1].Title, mockedFilms[1].Description, mockedFilms[1].ReleaseYear, mockedFilms[1].LanguageID,
				mockedFilms[1].RentalDuration, mockedFilms[1].RentalRate, mockedFilms[1].Length, mockedFilms[1].ReplacementCost,
				mockedFilms[1].Rating, time.Time{}, mockedFilms[1].SpecialFeatures, mockedFilms[1].FullText))

		mock.ExpectBegin()
		mock.
			ExpectExec(regexp.QuoteMeta(`DELETE FROM "film" WHERE "film"."film_id" = $1`)).
			WithArgs(1).
			WillReturnError(expectedError)
		mock.ExpectRollback()

		ctx.AddParam("id", "1")

		err = model.DeleteById(ctx, db)

		assert.ErrorIs(t, expectedError, err)
	})

	t.Run("Must delete a film when called", func(t *testing.T) {
		var model = &FilmModel{}
		var db = GetDBTestConnection()
		var ctx = MockContext()
		ctx.AddParam("id", "1")

		db.Exec(`TRUNCATE TABLE film RESTART IDENTITY CASCADE`)
		db.Create(&mockedFilms[0])
		err := model.DeleteById(ctx, db)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), model.Films[0].ID)
		assert.Equal(t, mockedFilms[0].Title, model.Films[0].Title)

		db.Delete(&Film{ID: 1})
	})
}
