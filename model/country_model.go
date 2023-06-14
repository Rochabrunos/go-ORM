package model

import "time"

type Country struct {
	ID	uint	`gorm:"primaryKey;column:country_id"`
	Country	string	`gorm:"size:100"`
	LastUpdate time.Time	`gorm:"autoUpdateTime"`
}

func (Country) TableName() string {
	return "country"
}