package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	var db = GetDBTestConnection()
	fmt.Print("Migrating the Model Category to the Test Database\n")
	if err := db.AutoMigrate(&Category{}); err != nil {
		fmt.Errorf("Fail to migrate the model Category: %v\n", err)
	}
	fmt.Printf("Migration has been successful\n")
}

func TestGetCategoryById(t *testing.T) {
	var db = GetDBTestConnection()
	var wants = []Case[Category]{
		{Title: "Should return an error when called with empty database",
			Input:   []Category{},
			Error:   "record not found",
			Context: []gin.Param{{Key: "id", Value: "1"}}},
		{Title: "Should return an error when passing an invalid ID",
			Input:   []Category{},
			Error:   "invalid id, make sure to pass a number",
			Context: []gin.Param{{Key: "id", Value: "a"}}},
		{Title: "Should return an Category{} when called with a valid ID",
			Input:   []Category{{Name: "Drama"}},
			Context: []gin.Param{{Key: "id", Value: "1"}}},
	}
	for _, want := range wants {
		t.Run(want.Title, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			db.Exec("TRUNCATE TABLE category RESTART IDENTITY CASCADE;")

			for i := range want.Input {
				db.Create(&want.Input[i])
			}

			categor, err := GetById(ctx, db)

			if err != nil && err.Error() != want.Error {
				t.Errorf("The error fail to meet the expectation, want: %s, got: %s", want.Error, err.Error())
			}
			for _, wantCategory := range want.Input {
				if categor != nil && categor.Name != wantCategory.Name {
					t.Errorf("The return fail to meet the expectation, want: %v, got: %v", wantCategory, categor)
				}
			}

		})
	}
}

func TestGetAllCategories(t *testing.T) {
	var db = GetDBTestConnection()
	var wants = []Case[Category]{
		{Title: "Shouldn't return error when called (empty database)",
			Context: []gin.Param{{Key: "p", Value: "0"}}},
		{Title: "Should return an []Category{} when called (non-empty database)",
			Input:   []Category{{Name: "Drama"}, {Name: "Action"}},
			Context: []gin.Param{{Key: "p", Value: "0"}}},
	}
	for _, want := range wants {
		t.Run(want.Title, func(t *testing.T) {
			ctx := MockContext()
			ctx.Params = want.Context

			db.Exec("TRUNCATE TABLE category RESTART IDENTITY CASCADE;")

			for i := range want.Input {
				db.Create(&want.Input[i])
			}

			categories, err := GetAll(ctx, db)

			if err != nil && err.Error() != want.Error {
				t.Errorf("The error fail to meet the expectation, want: %s, got: %s", want.Error, err.Error())
			}
			for index, wantCategory := range want.Input {
				if categories != nil && (*categories)[index].Name != wantCategory.Name {
					t.Errorf("The return fail to meet the expectation, want: %v, got: %v", wantCategory, categories)
				}
			}

		})
	}
}

func TestCreateNewCategory(t *testing.T) {
	var db = GetDBTestConnection()
	var wants = []Case[Category]{
		{Title: "Should return an error when called with body empty",
			Error: "invalid request"},
		{Title: "Should return an error when called with an duplicate ID",
			Input: []Category{{ID: 1, Name: "Drama"}},
			Error: "ERROR: duplicate key value violates unique constraint \"category_pkey\" (SQLSTATE 23505)"},
		{Title: "Should return a Category{} when called",
			Input: []Category{{Name: "Action"}}},
	}

	db.Exec("TRUNCATE TABLE category RESTART IDENTITY CASCADE;")
	db.Create(&Category{Name: "Drama"})

	for _, want := range wants {
		t.Run(want.Title, func(t *testing.T) {

			ctx := MockContext()
			ctx.Request = &http.Request{Header: http.Header{}}
			ctx.Request.Header.Set("content-type", "application/json")

			if want.Input != nil {
				jsonBytes, _ := json.Marshal(want.Input[0])
				ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
			}

			got, err := CreateNew(ctx, db)

			if err != nil && err.Error() != want.Error {
				t.Errorf("The error fail to meet the expectation, want: %s, got: %s", want.Error, err.Error())
			}

			if got != nil && got.Name != want.Input[0].Name {
				t.Errorf("The return fail to meet the expectation, want: %v, got: %v", want.Input[0].Name, got)
			}

		})
	}
}

func TestUpdateCategory(t *testing.T) {
	var db = GetDBTestConnection()
	var wants = []Case[Category]{
		{Title: "Should return an error when called with a non-existent ID",
			Context: []gin.Param{{Key: "id", Value: "2"}},
			Error:   "record not found",
		},
		{Title: "Shoud return error when called with an invalid ID",
			Context: []gin.Param{{Key: "id", Value: "a"}},
			Error:   "invalid id, make sure to pass a number",
		},
		{Title: "Shoud return error when called with an invalid Category{}",
			Input:   nil,
			Context: []gin.Param{{Key: "id", Value: "1"}},
			Error:   "invalid request",
		},
		{Title: "Shoud return a Category{} updated when called",
			Input:   []Category{{Name: "Action"}},
			Context: []gin.Param{{Key: "id", Value: "1"}},
		},
	}

	for _, want := range wants {
		db.Exec("TRUNCATE TABLE category RESTART IDENTITY CASCADE;")
		db.Create(&Category{Name: "Drama"})
		t.Run(want.Title, func(t *testing.T) {

			ctx := MockContext()
			ctx.Params = want.Context
			ctx.Request = &http.Request{Header: http.Header{}}
			ctx.Request.Header.Set("content-type", "application/json")

			if want.Input != nil {
				jsonBytes, _ := json.Marshal(want.Input[0])
				ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
			}

			got, err := UpdateCategoryById(ctx, db)

			if err != nil && err.Error() != want.Error {
				t.Errorf("The error fail to meet the expectation, want: %s, got: %s", want.Error, err.Error())
			}

			if got != nil && got.Name != want.Input[0].Name {
				t.Errorf("The return fail to meet the expectation, want: %v, got: %v", want.Input[0].Name, got)
			}
		})
	}
}

func TestDeleteCategoryById(t *testing.T) {
	var db = GetDBTestConnection()
	var wants = []Case[Category]{
		{Title: "Should return an error when called with an invalid ID",
			Context: []gin.Param{{Key: "id", Value: "a"}},
			Error:   "invalid id, make sure to pass a number",
		},
		{Title: "Should return an error when called with a non-existent ID",
			Context: []gin.Param{{Key: "id", Value: "1"}},
			Error:   "record not found",
		},
		{Title: "Should return the deleted Category{} when called",
			Input:   []Category{{Name: "Drama"}},
			Context: []gin.Param{{Key: "id", Value: "1"}},
		},
	}
	for _, want := range wants {
		db.Exec("TRUNCATE TABLE category RESTART IDENTITY CASCADE;")
		if want.Input != nil {
			db.Create(&want.Input)
		}
		t.Run(want.Title, func(t *testing.T) {
			ctx := MockContext()
			ctx.Params = want.Context

			got, err := DeleteCategoryById(ctx, db)

			if err != nil && err.Error() != want.Error {
				t.Errorf("The error fail to meet the expectation, want: %s, got: %s", want.Error, err.Error())
			}

			if got != nil && got.Name != want.Input[0].Name {
				t.Errorf("The return fail to meet the expectation, want: %v, got: %v", want.Input[0].Name, got)
			}
		})
	}
}
