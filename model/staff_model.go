package model

import "time"

type Staff struct {
	ID	uint	`gorm:"primaryKey;column:staff_id"`
	FirstName	string	`gorm:"size:40"`
	LastName	string	`gorm:"size:60"`
	AddressID	uint
	Address	*Address	`json:"-" gorm:"foreignKey:AddressID"`
	Email	string	`gorm:"size:100"`
	StoreID	uint
	Store	*Store	`json:"-" gorm:"foreignKey:StoreID"`
	Username	string	`gorm:"size:30"`
	Password	string	`gorm:"size:32"`
	LastUpdate	time.Time	`gorm:"autoUpdateTime"`
	Picture	[]byte	`gorm:"bytea"`
}

func (Staff) TableName() string {
	return "staff"
}