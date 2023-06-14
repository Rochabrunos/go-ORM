package model

import "time"

type Payment struct {
	ID          uint `gorm:"primaryKey;column:payment_id"`
	CustomerID  uint
	Customer    *Customer `json:"-" gorm:"foreignKey:CustomerID"`
	StaffID     uint
	Staff       *Staff `json:"-" gorm:"foreignKey:StaffID"`
	RentalID    uint
	Rental      *Rental `gorm:"foreignKey:RentalID"`
	Amount      float32
	PaymentDate time.Time
}

func (Payment) TableName() string {
	return "payment"
}
