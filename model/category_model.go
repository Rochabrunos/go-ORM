package model

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Category struct {
	ID         uint      `gorm:"primaryKey;column:category_id"`
	Name       string    `gorm:"size:100"`
	LastUpdate time.Time `gorm:"autoUpdateTime"`
}

func (Category) TableName() string {
	return "category"
}

func GetCategoryById(c *gin.Context) (*Category, error) {
	var category Category
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return nil, errors.New("invalid id, please, make sure to pass a number")
	}

	category.ID = uint(id)
	result := DB.First(&category)
	return &category, result.Error
}

func GetAllCategories(c *gin.Context) (*[]Category, error) {
	var category []Category
	page, _ := strconv.Atoi(c.DefaultQuery("p", "0"))
	result := DB.Offset(page * 10).Limit(10).Find(&category)
	return &category, result.Error
}

func CreateNewCategory(c *gin.Context) (*Category, error) {
	var category Category
	if err := c.ShouldBindJSON(&category); err != nil {
		return nil, err
	}
	result := DB.Create(&category)
	return &category, result.Error
}

func UpdateCategoryById(c *gin.Context) (*Category, error) {
	obj, err := GetCategoryById(c)
	if err != nil {
		return nil, err
	}
	if err := c.ShouldBindJSON(obj); err != nil {
		return nil, err
	}
	result := DB.Save(obj)
	return obj, result.Error
}

func DeleteCategoryById(c *gin.Context) (*Category, error) {
	obj, err := GetCategoryById(c)
	if err != nil {
		return nil, err
	}
	result := DB.Delete(obj)
	return obj, result.Error
}
