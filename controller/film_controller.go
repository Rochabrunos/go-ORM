package controller

import (
	"net/http"
	model "orm-golang/model"
	service "orm-golang/service"

	"github.com/gin-gonic/gin"
)

func GetFilmByIdEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	obj, err := model.GetFilmById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func GetAllFilmsEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	obj, err := model.GetAllFilms(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func CreateFilmEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	obj, err := model.CreateNewFilm(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func ModifyFilmEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	obj, err := model.UpdateFilmById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func DeleteFilmEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	obj, err := model.DeleteFilmById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}
