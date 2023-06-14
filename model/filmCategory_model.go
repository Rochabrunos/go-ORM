package model

import (
	"time"
)

type FilmCategory struct {
	FilmID     uint      `gorm:"primaryKey"`
	Film       *Film     `json:"-" gorm:"foreignKey:"FilmID"`
	CategoryID uint      `gorm:"primaryKey"`
	Category   *Category `json:"-" gorm:"foreignKey:CategoryID"`
	LastUpdate time.Time `gorm:"autoUpdateTime"`
}

func (FilmCategory) TableName() string {
	return "film_category"
}
