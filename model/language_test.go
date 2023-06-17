package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func init() {
	var err error
	var dbUser string = os.Getenv("TEST_DB_USER")
	var dbPass string = os.Getenv("TEST_DB_PASSWORD")
	var dbName string = os.Getenv("TEST_DB_NAME")
	var dbHost string = os.Getenv("TEST_DB_HOST")

	fmt.Println("Initilizing database connection")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5433 sslmode=disable Timezone=America/Sao_Paulo", dbHost, dbUser, dbPass, dbName)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Connection was successful")

	fmt.Print("Migrating the Model Language to the Test Database\n")
	if err := DB.AutoMigrate(&Language{}); err != nil {
		fmt.Errorf("Fail to migrate the model Language: %v\n", err)
	}
	fmt.Printf("Migration has been successful\n")
}
func TestGetLanguageById(t *testing.T) {
	var wants = []Case[Language]{
		{Title: "Should return an error when called with empty database",
			Input:   []Language{},
			Error:   "record not found",
			Context: []gin.Param{{Key: "id", Value: "1"}}},
		{Title: "Should return an error when passing an invalid ID",
			Input:   []Language{},
			Error:   "invalid id, make sure to pass a number",
			Context: []gin.Param{{Key: "id", Value: "a"}}},
		{Title: "Should return an Language{} when called with a valid ID",
			Input:   []Language{{Name: "English"}},
			Context: []gin.Param{{Key: "id", Value: "1"}}},
	}
	for _, want := range wants {
		t.Run(want.Title, func(t *testing.T) {
			ctx := MockContext()
			ctx.Params = want.Context

			DB.Exec("TRUNCATE TABLE language RESTART IDENTITY CASCADE;")

			for i := range want.Input {
				DB.Create(&want.Input[i])
			}

			lang, err := GetLanguageById(ctx, DB)

			if err != nil && err.Error() != want.Error {
				t.Errorf("The error fail to meet the expectation, want: %s, got: %s", want.Error, err.Error())
			}
			for _, wantLanguage := range want.Input {
				if lang != nil && lang.Name != wantLanguage.Name {
					t.Errorf("The return fail to meet the expectation, want: %v, got: %v", wantLanguage, lang)
				}
			}

		})
	}
}

func TestGetAllLanguages(t *testing.T) {
	var wants = []Case[Language]{
		{Title: "Shouldn't return error when called (empty database)",
			Context: []gin.Param{{Key: "p", Value: "0"}}},
		{Title: "Should return an []Language{} when called (non-empty database)",
			Input:   []Language{{Name: "English"}, {Name: "Portuguese"}},
			Context: []gin.Param{{Key: "p", Value: "0"}}},
	}
	for _, want := range wants {
		t.Run(want.Title, func(t *testing.T) {
			ctx := MockContext()
			ctx.Params = want.Context

			DB.Exec("TRUNCATE TABLE language RESTART IDENTITY CASCADE;")

			for i := range want.Input {
				DB.Create(&Language{Name: want.Input[i].Name})
			}

			langs, err := GetAllLanguages(ctx, DB)

			if err != nil && err.Error() != want.Error {
				t.Errorf("The error fail to meet the expectation, want: %s, got: %s", want.Error, err.Error())
			}
			for index, wantLanguage := range want.Input {
				if langs != nil && (*langs)[index].Name != wantLanguage.Name {
					t.Errorf("The return fail to meet the expectation, want: %v, got: %v", wantLanguage, langs)
				}
			}

		})
	}
}

func TestCreateNewLanguage(t *testing.T) {
	var wants = []Case[Language]{
		{Title: "Should return an error when called with body empty",
			Error: "invalid request"},
		{Title: "Should return an error when called with an duplicate ID",
			Input: []Language{{ID: 1, Name: "Portuguese"}},
			Error: "ERROR: duplicate key value violates unique constraint \"language_pkey\" (SQLSTATE 23505)"},
		{Title: "Should return a Language{} when called",
			Input: []Language{{Name: "Portuguese"}}},
	}

	DB.Exec("TRUNCATE TABLE language RESTART IDENTITY CASCADE;")
	DB.Create(&Language{Name: "English"})

	for _, want := range wants {
		t.Run(want.Title, func(t *testing.T) {

			ctx := MockContext()
			ctx.Request = &http.Request{Header: http.Header{}}
			ctx.Request.Header.Set("content-type", "application/json")

			if want.Input != nil {
				jsonBytes, _ := json.Marshal(want.Input[0])
				ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
			}

			got, err := CreateNewLanguage(ctx, DB)
			if err != nil && err.Error() != want.Error {
				t.Errorf("The error fail to meet the expectation, want: %s, got: %s", want.Error, err.Error())
			}

			if got != nil && got.Name != want.Input[0].Name {
				t.Errorf("The return fail to meet the expectation, want: %v, got: %v", want.Input[0].Name, got)
			}

		})
	}
}

func TestUpdateLanguage(t *testing.T) {
	var wants = []Case[Language]{
		{Title: "Should return an error when called with a non-existent ID",
			Context: []gin.Param{{Key: "id", Value: "2"}},
			Error:   "record not found",
		},
		{Title: "Shoud return error when called with an invalid ID",
			Context: []gin.Param{{Key: "id", Value: "a"}},
			Error:   "invalid id, make sure to pass a number",
		},
		{Title: "Shoud return error when called with an invalid Language{}",
			Input:   nil,
			Context: []gin.Param{{Key: "id", Value: "1"}},
			Error:   "invalid request",
		},
		{Title: "Shoud return a Language{} updated when called",
			Input:   []Language{{Name: "Portuguse"}},
			Context: []gin.Param{{Key: "id", Value: "1"}},
		},
	}

	for _, want := range wants {
		DB.Exec("TRUNCATE TABLE language RESTART IDENTITY CASCADE;")
		DB.Create(&Language{Name: "English"})
		t.Run(want.Title, func(t *testing.T) {

			ctx := MockContext()
			ctx.Params = want.Context
			ctx.Request = &http.Request{Header: http.Header{}}
			ctx.Request.Header.Set("content-type", "application/json")

			if want.Input != nil {
				jsonBytes, _ := json.Marshal(want.Input[0])
				ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
			}

			got, err := UpdateLanguageById(ctx, DB)

			if err != nil && err.Error() != want.Error {
				t.Errorf("The error fail to meet the expectation, want: %s, got: %s", want.Error, err.Error())
			}

			if got != nil && got.Name != want.Input[0].Name {
				t.Errorf("The return fail to meet the expectation, want: %v, got: %v", want.Input[0].Name, got)
			}
		})
	}
}

func TestDeleteLanguageById(t *testing.T) {
	var wants = []Case[Language]{
		{Title: "Should return an error when called with an invalid ID",
			Context: []gin.Param{{Key: "id", Value: "a"}},
			Error:   "invalid id, make sure to pass a number",
		},
		{Title: "Should return an error when called with a non-existent ID",
			Context: []gin.Param{{Key: "id", Value: "1"}},
			Error:   "record not found",
		},
		{Title: "Should return the deleted Language{} when called",
			Input:   []Language{{Name: "English"}},
			Context: []gin.Param{{Key: "id", Value: "1"}},
		},
	}
	for _, want := range wants {
		DB.Exec("TRUNCATE TABLE language RESTART IDENTITY CASCADE;")
		if want.Input != nil {
			DB.Create(&want.Input)
		}
		t.Run(want.Title, func(t *testing.T) {
			ctx := MockContext()
			ctx.Params = want.Context

			got, err := DeleteLanguageById(ctx, DB)

			if err != nil && err.Error() != want.Error {
				t.Errorf("The error fail to meet the expectation, want: %s, got: %s", want.Error, err.Error())
			}

			if got != nil && got.Name != want.Input[0].Name {
				t.Errorf("The return fail to meet the expectation, want: %v, got: %v", want.Input[0].Name, got)
			}
		})
	}
}
