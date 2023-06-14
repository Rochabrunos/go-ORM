package controller

import (
	"net/http"
	"github.com/gin-gonic/gin"
	model"orm-golang/model"
)

func GetFilmByIdEndpoint(c *gin.Context) {
	obj, err := model.GetFilmById(c)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func GetAllFilmsEndpoint(c *gin.Context) {
	obj, err := model.GetAllFilms(c)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func CreateFilmEndpoint(c *gin.Context) {
	obj, err := model.CreateNewFilm(c)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func ModifyFilmEndpoint(c *gin.Context) {
	obj, err := model.UpdateFilmById(c)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func DeleteFilmEndpoint(c *gin.Context) {
	obj, err := model.DeleteFilmById(c)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}