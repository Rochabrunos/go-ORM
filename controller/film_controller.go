package controller

import (
	"net/http"
	model "orm-golang/model"
	service "orm-golang/service"

	"github.com/gin-gonic/gin"
)

var daoFilm = DAO{}.New(&model.FilmModel{}, service.GetDBConnection())

func GetByIdFilmEndpoint(c *gin.Context) {
	err := daoFilm.Model.GetById(c, daoFilm.Conn)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoFilm.Model)
}

func GetAllFilmsEndpoint(c *gin.Context) {
	err := daoFilm.Model.GetAll(c, daoFilm.Conn)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoFilm.Model)
}

func CreateFilmEndpoint(c *gin.Context) {
	err := daoFilm.Model.CreateNew(c, daoFilm.Conn)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoFilm.Model)
}

func ModifyFilmEndpoint(c *gin.Context) {
	err := daoFilm.Model.UpdateById(c, daoFilm.Conn)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoFilm.Model)
}

func DeleteFilmEndpoint(c *gin.Context) {
	err := daoFilm.Model.DeleteById(c, daoFilm.Conn)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoFilm.Model)
}
