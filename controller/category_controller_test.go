package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"orm-golang/model"
	mocks "orm-golang/model/mock"
	"orm-golang/service"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	db := GetDBTestConnection()
	db.Migrator().DropTable(&model.Category{})
	db.AutoMigrate(&model.Category{})
}
func TestGetByIdCategoryEndpoint(t *testing.T) {
	t.Run("Test failing to get a category", func(t *testing.T) {
		daoCategory = DAO{}.New(&mocks.MockedCategoryModel{}, GetDBTestConnection())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.AddParam("id", "1")

		GetByIdCategoryEndpoint(c)

		strJsonError, err := json.Marshal(gin.H{"error": mocks.MockedErrorMessage})
		assert.NoError(t, err)
		assert.Equal(t, w.Body.Bytes(), strJsonError)
		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

	})
	t.Run("Test return a category as expect", func(t *testing.T) {
		var result model.CategoryModel

		daoCategory = DAO{}.New(&model.CategoryModel{}, service.GetDBConnection())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.AddParam("id", "1")

		GetByIdCategoryEndpoint(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusOK, "status code fail to meet the expectation")
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err, "return fail to meet the expectation")
		assert.Equal(t, 1, len(result.Categories))
		assert.Equal(t, uint(1), result.Categories[0].ID)
	})
}

func TestGetAllCategoryEndpoint(t *testing.T) {
	t.Run("Test failing to get the categories", func(t *testing.T) {
		daoCategory = DAO{}.New(&mocks.MockedCategoryModel{}, GetDBTestConnection())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		GetAllCategoryEndpoint(c)

		strJsonError, err := json.Marshal(gin.H{"error": mocks.MockedErrorMessage})
		assert.NoError(t, err)
		assert.Equal(t, w.Body.Bytes(), strJsonError)
		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

	})
	t.Run("Test return the categories as expect", func(t *testing.T) {
		var result model.CategoryModel

		daoCategory = DAO{}.New(&model.CategoryModel{}, service.GetDBConnection())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		GetAllCategoryEndpoint(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusOK, "status code fail to meet the expectation")
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err, "return fail to meet the expectation")
	})
}

func TestCreateCategoryEndpoint(t *testing.T) {
	t.Run("Test fail to create the category", func(t *testing.T) {
		daoCategory = DAO{}.New(&mocks.MockedCategoryModel{}, GetDBTestConnection())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		CreateCategoryEndpoint(c)
		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest, "status code fail to meet the expectation")

		strJsonError, err := json.Marshal(gin.H{"error": mocks.MockedErrorMessage})
		assert.NoError(t, err)
		assert.Equal(t, w.Body.Bytes(), strJsonError)
	})

	t.Run("Test successfully creates the category", func(t *testing.T) {
		db := GetDBTestConnection()
		db.Exec("TRUNCATE TABLE category RESTART IDENTITY CASCADE")

		daoCategory = DAO{}.New(&model.CategoryModel{}, db)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		data := map[string]string{"Name": "Action"}
		jsonBytes, err := json.Marshal(data)
		assert.NoError(t, err)
		c.Request = &http.Request{}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))

		CreateCategoryEndpoint(c)
		assert.Equal(t, w.Result().StatusCode, http.StatusOK, "status code fail to meet the expectation")

		var result model.CategoryModel
		err = json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, data["Name"], result.Categories[0].Name)
	})
}

func TestModifyCategoryEndpoint(t *testing.T) {
	t.Run("Test fail to modify the category", func(t *testing.T) {
		daoCategory = DAO{}.New(&mocks.MockedCategoryModel{}, GetDBTestConnection())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		ModifyCategoryEndpoint(c)
		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest, "status code fail to meet the expectation")

		strJsonError, err := json.Marshal(gin.H{"error": mocks.MockedErrorMessage})
		assert.NoError(t, err)
		assert.Equal(t, w.Body.Bytes(), strJsonError)
	})

	t.Run("Test successfully modifies the category", func(t *testing.T) {
		db := GetDBTestConnection()
		db.Exec("TRUNCATE TABLE category RESTART IDENTITY CASCADE")
		db.Create(&model.Category{Name: "Action"})

		daoCategory = DAO{}.New(&model.CategoryModel{}, db)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.AddParam("id", "1")
		data := map[string]string{"Name": "Drama"}
		jsonBytes, err := json.Marshal(data)
		assert.NoError(t, err)
		c.Request = &http.Request{}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))

		CreateCategoryEndpoint(c)
		assert.Equal(t, w.Result().StatusCode, http.StatusOK, "status code fail to meet the expectation")

		var result model.CategoryModel
		err = json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, data["Name"], result.Categories[0].Name)
	})
}

func TestDeleteCategoryEndpoint(t *testing.T) {
	t.Run("Test fail to delete the category", func(t *testing.T) {
		daoCategory = DAO{}.New(&mocks.MockedCategoryModel{}, GetDBTestConnection())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		DeleteCategoryEndpoint(c)
		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest, "status code fail to meet the expectation")

		strJsonError, err := json.Marshal(gin.H{"error": mocks.MockedErrorMessage})
		assert.NoError(t, err)
		assert.Equal(t, w.Body.Bytes(), strJsonError)
	})

	t.Run("Test successfully deletes the category", func(t *testing.T) {
		db := GetDBTestConnection()
		db.Exec("TRUNCATE TABLE category RESTART IDENTITY CASCADE")
		categoryData := model.Category{Name: "Action"}
		db.Create(&categoryData)

		daoCategory = DAO{}.New(&model.CategoryModel{}, db)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.AddParam("id", "1")

		DeleteCategoryEndpoint(c)
		assert.Equal(t, w.Result().StatusCode, http.StatusOK, "status code fail to meet the expectation")

		var result model.CategoryModel
		err := json.Unmarshal(w.Body.Bytes(), &result)

		assert.NoError(t, err)
		assert.Equal(t, categoryData.Name, result.Categories[0].Name)
	})
}
