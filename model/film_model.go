package model

import (
	"time"
	"strconv"
	"errors"
	"github.com/gin-gonic/gin"
)

type Film struct {
	ID uint `gorm:"primaryKey;column:film_id"`
	Title string `gorm:"size:500"`
	ReleaseYear string
	LanguageID uint
	Language *Language `json:"-" gorm:"foreignKey:LanguageID"`
	RentalDuration string
	RentalRate string
	Length string  
	ReplacementCost float64
	Rating string
	LastUpdate time.Time `gorm:"autoUpdateTime"`
	SpecialFeatures string
	FullText string
}

func (Film) TableName() string {
	return "film"
}

func GetFilmById(c *gin.Context) (*Film, error) {
	var film Film
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return nil, errors.New("invalid id, please, make sure to pass a number")
	}
	
	film.ID = uint(id)
	result := DB.First(&film)
	return &film, result.Error
}
	
func GetAllFilms(c *gin.Context) (*[]Film, error) {
	var films []Film
	page, _ := strconv.Atoi(c.DefaultQuery("p", "0"))
	result := DB.Offset(page*10).Limit(10).Find(&films)
	return &films, result.Error
}

func CreateNewFilm(c *gin.Context) (*Film, error) {
	var film Film
	if err := c.ShouldBindJSON(&film); err != nil {
		return nil, err
	}
	result := DB.Create(&film)
	return &film, result.Error
}

func UpdateFilmById(c *gin.Context) (*Film, error){
	film, err := GetFilmById(c)
	if err != nil {
		return nil, err
	}
	if err := c.ShouldBindJSON(film); err != nil {
		return nil, err
	}
	result := DB.Save(film)
	return film, result.Error	
}

func DeleteFilmById(c *gin.Context) (*Film, error) {
	film, err := GetFilmById(c)
	if err != nil {
		return nil, err
	}
	result := DB.Delete(film)
	return film, result.Error
}