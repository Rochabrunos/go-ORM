package model

import "time"

type Store struct {
	ID            uint `gorm:"primaryKey;column:store_id"`
	ManageStaffID uint
	Staff         *Staff `json:"-" gorm:"foreignKey:ManageStaffID"`
	AddressID     uint
	Address       *Address   `json:"-" gorm:"foreignKey:AddressID"`
	LastUpdate    *time.Time `gorm:"autoUpdateTime"`
}

func (Store) TableName() string {
	return "store"
}
