package model

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

func GetFilmById(c *gin.Context) (*Film, error) {
	var film Film
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return nil, errors.New("invalid id, make sure to pass a number")
	}

	film.ID = uint(id)
	if result := DB.First(&film); result.Error != nil {
		return nil, result.Error
	}
	return &film, nil
}

func GetAllFilms(c *gin.Context) (*[]Film, error) {
	var films []Film
	page, _ := strconv.Atoi(c.DefaultQuery("p", "0"))
	if result := DB.Offset(page * 10).Limit(10).Find(&films); result != nil {
		return nil, result.Error
	}
	return &films, nil
}

func CreateNewFilm(c *gin.Context) (*Film, error) {
	var film Film
	if err := c.ShouldBindJSON(&film); err != nil {
		return nil, err
	}
	if result := DB.Create(&film); result != nil {
		return nil, result.Error
	}
	return &film, nil
}

func UpdateFilmById(c *gin.Context) (*Film, error) {
	film, err := GetFilmById(c)
	if err != nil {
		return nil, err
	}
	if err := c.ShouldBindJSON(film); err != nil {
		return nil, err
	}
	if result := DB.Save(film); result != nil {
		return nil, result.Error
	}
	return film, nil
}

func DeleteFilmById(c *gin.Context) (*Film, error) {
	film, err := GetFilmById(c)
	if err != nil {
		return nil, err
	}
	if result := DB.Delete(film); result.Error != nil {
		return nil, result.Error
	}
	return film, nil
}
