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

var mockedFilms = []model.Film{
	{
		Title:           "Chamber Italian",
		Description:     "A Fateful Reflection of a Moose And a Husband who must Overcome a Monkey in Nigeria",
		ReleaseYear:     "2006",
		Language:        model.Language{ID: 1, Name: "English"},
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
		Language:        model.Language{ID: 1, Name: "English"},
		RentalDuration:  5,
		RentalRate:      4.99,
		Length:          49,
		ReplacementCost: 19.99,
		Rating:          "R",
		SpecialFeatures: []string{"Behind the Scenes"},
		FullText:        `'australia':18 'cat':8 'drama':5 'epic':4 'explor':11 'gross':1 'moos':16 'must':13 'redeem':14 'wonder':2`,
	},
}

func init() {
	db := GetDBTestConnection()
	db.Migrator().DropTable(&model.Language{})
	db.Migrator().DropTable(&model.Film{})
	db.AutoMigrate(&model.Language{})
	db.Create(&model.Language{Name: "English"})
	db.AutoMigrate(&model.Film{})
}
func TestGetByIdFilmEndpoint(t *testing.T) {
	t.Run("Test failing to get a film", func(t *testing.T) {
		daoFilm = DAO{}.New(&mocks.MockedFilmModel{}, GetDBTestConnection())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.AddParam("id", "1")

		GetByIdFilmEndpoint(c)

		strJsonError, err := json.Marshal(gin.H{"error": mocks.MockedErrorMessage})
		assert.NoError(t, err)
		assert.Equal(t, w.Body.Bytes(), strJsonError)
		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

	})
	t.Run("Test return a film as expect", func(t *testing.T) {
		var result model.FilmModel

		daoFilm = DAO{}.New(&model.FilmModel{}, service.GetDBConnection())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.AddParam("id", "1")

		GetByIdFilmEndpoint(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusOK, "status code fail to meet the expectation")
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err, "return fail to meet the expectation")
		assert.Equal(t, 1, len(result.Films))
		assert.Equal(t, uint(1), result.Films[0].ID)
	})
}

func TestGetAllFilmEndpoint(t *testing.T) {
	t.Run("Test failing to get the categories", func(t *testing.T) {
		daoFilm = DAO{}.New(&mocks.MockedFilmModel{}, GetDBTestConnection())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		GetAllFilmsEndpoint(c)

		strJsonError, err := json.Marshal(gin.H{"error": mocks.MockedErrorMessage})
		assert.NoError(t, err)
		assert.Equal(t, w.Body.Bytes(), strJsonError)
		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

	})
	t.Run("Test return the categories as expect", func(t *testing.T) {
		var result model.FilmModel

		daoFilm = DAO{}.New(&model.FilmModel{}, service.GetDBConnection())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		GetAllFilmsEndpoint(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusOK, "status code fail to meet the expectation")
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err, "return fail to meet the expectation")
	})
}

func TestCreateFilmEndpoint(t *testing.T) {
	t.Run("Test fail to create the film", func(t *testing.T) {
		daoFilm = DAO{}.New(&mocks.MockedFilmModel{}, GetDBTestConnection())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		CreateFilmEndpoint(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest, "status code fail to meet the expectation")
		strJsonError, err := json.Marshal(gin.H{"error": mocks.MockedErrorMessage})
		assert.NoError(t, err)
		assert.Equal(t, w.Body.Bytes(), strJsonError)
	})

	t.Run("Test successfully creates the film", func(t *testing.T) {
		var result model.FilmModel

		db := GetDBTestConnection()
		db.Exec("TRUNCATE TABLE film RESTART IDENTITY CASCADE")

		daoFilm = DAO{}.New(&model.FilmModel{}, db)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		data := mockedFilms[0]
		jsonBytes, err := json.Marshal(data)
		assert.NoError(t, err)
		c.Request = &http.Request{}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))

		CreateFilmEndpoint(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusOK, "status code fail to meet the expectation")
		err = json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, data.Title, result.Films[0].Title)
	})
}

func TestModifyFilmEndpoint(t *testing.T) {
	t.Run("Test fail to modify the film", func(t *testing.T) {
		daoFilm = DAO{}.New(&mocks.MockedFilmModel{}, GetDBTestConnection())
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		ModifyFilmEndpoint(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest, "status code fail to meet the expectation")
		strJsonError, err := json.Marshal(gin.H{"error": mocks.MockedErrorMessage})
		assert.NoError(t, err)
		assert.Equal(t, w.Body.Bytes(), strJsonError)
	})

	t.Run("Test successfully modifies the film", func(t *testing.T) {
		var result model.FilmModel

		db := GetDBTestConnection()
		db.Exec("TRUNCATE TABLE film RESTART IDENTITY CASCADE")
		filmData := mockedFilms[0]
		db.Create(&filmData)

		daoFilm = DAO{}.New(&model.FilmModel{}, db)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.AddParam("id", "1")
		data := mockedFilms[1]
		jsonBytes, err := json.Marshal(data)
		assert.NoError(t, err)
		c.Request = &http.Request{}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))

		ModifyFilmEndpoint(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusOK, "status code fail to meet the expectation")
		err = json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, data.Title, result.Films[0].Title)
	})
}

func TestDeleteFilmEndpoint(t *testing.T) {
	t.Run("Test fail to delete the film", func(t *testing.T) {
		daoFilm = DAO{}.New(&mocks.MockedFilmModel{}, GetDBTestConnection())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		DeleteFilmEndpoint(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest, "status code fail to meet the expectation")
		strJsonError, err := json.Marshal(gin.H{"error": mocks.MockedErrorMessage})
		assert.NoError(t, err)
		assert.Equal(t, w.Body.Bytes(), strJsonError)
	})

	t.Run("Test successfully deletes the film", func(t *testing.T) {
		var result model.FilmModel

		db := GetDBTestConnection()
		db.Exec("TRUNCATE TABLE film RESTART IDENTITY CASCADE")
		filmData := mockedFilms[0]
		db.Create(&filmData)

		daoFilm = DAO{}.New(&model.FilmModel{}, db)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.AddParam("id", "1")

		DeleteFilmEndpoint(c)

		assert.Equal(t, w.Result().StatusCode, http.StatusOK, "status code fail to meet the expectation")
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, filmData.Title, result.Films[0].Title)
	})
}
