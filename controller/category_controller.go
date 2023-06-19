package controller

import (
	"net/http"
	model "orm-golang/model"
	service "orm-golang/service"

	"github.com/gin-gonic/gin"
)

var handler = Handler{}.New(&model.Category{})

func GetCategByIdEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := handler.Model.GetCategoryById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, handler.Model)
}

// func GetAllCategsEndpoint(c *gin.Context) {
// 	connDB := service.GetDBConnection()
// 	err := handler.Model.GetAllCategories(c, connDB)
// 	if err != nil {
// 		SendError(c, http.StatusBadRequest, err)
// 		return
// 	}
// 	c.JSON(http.StatusOK, handler.Model)
// }

func CreateCategEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := handler.Model.CreateNewCategory(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, handler.Model)
}

func ModifyCategEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := handler.Model.UpdateCategoryById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, handler.Model)
}

func DeleteCategEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := handler.Model.DeleteCategoryById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, handler.Model)
}
