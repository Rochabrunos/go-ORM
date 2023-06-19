package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"orm-golang/model"
	mocks "orm-golang/model/mock"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetByIdEndpoint(t *testing.T) {
	t.Run("Test failing to get a category", func(t *testing.T) {
		dao = DAO{}.New(&mocks.MockedCategoryModel{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.AddParam("id", "1")
		GetByIdEndpoint(c)
		strJsonError, err := json.Marshal(gin.H{"error": mocks.MockedErrorMessage})
		assert.NoError(t, err)
		assert.Equal(t, w.Body.Bytes(), strJsonError)
		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

	})
	t.Run("Test return a category as expect", func(t *testing.T) {
		var result model.CategoryModel
		dao = DAO{}.New(&model.CategoryModel{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.AddParam("id", "1")
		GetByIdEndpoint(c)
		assert.Equal(t, w.Result().StatusCode, http.StatusOK, "status code fail to meet the expectation")
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err, "return fail to meet the expectation")
		assert.Equal(t, 1, len(result.Categories))
		assert.Equal(t, uint(1), result.Categories[0].ID)
	})
}

func TestGetAllEndpoint(t *testing.T) {
	t.Run("Test failing to get the categories", func(t *testing.T) {
		dao = DAO{}.New(&mocks.MockedCategoryModel{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		GetAllEndpoint(c)
		strJsonError, err := json.Marshal(gin.H{"error": mocks.MockedErrorMessage})
		assert.NoError(t, err)
		assert.Equal(t, w.Body.Bytes(), strJsonError)
		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

	})
	t.Run("Test return the categories as expect", func(t *testing.T) {
		var result model.CategoryModel
		dao = DAO{}.New(&model.CategoryModel{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		GetAllEndpoint(c)
		assert.Equal(t, w.Result().StatusCode, http.StatusOK, "status code fail to meet the expectation")
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err, "return fail to meet the expectation")
	})
}
