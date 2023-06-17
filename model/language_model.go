package model

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Language struct {
	ID         uint       `json:",omitempty" gorm:"primaryKey;column:language_id"`
	Name       string     `gorm:"size:50"`
	LastUpdate *time.Time `json:",omitempty" gorm:"autoUpdateTime:mili"`
}

func (Language) TableName() string {
	return "language"
}

func GetLanguageById(c *gin.Context, db *gorm.DB) (*Language, error) {
	var lang Language
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return nil, errors.New("invalid id, make sure to pass a number")
	}

	lang.ID = uint(id)
	if result := db.First(&lang); result.Error != nil {
		return nil, result.Error
	}
	return &lang, nil
}

func GetAllLanguages(c *gin.Context, db *gorm.DB) (*[]Language, error) {
	var langs []Language
	page, _ := strconv.Atoi(c.DefaultQuery("p", "0"))
	if result := db.Offset(page * 10).Limit(10).Find(&langs); result.Error != nil {
		return nil, result.Error
	}
	return &langs, nil
}

func CreateNewLanguage(c *gin.Context, db *gorm.DB) (*Language, error) {
	var newLang Language
	if err := c.ShouldBindJSON(&newLang); err != nil {
		return nil, err
	}
	if result := db.Create(&newLang); result.Error != nil {
		return nil, result.Error
	}
	return &newLang, nil
}

func UpdateLanguageById(c *gin.Context, db *gorm.DB) (*Language, error) {
	lang, err := GetLanguageById(c, db)
	if err != nil {
		return nil, err
	}
	if err := c.ShouldBindJSON(lang); err != nil {
		return nil, err
	}
	if result := db.Save(lang); result.Error != nil {
		return nil, result.Error
	}
	return lang, nil
}

func DeleteLanguageById(c *gin.Context, db *gorm.DB) (*Language, error) {
	lang, err := GetLanguageById(c, db)
	if err != nil {
		return nil, err
	}
	if result := db.Delete(lang); result.Error != nil {
		return nil, result.Error
	}
	return lang, nil
}
