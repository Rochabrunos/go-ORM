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
	Name       string    `gorm:"size:100"`
	LastUpdate time.Time `json:",omitempty" gorm:"autoUpdateTime"`
}

func (Category) TableName() string {
	return "category"
}

func (c *Category) GetCategoryById(ctx *gin.Context, db *gorm.DB) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return errors.New("invalid id, make sure to pass a number")
	}
	c.ID = uint(id)

	if result := db.First(c); result.Error != nil {
		return result.Error
	}
	return nil
}

// func (c *Category) GetAllCategories(ctx *gin.Context, db *gorm.DB) error {
// 	page, _ := strconv.Atoi(ctx.DefaultQuery("p", "0"))

// 	if result := db.Offset(page * 10).Limit(10).Find(&categories); result.Error != nil {
// 		return result.Error
// 	}
// 	return nil
// }

func (c *Category) CreateNewCategory(ctx *gin.Context, db *gorm.DB) error {
	if err := ctx.ShouldBindJSON(c); err != nil {
		return err
	}
	if result := db.Create(c); result.Error != nil {
		return result.Error
	}
	return nil
}
func (c *Category) UpdateCategoryById(ctx *gin.Context, db *gorm.DB) error {
	err := c.GetCategoryById(ctx, db)
	if err != nil {
		return err
	}
	if err := ctx.ShouldBindJSON(c); err != nil {
		return err
	}
	if result := db.Save(c); result.Error != nil {
		return result.Error
	}
	return nil
}

func (c *Category) DeleteCategoryById(ctx *gin.Context, db *gorm.DB) error {
	err := c.GetCategoryById(ctx, db)
	if err != nil {
		return err
	}
	if result := db.Delete(c); result.Error != nil {
		return result.Error
	}
	return nil
}
