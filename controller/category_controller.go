package controller

import (
	"net/http"
	model "orm-golang/model"
	service "orm-golang/service"

	"github.com/gin-gonic/gin"
)

var daoCategory = DAO{}.New(&model.CategoryModel{})

func GetByIdEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := daoCategory.Model.GetById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoCategory.Model)
}

func GetAllEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := daoCategory.Model.GetAll(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoCategory.Model)
}

func CreateEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := daoCategory.Model.CreateNew(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoCategory.Model)
}

func ModifyEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := daoCategory.Model.UpdateById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoCategory.Model)
}

func DeleteEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := daoCategory.Model.DeleteById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoCategory.Model)
}
