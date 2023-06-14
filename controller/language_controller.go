package controller

import (
	"net/http"
	model "orm-golang/model"

	"github.com/gin-gonic/gin"
)

func GetLangByIdEndpoint(c *gin.Context) {
	obj, err := model.GetLanguageById(c)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func GetAllLangsEndpoint(c *gin.Context) {
	obj, err := model.GetAllLanguages(c)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func CreateLangEndpoint(c *gin.Context) {
	obj, err := model.CreateNewLanguage(c)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func ModifyLangEndpoint(c *gin.Context) {
	obj, err := model.UpdateLanguageById(c)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func DeleteLangEndpoint(c *gin.Context) {
	obj, err := model.DeleteLanguageById(c)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}
