package model

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Case struct {
	Title   string
	Input   []Language
	Error   string
	Context []gin.Param
}

func init() {
	var err error
	var dbUser string = os.Getenv("TEST_DB_USER")
	var dbPass string = os.Getenv("TEST_DB_PASSWORD")
	var dbName string = os.Getenv("TEST_DB_NAME")
	var dbHost string = os.Getenv("TEST_DB_HOST")

	fmt.Println("Initilizing database connection")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5433 sslmode=disable Timezone=America/Sao_Paulo", dbHost, dbUser, dbPass, dbName)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Connection was successfuly")

	fmt.Print("Migrating Model Language to Test Database\n")
	if err := DB.AutoMigrate(&Language{}); err != nil {
		fmt.Errorf("Fail to migrate the model Language: %v\n", err)
	}
	fmt.Printf("Migration succeeded\n")
}

func MockContext() *gin.Context {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	return c
}
func TestGetLanguageById(t *testing.T) {
	var wants = []Case{
		{Title: "Empty database should return an error", Input: []Language{}, Error: "record not found", Context: []gin.Param{{Key: "id", Value: "1"}}},
		{Title: "Passing an invalid id should return an error", Input: []Language{}, Error: "invalid id, make sure to pass a number", Context: []gin.Param{{Key: "id", Value: "a"}}},
		{Title: "Should return an object Language when callled with appropriate id", Input: []Language{{Name: "English"}}, Context: []gin.Param{{Key: "id", Value: "1"}}},
	}
	for _, want := range wants {
		t.Run(want.Title, func(t *testing.T) {
			ctx := MockContext()
			ctx.Params = want.Context

			DB.Exec("TRUNCATE TABLE language RESTART IDENTITY;")

			for i := range want.Input {
				DB.Create(&want.Input[i])
			}

			lang, err := GetLanguageById(ctx)

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
	var wants = []Case{
		{Title: "Shouldn't return error when called with an empty DB", Input: []Language{}, Context: []gin.Param{{Key: "p", Value: "0"}}},
		{Title: "Should return an array when called with an non-empty DB", Input: []Language{{Name: "English"}, {Name: "Portuguese"}}, Context: []gin.Param{{Key: "p", Value: "0"}}},
	}
	for _, want := range wants {
		t.Run(want.Title, func(t *testing.T) {
			ctx := MockContext()
			ctx.Params = want.Context

			DB.Exec("TRUNCATE TABLE language RESTART IDENTITY;")

			for i := range want.Input {
				DB.Create(&want.Input[i])
			}

			langs, err := GetAllLanguages(ctx)

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
