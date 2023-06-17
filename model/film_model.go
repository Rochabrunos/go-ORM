package model

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Film struct {
	ID              uint   `json:",omitempty" gorm:"primaryKey;column:film_id"`
	Title           string `gorm:"size:255"`
	Description     string
	ReleaseYear     string
	LanguageID      uint
	Language        *Language `json:"-" gorm:"foreignKey:LanguageID"`
	RentalDuration  uint      `gorm:"default:3"`
	RentalRate      float32   `gorm:"default:4.99"`
	Length          uint
	ReplacementCost float32   `gorm:"default:19.99"`
	Rating          string    `gorm:"default:G"`
	LastUpdate      time.Time `gorm:"autoUpdateTime"`
	SpecialFeatures []string  `gorm:"serializer:json"`
	FullText        string
}

func (Film) TableName() string {
	return "film"
}

func GetFilmById(c *gin.Context, db *gorm.DB) (*Film, error) {
	var film Film
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return nil, errors.New("invalid id, make sure to pass a number")
	}

	film.ID = uint(id)
	if result := db.First(&film); result.Error != nil {
		return nil, result.Error
	}
	return &film, nil
}

func GetAllFilms(c *gin.Context, db *gorm.DB) (*[]Film, error) {
	var films []Film

	if result := db.Find(&films); result.Error != nil {
		return nil, result.Error
	}
	return &films, nil
}

func CreateNewFilm(c *gin.Context, db *gorm.DB) (*Film, error) {
	var film Film
	if err := c.ShouldBindJSON(&film); err != nil {
		return nil, err
	}

	if result := db.Create(&film); result.Error != nil {
		return nil, result.Error
	}
	return &film, nil
}

func UpdateFilmById(c *gin.Context, db *gorm.DB) (*Film, error) {
	film, err := GetFilmById(c, db)
	if err != nil {
		return nil, err
	}
	if err := c.ShouldBindJSON(film); err != nil {
		return nil, err
	}

	if result := db.Save(film); result != nil {
		return nil, result.Error
	}
	return film, nil
}

func DeleteFilmById(c *gin.Context, db *gorm.DB) (*Film, error) {
	film, err := GetFilmById(c, db)
	if err != nil {
		return nil, err
	}

	if result := db.Delete(film); result.Error != nil {
		return nil, result.Error
	}
	return film, nil
}
