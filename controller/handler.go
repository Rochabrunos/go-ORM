package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DAO struct {
	Model Model
}

type Model interface {
	GetById(*gin.Context, *gorm.DB) error
	GetAll(*gin.Context, *gorm.DB) error
	CreateNew(*gin.Context, *gorm.DB) error
	UpdateById(*gin.Context, *gorm.DB) error
	DeleteById(*gin.Context, *gorm.DB) error
}

func (DAO) New(c Model) DAO {
	h := DAO{
		Model: c,
	}

	return h
}
