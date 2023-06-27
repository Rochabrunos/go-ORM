package controller

import (
	"net/http"
	model "orm-golang/model"
	service "orm-golang/service"

	"github.com/gin-gonic/gin"
)

var daoLanguage = DAO{}.New(&model.LanguageModel{}, service.GetDBConnection())

func GetByIdLanguageEndpoint(c *gin.Context) {
	err := daoLanguage.Model.GetById(c, daoLanguage.Conn)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoLanguage.Model)
}

func GetAllLanguageEndpoint(c *gin.Context) {
	err := daoLanguage.Model.GetAll(c, daoLanguage.Conn)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoLanguage.Model)
}

func CreateLanguageEndpoint(c *gin.Context) {
	err := daoLanguage.Model.CreateNew(c, daoLanguage.Conn)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoLanguage.Model)
}

func ModifyLanguageEndpoint(c *gin.Context) {
	err := daoLanguage.Model.UpdateById(c, daoLanguage.Conn)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoLanguage.Model)
}

func DeleteLanguageEndpoint(c *gin.Context) {
	err := daoLanguage.Model.DeleteById(c, daoLanguage.Conn)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoLanguage.Model)
}
