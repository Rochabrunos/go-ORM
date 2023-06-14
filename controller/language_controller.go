package controller

import (
	"net/http"
	"github.com/gin-gonic/gin"
	model"orm-golang/model"
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
	obj, err := model.GetAllLanguanges(c)
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
