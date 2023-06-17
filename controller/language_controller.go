package controller

import (
	"net/http"
	model "orm-golang/model"
	service "orm-golang/service"

	"github.com/gin-gonic/gin"
)

func GetLangByIdEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	obj, err := model.GetLanguageById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func GetAllLangsEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	obj, err := model.GetAllLanguages(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func CreateLangEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	obj, err := model.CreateNewLanguage(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func ModifyLangEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	obj, err := model.UpdateLanguageById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func DeleteLangEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	obj, err := model.DeleteLanguageById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}
