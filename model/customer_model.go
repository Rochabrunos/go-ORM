package model

import "time"

type Customer struct {
	ID         uint `gorm:"primaryKey;column:customer_id"`
	StoreID    uint
	Store      *Store `gorm:"foreignKey:StoreID"`
	FirstName  string `gorm:"size:40"`
	LastName   string `gorm:"size:60"`
	Email      string `gorm:"size:100"`
	AddressID  uint
	Address    *Address `json:"-" gorm:"foreignKey:AddressID"`
	Activebool bool
	CreateDate time.Time `gorm:"autoCreateTime"`
	LastUpdate time.Time `gorm:"autoUpdateTime"`
	Active     bool
}

func (Customer) TableName() string {
	return "customer"
}
