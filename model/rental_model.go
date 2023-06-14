package model

import "time"

type Rental struct {
	ID          uint `gorm:"primaryKey;column:rental_id"`
	RentalDate  time.Time
	InventoryID uint
	Inventory   *Inventory `json:"-" gorm:"foreignKey:InventoryID"`
	CustomerID  uint
	Customer    *Customer `json:"-" gorm:"foreignKey:CustomerID"`
	ReturnDate  time.Time
	StaffID     uint
	Staff       *Staff    `json:"-" gorm:"foreignKey:StaffID"`
	LastUpdate  time.Time `gorm:"autoUpdateTime"`
}

func (Rental) TableName() string {
	return "rental"
}
