package controller

import (
	"net/http"
	model "orm-golang/model"
	service "orm-golang/service"

	"github.com/gin-gonic/gin"
)

var daoFilm = DAO{}.New(&model.CategoryModel{})

func GetFilmByIdEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := daoFilm.Model.GetById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoFilm.Model)
}

func GetAllFilmsEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := daoFilm.Model.GetAll(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoFilm.Model)
}

func CreateFilmEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := daoFilm.Model.CreateNew(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoFilm.Model)
}

func ModifyFilmEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := daoFilm.Model.UpdateById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoFilm.Model)
}

func DeleteFilmEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := daoFilm.Model.DeleteById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoFilm.Model)
}
