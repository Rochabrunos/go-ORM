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
		return nil, errors.New("invalid id, make sure to pass a number")
	}
	category.ID = uint(id)

	if result := DB.First(&category); result.Error != nil {
		return nil, result.Error
	}
	return &category, nil
}

func GetAllCategories(c *gin.Context) (*[]Category, error) {
	var category []Category
	page, _ := strconv.Atoi(c.DefaultQuery("p", "0"))

	if result := DB.Offset(page * 10).Limit(10).Find(&category); result.Error != nil {
		return nil, result.Error
	}
	return &category, nil
}

func CreateNewCategory(c *gin.Context) (*Category, error) {
	var category Category
	if err := c.ShouldBindJSON(&category); err != nil {
		return nil, err
	}
	if result := DB.Create(&category); result.Error != nil {
		return nil, result.Error
	}
	return &category, nil
}

func UpdateCategoryById(c *gin.Context) (*Category, error) {
	category, err := GetCategoryById(c)
	if err != nil {
		return nil, err
	}
	if err := c.ShouldBindJSON(category); err != nil {
		return nil, err
	}
	if result := DB.Save(category); result.Error != nil {
		return nil, result.Error
	}
	return category, nil
}

func DeleteCategoryById(c *gin.Context) (*Category, error) {
	obj, err := GetCategoryById(c)
	if err != nil {
		return nil, err
	}
	if result := DB.Delete(obj); result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}
