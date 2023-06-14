package controller

import (
	"net/http"
	"github.com/gin-gonic/gin"
	model"orm-golang/model"
)

func GetCategByIdEndpoint(c *gin.Context) {
	obj, err := model.GetCategoryById(c)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func GetAllCategsEndpoint(c *gin.Context) {
	obj, err := model.GetAllCategories(c)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func CreateCategEndpoint(c *gin.Context) {
	obj, err := model.CreateNewCategory(c)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func ModifyCategEndpoint(c *gin.Context) {
	obj, err := model.UpdateCategoryById(c)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}

func DeleteCategEndpoint(c *gin.Context) {
	obj, err := model.DeleteCategoryById(c)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, obj)
}