package model

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Film struct {
<<<<<<< HEAD
<<<<<<< HEAD
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
=======
=======
>>>>>>> caa7177 (Added unit tests for Film{} model)
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
<<<<<<< HEAD
>>>>>>> caa717724aed354e674ae8e21ea91d60755186a0
=======
>>>>>>> caa7177 (Added unit tests for Film{} model)
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
<<<<<<< HEAD
<<<<<<< HEAD
	result := db.First(&film)
	return &film, result.Error
}

func GetAllFilms(c *gin.Context, db *gorm.DB) (*[]Film, error) {
	var films []Film
	page, _ := strconv.Atoi(c.DefaultQuery("p", "0"))
	result := db.Offset(page * 10).Limit(10).Find(&films)
	return &films, result.Error
=======
	if result := DB.First(&film); result.Error != nil {
		return nil, result.Error
	}
	return &film, nil
=======
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
>>>>>>> caa7177 (Added unit tests for Film{} model)
}

func GetAllFilms(c *gin.Context) (*[]Film, error) {
	var films []Film
	page, _ := strconv.Atoi(c.DefaultQuery("p", "0"))
	if result := DB.Offset(page * 10).Limit(10).Find(&films); result != nil {
		return nil, result.Error
	}
	return &films, nil
>>>>>>> caa717724aed354e674ae8e21ea91d60755186a0
}

func CreateNewFilm(c *gin.Context, db *gorm.DB) (*Film, error) {
	var film Film
	if err := c.ShouldBindJSON(&film); err != nil {
		return nil, err
	}
<<<<<<< HEAD
<<<<<<< HEAD
	result := db.Create(&film)
	return &film, result.Error
}

func UpdateFilmById(c *gin.Context, db *gorm.DB) (*Film, error) {
	film, err := GetFilmById(c, db)
=======
	if result := DB.Create(&film); result != nil {
		return nil, result.Error
	}
	return &film, nil
}

=======
	if result := DB.Create(&film); result != nil {
		return nil, result.Error
	}
	return &film, nil
}

>>>>>>> caa7177 (Added unit tests for Film{} model)
func UpdateFilmById(c *gin.Context) (*Film, error) {
	film, err := GetFilmById(c)
>>>>>>> caa717724aed354e674ae8e21ea91d60755186a0
	if err != nil {
		return nil, err
	}
	if err := c.ShouldBindJSON(film); err != nil {
		return nil, err
	}
<<<<<<< HEAD
<<<<<<< HEAD
	result := db.Save(film)
	return film, result.Error
=======
=======
>>>>>>> caa7177 (Added unit tests for Film{} model)
	if result := DB.Save(film); result != nil {
		return nil, result.Error
	}
	return film, nil
<<<<<<< HEAD
>>>>>>> caa717724aed354e674ae8e21ea91d60755186a0
=======
>>>>>>> caa7177 (Added unit tests for Film{} model)
}

func DeleteFilmById(c *gin.Context, db *gorm.DB) (*Film, error) {
	film, err := GetFilmById(c, db)
	if err != nil {
		return nil, err
	}
<<<<<<< HEAD
<<<<<<< HEAD
	result := db.Delete(film)
	return film, result.Error
=======
=======
>>>>>>> caa7177 (Added unit tests for Film{} model)
	if result := DB.Delete(film); result.Error != nil {
		return nil, result.Error
	}
	return film, nil
<<<<<<< HEAD
>>>>>>> caa717724aed354e674ae8e21ea91d60755186a0
=======
>>>>>>> caa7177 (Added unit tests for Film{} model)
}
