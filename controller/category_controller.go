package controller

import (
	"net/http"
	model "orm-golang/model"
	service "orm-golang/service"

	"github.com/gin-gonic/gin"
)

var daoCategory = DAO{}.New(&model.CategoryModel{}, service.GetDBConnection())

func GetByIdCategoryEndpoint(c *gin.Context) {
	err := daoCategory.Model.GetById(c, daoCategory.Conn)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoCategory.Model)
}

func GetAllCategoryEndpoint(c *gin.Context) {
	err := daoCategory.Model.GetAll(c, daoCategory.Conn)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoCategory.Model)
}

func CreateCategoryEndpoint(c *gin.Context) {
	err := daoCategory.Model.CreateNew(c, daoCategory.Conn)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoCategory.Model)
}

func ModifyCategoryEndpoint(c *gin.Context) {
	err := daoCategory.Model.UpdateById(c, daoCategory.Conn)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoCategory.Model)
}

func DeleteCategoryEndpoint(c *gin.Context) {
	err := daoCategory.Model.DeleteById(c, daoCategory.Conn)
	if err != nil {
		SendError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, daoCategory.Model)
}
