package model

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Category struct {
	ID         uint      `json:",omitempty" gorm:"primaryKey;column:category_id"`
	Name       string    `gorm:"size:100" binding:"required"`
	LastUpdate time.Time `json:",omitempty" gorm:"autoUpdateTime"`
}

type CategoryModel struct {
	Categories []Category
}

func (Category) TableName() string {
	return "category"
}

func (c *CategoryModel) GetById(ctx *gin.Context, db *gorm.DB) error {
	var category Category
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return errors.New("invalid id, make sure to pass a number")
	}
	category.ID = uint(id)

	if result := db.First(&category); result.Error != nil {
		return result.Error
	}
	c.Categories = []Category{category}
	return nil
}

func (c *CategoryModel) GetAll(ctx *gin.Context, db *gorm.DB) error {
	var categories []Category
	page, _ := strconv.Atoi(ctx.DefaultQuery("p", "0"))

	if result := db.Offset(page * 10).Limit(10).Find(&categories); result.Error != nil {
		return result.Error
	}
	c.Categories = categories
	return nil
}

func (c *CategoryModel) CreateNew(ctx *gin.Context, db *gorm.DB) error {
	var category Category

	if err := ctx.ShouldBindJSON(&category); err != nil {
		return err
	}
	if result := db.Create(&category); result.Error != nil {
		return result.Error
	}
	c.Categories = []Category{category}
	return nil
}
func (c *CategoryModel) UpdateById(ctx *gin.Context, db *gorm.DB) error {
	var category Category
	err := c.GetById(ctx, db)
	if err != nil {
		return err
	}
	if err := ctx.ShouldBindJSON(category); err != nil {
		return err
	}
	if result := db.Save(&category); result.Error != nil {
		return result.Error
	}
	c.Categories = []Category{category}
	return nil
}

func (c *CategoryModel) DeleteById(ctx *gin.Context, db *gorm.DB) error {
	var category Category
	err := c.GetById(ctx, db)
	category = c.Categories[0]
	if err != nil {
		return err
	}
	if result := db.Delete(&category); result.Error != nil {
		return result.Error
	}
	c.Categories = []Category{category}
	return nil
}
