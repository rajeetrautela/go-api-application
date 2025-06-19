package model

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	ID    string `json:"id"`
	Name  string `json:"name" gorm:"uniqueIndex"`
	Price int    `json:"price"`
}
