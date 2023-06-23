package controller

import (
	"net/http"
	model "orm-golang/model"
	service "orm-golang/service"

	"github.com/gin-gonic/gin"
)

var daoLanguage = DAO{}.New(&model.LanguageModel{})

func GetByIdLanguageEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := daoLanguage.Model.GetById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoLanguage.Model)
}

func GetAllLanguageEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := daoLanguage.Model.GetAll(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoLanguage.Model)
}

func CreateLanguageEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := daoLanguage.Model.CreateNew(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoLanguage.Model)
}

func ModifyLanguageEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := daoLanguage.Model.UpdateById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoLanguage.Model)
}

func DeleteLanguageEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := daoLanguage.Model.DeleteById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoLanguage.Model)
}
