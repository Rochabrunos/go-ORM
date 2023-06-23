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

func TestLanguageGetByIdEndpoint(t *testing.T) {
	t.Run("Test failing to get a language", func(t *testing.T) {
		daoLanguage = DAO{}.New(&mocks.MockedLanguageModel{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.AddParam("id", "1")
		GetByIdLanguageEndpoint(c)
		strJsonError, err := json.Marshal(gin.H{"error": mocks.MockedErrorMessage})
		assert.NoError(t, err)
		assert.Equal(t, w.Body.Bytes(), strJsonError)
		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

	})
	t.Run("Test return a language as expect", func(t *testing.T) {
		var result model.LanguageModel
		daoLanguage = DAO{}.New(&model.LanguageModel{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.AddParam("id", "1")
		GetByIdLanguageEndpoint(c)
		assert.Equal(t, w.Result().StatusCode, http.StatusOK, "status code fail to meet the expectation")
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err, "return fail to meet the expectation")
		assert.Equal(t, 1, len(result.Languages))
		assert.Equal(t, uint(1), result.Languages[0].ID)
	})
}

func TestLanguageGetAllEndpoint(t *testing.T) {
	t.Run("Test failing to get the language", func(t *testing.T) {
		daoLanguage = DAO{}.New(&mocks.MockedLanguageModel{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		GetAllLanguageEndpoint(c)
		strJsonError, err := json.Marshal(gin.H{"error": mocks.MockedErrorMessage})
		assert.NoError(t, err)
		assert.Equal(t, w.Body.Bytes(), strJsonError)
		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

	})
	t.Run("Test return the language as expect", func(t *testing.T) {
		var result model.LanguageModel
		daoLanguage = DAO{}.New(&model.LanguageModel{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		GetAllLanguageEndpoint(c)
		assert.Equal(t, w.Result().StatusCode, http.StatusOK, "status code fail to meet the expectation")
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err, "return fail to meet the expectation")
	})
}
