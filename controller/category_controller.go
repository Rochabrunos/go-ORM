package controller

import (
	"net/http"
	model "orm-golang/model"
	service "orm-golang/service"

	"github.com/gin-gonic/gin"
)

func GetCategByIdEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	obj, err := model.GetCategoryById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func GetAllCategsEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	obj, err := model.GetAllCategories(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func CreateCategEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	obj, err := model.CreateNewCategory(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func ModifyCategEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	obj, err := model.UpdateCategoryById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func DeleteCategEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	obj, err := model.DeleteCategoryById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}
