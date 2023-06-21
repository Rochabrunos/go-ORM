package model

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Language struct {
	ID         uint      `json:",omitempty" gorm:"primaryKey;column:language_id"`
	Name       string    `gorm:"size:50" binding:"required"`
	LastUpdate time.Time `json:",omitempty" gorm:"autoUpdateTime"`
}

func (Language) TableName() string {
	return "language"
}

type LanguageModel struct {
	Languages []Language
}

func (l *LanguageModel) GetById(c *gin.Context, db *gorm.DB) error {
	var language Language
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.New("invalid id, make sure to pass a number")
	}
	if result := db.First(&language, uint(id)); result.Error != nil {
		return result.Error
	}

	l.Languages = []Language{language}
	return nil
}

func (l *LanguageModel) GetAll(c *gin.Context, db *gorm.DB) error {
	var languages []Language
	page, _ := strconv.Atoi(c.DefaultQuery("p", "0"))
	if result := db.Offset(page * 10).Limit(10).Find(&languages); result.Error != nil {
		return result.Error
	}
	l.Languages = languages
	return nil
}

func (l *LanguageModel) CreateNew(c *gin.Context, db *gorm.DB) error {
	var language Language
	if err := c.ShouldBindJSON(&language); err != nil {
		return err
	}
	if result := db.Create(&language); result.Error != nil {
		return result.Error
	}
	l.Languages = []Language{language}
	return nil
}

func (l *LanguageModel) UpdateById(c *gin.Context, db *gorm.DB) error {
	err := l.GetById(c, db)
	if err != nil {
		return err
	}
	language := Language{ID: l.Languages[0].ID}
	if err := c.ShouldBindJSON(&language); err != nil {
		return err
	}
	if result := db.Save(&language); result.Error != nil {
		return result.Error
	}
	l.Languages = []Language{language}
	return nil
}

func (l *LanguageModel) DeleteById(c *gin.Context, db *gorm.DB) error {
	err := l.GetById(c, db)
	if err != nil {
		return err
	}
	language := l.Languages[0]
	if result := db.Delete(language); result.Error != nil {
		return result.Error
	}
	return nil
}
