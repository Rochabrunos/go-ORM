package model

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type SpecialFeatures []string

type Film struct {
	ID              uint   `json:",omitempty" gorm:"primaryKey;column:film_id"`
	Title           string `gorm:"size:255" binding:"required"`
	Description     string `binding:"required"`
	ReleaseYear     string `binding:"required"`
	LanguageID      uint
	Language        Language       `gorm:"foreignKey:LanguageID" binding:"required"`
	RentalDuration  uint           `gorm:"default:3" binding:"required"`
	RentalRate      float32        `gorm:"default:4.99" binding:"required"`
	Length          uint           `binding:"required"`
	ReplacementCost float32        `gorm:"default:19.99" binding:"required"`
	Rating          string         `gorm:"default:G" binding:"required"`
	LastUpdate      time.Time      `gorm:"autoUpdateTime"`
	SpecialFeatures pq.StringArray `gorm:"type:text[];serialize:json"`
	FullText        string
}

type FilmModel struct {
	Films []Film
}

func (Film) TableName() string {
	return "film"
}

func (f *FilmModel) GetById(c *gin.Context, db *gorm.DB) error {
	var film Film
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.New("invalid id, make sure to pass a number")
	}
	if result := db.First(&film, uint(id)); result.Error != nil {
		return result.Error
	}
	f.Films = []Film{film}
	return nil
}

func (f *FilmModel) GetAll(c *gin.Context, db *gorm.DB) error {
	var films []Film
	page, _ := strconv.Atoi(c.DefaultQuery("p", "0"))
	if result := db.Offset(page * 10).Limit(10).Find(&films); result.Error != nil {
		return result.Error
	}
	f.Films = films
	return nil
}

func (f *FilmModel) CreateNew(c *gin.Context, db *gorm.DB) error {
	var film Film
	fmt.Println(c.Request.Body)
	if err := c.ShouldBindJSON(&film); err != nil {
		return err
	}
	if result := db.Create(&film); result.Error != nil {
		return result.Error
	}
	f.Films = []Film{film}
	return nil
}

func (f *FilmModel) UpdateById(c *gin.Context, db *gorm.DB) error {
	var film Film
	err := f.GetById(c, db)
	if err != nil {
		return err
	}
	if err := c.ShouldBindJSON(&film); err != nil {
		return err
	}
	film.ID = f.Films[0].ID
	if result := db.Save(&film); result.Error != nil {
		return result.Error
	}
	f.Films[0] = film
	return nil
}

func (f *FilmModel) DeleteById(c *gin.Context, db *gorm.DB) error {
	var film Film
	err := f.GetById(c, db)
	if err != nil {
		return err
	}
	film.ID = f.Films[0].ID
	if result := db.Delete(&film); result.Error != nil {
		return result.Error
	}
	return nil
}
