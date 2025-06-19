package model

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}
