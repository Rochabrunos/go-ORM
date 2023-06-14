package model

import "time"

type FilmActor struct {
	ID         uint      `gorm:"primaryKey"`
	Actor      *Actor    `json:"-" gorm:"foreignKey:ActorID"`
	FilmID     uint      `gorm:"primaryKey"`
	Film       *Film     `json:"-" gorm:"foreignKey:FilmID"`
	LastUpdate time.Time `gorm:"autoUpdateTime"`
}

func (FilmActor) TableName() string {
	return "film_actor"
}
