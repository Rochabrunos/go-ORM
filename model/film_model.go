package model

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Film struct {
	ID              uint   `gorm:"primaryKey;column:film_id"`
	Title           string `gorm:"size:500"`
	ReleaseYear     string
	LanguageID      uint
	Language        *Language `json:"-" gorm:"foreignKey:LanguageID"`
	RentalDuration  string
	RentalRate      string
	Length          string
	ReplacementCost float64
	Rating          string
	LastUpdate      time.Time `gorm:"autoUpdateTime"`
	SpecialFeatures string
	FullText        string
}

func (Film) TableName() string {
	return "film"
}

func GetFilmById(c *gin.Context, db *gorm.DB) (*Film, error) {
	var film Film
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return nil, errors.New("invalid id, please, make sure to pass a number")
	}

	film.ID = uint(id)
	result := db.First(&film)
	return &film, result.Error
}

func GetAllFilms(c *gin.Context, db *gorm.DB) (*[]Film, error) {
	var films []Film
	page, _ := strconv.Atoi(c.DefaultQuery("p", "0"))
	result := db.Offset(page * 10).Limit(10).Find(&films)
	return &films, result.Error
}

func CreateNewFilm(c *gin.Context, db *gorm.DB) (*Film, error) {
	var film Film
	if err := c.ShouldBindJSON(&film); err != nil {
		return nil, err
	}
	result := db.Create(&film)
	return &film, result.Error
}

func UpdateFilmById(c *gin.Context, db *gorm.DB) (*Film, error) {
	film, err := GetFilmById(c, db)
	if err != nil {
		return nil, err
	}
	if err := c.ShouldBindJSON(film); err != nil {
		return nil, err
	}
	result := db.Save(film)
	return film, result.Error
}

func DeleteFilmById(c *gin.Context, db *gorm.DB) (*Film, error) {
	film, err := GetFilmById(c, db)
	if err != nil {
		return nil, err
	}
	result := db.Delete(film)
	return film, result.Error
}
