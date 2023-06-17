package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
)

func init() {
	var db = GetDBTestConnection()

	db.AutoMigrate(&Language{})
	if err := db.AutoMigrate(&Film{}); err != nil {
		fmt.Errorf("Fail to migrate the model Film: %v\n", err)
		panic(err)
	}
	db.Migrator().CreateConstraint(&Film{}, "Language")
	fmt.Printf("Migration has been successful\n")
}

func createMockFilm(id uint) Film {
	var f Film
	f.ID = id
	f.Title = uniuri.NewLen(50)
	f.Description = "An interest film"
	f.ReleaseYear = "2013"
	f.LanguageID = 1
	f.RentalDuration = 3
	f.RentalRate = 5.99
	f.Length = 120
	f.ReplacementCost = 50.99
	f.SpecialFeatures = []string{"special features", "behind the scene"}
	f.Rating = "G"
	f.FullText = ""

	return f
}

func TestGetFilmById(t *testing.T) {
	var db = GetDBTestConnection()
	var wants = []Case[Film]{
		{Title: "Should return an error when called with empty database",
			Error:   "record not found",
			Context: []gin.Param{{Key: "id", Value: "1"}}},
		{Title: "Should return an error when passing an invalid ID",
			Error:   "invalid id, make sure to pass a number",
			Context: []gin.Param{{Key: "id", Value: "a"}}},
		{Title: "Should return an Film{} when called with a valid ID",
			Input:   []Film{createMockFilm(0)},
			Context: []gin.Param{{Key: "id", Value: "1"}}},
	}
	db.Create(&Language{ID: 1, Name: "English"})
	for _, want := range wants {
		t.Run(want.Title, func(t *testing.T) {
			ctx := MockContext()
			ctx.Params = want.Context

			db.Exec("TRUNCATE TABLE film RESTART IDENTITY CASCADE;")
			if want.Input != nil {
				for i := range want.Input {
					db.Create(&want.Input[i])
				}
			}

			film, err := GetFilmById(ctx, db)
			if err != nil && err.Error() != want.Error {
				t.Errorf("The error fail to meet the expectation, want: %s, got: %s", want.Error, err.Error())
			}
			for _, wantCategory := range want.Input {
				if film != nil && film.Title != wantCategory.Title {
					t.Errorf("The return fail to meet the expectation, want: %v, got: %v", wantCategory, film)
				}
			}

		})
	}
	db.Delete(&Language{ID: 1, Name: "English"})
}

func TestGetAllFilms(t *testing.T) {
	var db = GetDBTestConnection()
	var wants = []Case[Film]{
		{Title: "Shouldn't return error when called (empty database)",
			Context: []gin.Param{{Key: "p", Value: "0"}}},
		{Title: "Should return an []Film{} when called (non-empty database)",
			Input:   []Film{createMockFilm(0), createMockFilm(0)},
			Context: []gin.Param{{Key: "p", Value: "0"}}},
	}
	db.Create(&Language{ID: 1, Name: "English"})
	for _, want := range wants {
		t.Run(want.Title, func(t *testing.T) {
			ctx := MockContext()
			ctx.Params = want.Context

			db.Exec("TRUNCATE TABLE film RESTART IDENTITY CASCADE;")

			for i := range want.Input {
				db.Create(&want.Input[i])
			}

			films, err := GetAllFilms(ctx, db)

			if err != nil && err.Error() != want.Error {
				t.Errorf("The error fail to meet the expectation, want: %s, got: %s", want.Error, err.Error())
			}
			for index, wantCategory := range want.Input {
				if films != nil && (*films)[index].Title != wantCategory.Title {
					t.Errorf("The return fail to meet the expectation, want: %v, got: %v", wantCategory, films)
				}
			}

		})
	}
	db.Delete(&Language{ID: 1, Name: "English"})
}

func TestCreateNewFilm(t *testing.T) {
	var db = GetDBTestConnection()
	var wants = []Case[Film]{
		{Title: "Should return an error when called with body empty",
			Error: "invalid request"},
		{Title: "Should return an error when called with an duplicate ID",
			Input: []Film{createMockFilm(1)},
			Error: "ERROR: duplicate key value violates unique constraint \"film_pkey\" (SQLSTATE 23505)"},
		{Title: "Should return a Film{} when called",
			Input: []Film{createMockFilm(0)}},
	}

	db.Exec("TRUNCATE TABLE film RESTART IDENTITY CASCADE;")
	db.Create(&Language{ID: 1, Name: "English"})
	film := createMockFilm(0)
	db.Create(&film)

	for _, want := range wants {
		t.Run(want.Title, func(t *testing.T) {

			ctx := MockContext()
			ctx.Request = &http.Request{Header: http.Header{}}
			ctx.Request.Header.Set("content-type", "application/json")

			if want.Input != nil {
				jsonBytes, _ := json.Marshal(want.Input[0])
				ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
			}

			got, err := CreateNewFilm(ctx, db)

			if err != nil && err.Error() != want.Error {
				t.Errorf("The error fail to meet the expectation, want: %s, got: %s", want.Error, err.Error())
			}

			if got != nil && got.Title != want.Input[0].Title {
				t.Errorf("The return fail to meet the expectation, want: %v, got: %v", want.Input[0].Title, got.Title)
			}

		})
	}
	db.Delete(&Language{ID: 1, Name: "English"})
}
func TestUpdateFilm(t *testing.T) {
	var db = GetDBTestConnection()
	var wants = []Case[Film]{
		{Title: "Should return an error when called with a non-existent ID",
			Context: []gin.Param{{Key: "id", Value: "2"}},
			Error:   "record not found",
		},
		{Title: "Shoud return error when called with an invalid ID",
			Context: []gin.Param{{Key: "id", Value: "a"}},
			Error:   "invalid id, make sure to pass a number",
		},
		{Title: "Shoud return error when called with an invalid Film{}",
			Input:   nil,
			Context: []gin.Param{{Key: "id", Value: "1"}},
			Error:   "invalid request",
		},
		{Title: "Shoud return a Film{} updated when called",
			Input:   []Film{createMockFilm(1)},
			Context: []gin.Param{{Key: "id", Value: "1"}},
		},
	}

	for _, want := range wants {
		db.Exec("TRUNCATE TABLE film RESTART IDENTITY CASCADE;")
		db.Create(&Language{ID: 1, Name: "English"})
		film := createMockFilm(0)
		db.Create(&film)
		t.Run(want.Title, func(t *testing.T) {

			ctx := MockContext()
			ctx.Params = want.Context
			ctx.Request = &http.Request{Header: http.Header{}}
			ctx.Request.Header.Set("content-type", "application/json")

			if want.Input != nil {
				jsonBytes, _ := json.Marshal(want.Input[0])
				ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
			}

			got, err := UpdateFilmById(ctx, db)

			if err != nil && err.Error() != want.Error {
				t.Errorf("The error fail to meet the expectation, want: %s, got: %s", want.Error, err.Error())
			}

			if got != nil && got.Title != want.Input[0].Title {
				t.Errorf("The return fail to meet the expectation, want: %v, got: %v", want.Input[0].Title, got.Title)
			}
		})
	}
	db.Delete(&Language{ID: 1, Name: "English"})
}

func TestDeleteFIlmById(t *testing.T) {
	var db = GetDBTestConnection()
	var wants = []Case[Film]{
		{Title: "Should return an error when called with an invalid ID",
			Context: []gin.Param{{Key: "id", Value: "a"}},
			Error:   "invalid id, make sure to pass a number",
		},
		{Title: "Should return an error when called with a non-existent ID",
			Context: []gin.Param{{Key: "id", Value: "1"}},
			Error:   "record not found",
		},
		{Title: "Should return the deleted Film{} when called",
			Input:   []Film{createMockFilm(0)},
			Context: []gin.Param{{Key: "id", Value: "1"}},
		},
	}
	db.Create(&Language{ID: 1, Name: "English"})
	for _, want := range wants {
		db.Exec("TRUNCATE TABLE film  RESTART IDENTITY CASCADE;")
		if want.Input != nil {
			db.Create(&want.Input)
		}
		t.Run(want.Title, func(t *testing.T) {
			ctx := MockContext()
			ctx.Params = want.Context

			got, err := DeleteFilmById(ctx, db)

			if err != nil && err.Error() != want.Error {
				t.Errorf("The error fail to meet the expectation, want: %s, got: %s", want.Error, err.Error())
			}

			if got != nil && got.Title != want.Input[0].Title {
				t.Errorf("The return fail to meet the expectation, want: %v, got: %v", want.Input[0].Title, got.Title)
			}
		})
	}
	db.Delete(&Language{ID: 1, Name: "English"})
}
