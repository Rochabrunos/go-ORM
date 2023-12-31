package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	model "orm-golang/model"
	mocks "orm-golang/model/mock"
	"orm-golang/service"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	db := GetDBTestConnection()
	fmt.Println("Droping table")
	db.Migrator().DropTable(&model.Language{})
	fmt.Println("Migrating table")
	db.AutoMigrate(&model.Language{})
}

func TestLanguageGetByIdEndpoint(t *testing.T) {
	t.Run("Test failing to get a language", func(t *testing.T) {
		daoLanguage = DAO{}.New(&mocks.MockedLanguageModel{}, GetDBTestConnection())

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

		daoLanguage = DAO{}.New(&model.LanguageModel{}, service.GetDBConnection())

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
		daoLanguage = DAO{}.New(&mocks.MockedLanguageModel{}, GetDBTestConnection())

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

		daoLanguage = DAO{}.New(&model.LanguageModel{}, service.GetDBConnection())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		GetAllLanguageEndpoint(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusOK, "status code fail to meet the expectation")
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err, "return fail to meet the expectation")
	})
}

func TestCreateLanguageEndpoint(t *testing.T) {
	t.Run("Test fail to create the language", func(t *testing.T) {
		daoLanguage = DAO{}.New(&mocks.MockedLanguageModel{}, GetDBTestConnection())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		CreateLanguageEndpoint(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest, "status code fail to meet the expectation")
		strJsonError, err := json.Marshal(gin.H{"error": mocks.MockedErrorMessage})
		assert.NoError(t, err)
		assert.Equal(t, w.Body.Bytes(), strJsonError)
	})

	t.Run("Test successfully creates the language", func(t *testing.T) {
		var result model.LanguageModel

		db := GetDBTestConnection()
		db.Exec("TRUNCATE TABLE language RESTART IDENTITY CASCADE")
		daoLanguage = DAO{}.New(&model.LanguageModel{}, db)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		data := map[string]string{"Name": "English"}
		jsonBytes, err := json.Marshal(data)
		assert.NoError(t, err)
		c.Request = &http.Request{}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))

		CreateLanguageEndpoint(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusOK, "status code fail to meet the expectation")
		err = json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, data["Name"], result.Languages[0].Name)
	})
}

func TestModifyLanguageEndpoint(t *testing.T) {
	t.Run("Test fail to modify the language", func(t *testing.T) {
		daoLanguage = DAO{}.New(&mocks.MockedLanguageModel{}, GetDBTestConnection())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		ModifyLanguageEndpoint(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest, "status code fail to meet the expectation")
		strJsonError, err := json.Marshal(gin.H{"error": mocks.MockedErrorMessage})
		assert.NoError(t, err)
		assert.Equal(t, w.Body.Bytes(), strJsonError)
	})

	t.Run("Test successfully modifies the language", func(t *testing.T) {
		var result model.LanguageModel

		db := GetDBTestConnection()
		db.Exec("TRUNCATE TABLE language RESTART IDENTITY CASCADE")
		db.Create(&model.Language{Name: "Portuguese"})

		daoLanguage = DAO{}.New(&model.LanguageModel{}, db)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.AddParam("id", "1")
		data := map[string]string{"Name": "English"}
		jsonBytes, err := json.Marshal(data)
		assert.NoError(t, err)
		c.Request = &http.Request{}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))

		CreateLanguageEndpoint(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusOK, "status code fail to meet the expectation")
		err = json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, data["Name"], result.Languages[0].Name)
	})
}

func TestDeleteLanguageEndpoint(t *testing.T) {
	t.Run("Test fail to delete the language", func(t *testing.T) {
		daoLanguage = DAO{}.New(&mocks.MockedLanguageModel{}, GetDBTestConnection())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		DeleteLanguageEndpoint(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest, "status code fail to meet the expectation")
		strJsonError, err := json.Marshal(gin.H{"error": mocks.MockedErrorMessage})
		assert.NoError(t, err)
		assert.Equal(t, w.Body.Bytes(), strJsonError)
	})

	t.Run("Test successfully deletes the language", func(t *testing.T) {
		var result model.LanguageModel

		db := GetDBTestConnection()
		db.Exec("TRUNCATE TABLE language RESTART IDENTITY CASCADE")
		languageData := model.Language{Name: "Portuguse"}
		db.Create(&languageData)

		daoLanguage = DAO{}.New(&model.LanguageModel{}, db)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.AddParam("id", "1")

		DeleteLanguageEndpoint(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusOK, "status code fail to meet the expectation")
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, languageData.Name, result.Languages[0].Name)
	})
}
