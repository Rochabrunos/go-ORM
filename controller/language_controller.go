package controller

import (
	"net/http"
	model "orm-golang/model"
	service "orm-golang/service"

	"github.com/gin-gonic/gin"
)

var daoLanguage = DAO{}.New(&model.LanguageModel{})

func GetLangByIdEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := daoLanguage.Model.GetById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoLanguage.Model)
}

func GetAllLangsEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := daoLanguage.Model.GetAll(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoLanguage.Model)
}

func CreateLangEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := daoLanguage.Model.CreateNew(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoLanguage.Model)
}

func ModifyLangEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := daoLanguage.Model.UpdateById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoLanguage.Model)
}

func DeleteLangEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := daoLanguage.Model.DeleteById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoLanguage.Model)
}
