package model

import "time"

type City struct {
	ID         uint   `gorm:"primaryKey;column:city_id"`
	City       string `gorm:"size:100"`
	CountryID  uint
	Country    *Country  `json:"-" gorm:"foreignKey:CountryID"`
	LastUpdate time.Time `gorm:"autoUpdateTime"`
}

func (City) TableName() string {
	return "city"
}
