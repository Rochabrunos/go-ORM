package model

import "time"

type Inventory struct {
	ID         uint `gorm:"primaryKey;column:inventory_id"`
	FilmID     uint
	Film       *Film `json:"-" gorm:"foreignKey:FilmID"`
	StoreID    uint
	Store      *Store    `json:"-" gorm:"foreignKey:StoreID"`
	LastUpdate time.Time `gorm:"autoUpdateTime"`
}

func (Inventory) TableName() string {
	return "inventory"
}
