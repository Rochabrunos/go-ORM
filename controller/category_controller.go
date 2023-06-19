package controller

import (
	"net/http"
	model "orm-golang/model"
	service "orm-golang/service"

	"github.com/gin-gonic/gin"
)

var dao = DAO{}.New(&model.CategoryModel{})

func GetByIdEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := dao.Model.GetById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dao.Model)
}

func GetAllEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := dao.Model.GetAll(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dao.Model)
}

func CreateEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := dao.Model.CreateNew(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dao.Model)
}

func ModifyEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := dao.Model.UpdateById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dao.Model)
}

func DeleteEndpoint(c *gin.Context) {
	connDB := service.GetDBConnection()
	err := dao.Model.DeleteById(c, connDB)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, dao.Model)
}
