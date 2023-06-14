package model

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Language struct {
	ID         uint       `gorm:"primaryKey;column:language_id"`
	Name       string     `gorm:"size:50"`
	LastUpdate *time.Time `gorm:"autoUpdateTime:mili"`
}

func (Language) TableName() string {
	return "language"
}

func GetLanguageById(c *gin.Context) (*Language, error) {
	var lang Language
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return nil, errors.New("invalid id, make sure to pass a number")
	}

	lang.ID = uint(id)
	result := DB.First(&lang)
	if result.Error != nil {
		return nil, result.Error
	}
	return &lang, nil
}

func GetAllLanguages(c *gin.Context) (*[]Language, error) {
	var langs []Language
	page, _ := strconv.Atoi(c.DefaultQuery("p", "0"))
	result := DB.Offset(page * 10).Limit(10).Find(&langs)
	return &langs, result.Error
}

func CreateNewLanguage(c *gin.Context) (*Language, error) {
	var newLang Language
	if err := c.ShouldBindJSON(&newLang); err != nil {
		return nil, err
	}
	result := DB.Create(&newLang)
	return &newLang, result.Error
}

func UpdateLanguageById(c *gin.Context) (*Language, error) {
	obj, err := GetLanguageById(c)
	if err != nil {
		return nil, err
	}
	if err := c.ShouldBindJSON(obj); err != nil {
		return nil, err
	}
	result := DB.Save(obj)
	return obj, result.Error
}

func DeleteLanguageById(c *gin.Context) (*Language, error) {
	obj, err := GetLanguageById(c)
	if err != nil {
		return nil, err
	}
	result := DB.Delete(obj)
	return obj, result.Error
}
