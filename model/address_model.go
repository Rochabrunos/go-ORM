package model

import "time"

type Address struct {
	ID          uint   `gorm:"primaryKey;column:address_id`
	Address     string `gorm:"size:200"`
	Address2    string `gorm:"size:200"`
	District    string `gorm:"size:100"`
	CityID      uint
	City        *City     `json:"-" gorm:"foreignKey:CityID"`
	PostalCode  string    `gorm:"size:20"`
	Phone       string    `gorm:"size:40"`
	LasatUpdate time.Time `gorm:"autoUpdateTime"`
}

func (Address) TableName() string {
	return "address"
}
