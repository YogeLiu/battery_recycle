package models

import "time"

type Seller struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:100;not null"`
	Phone     string    `json:"phone" gorm:"size:20;not null"`
	Address   string    `json:"address" gorm:"size:255;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Seller) TableName() string {
	return "sellers"
}