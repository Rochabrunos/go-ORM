package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	Model Model
}

type Model interface {
	GetCategoryById(*gin.Context, *gorm.DB) error
	// GetAllCategories(*gin.Context, *gorm.DB) error
	CreateNewCategory(*gin.Context, *gorm.DB) error
	UpdateCategoryById(*gin.Context, *gorm.DB) error
	DeleteCategoryById(*gin.Context, *gorm.DB) error
}

func (Handler) New(c Model) Handler {
	h := Handler{
		Model: c,
	}

	return h
}
