package model

import "time"

type Actor struct {
	ID uint `gorm:"primaryKey;column:actor_id"`
	FistName string `gorm:"size:50"`
	LastName string `gorm:"size:50"`
	LastUpdate time.Time `gorm:"autoUpdateTime"`
}

func (Actor) TableName() string {
	return "actor"
}