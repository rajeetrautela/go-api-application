package model

import "gorm.io/gorm"

type RefreshToken struct {
	gorm.Model
	Token    string `gorm:"uniqueIndex"`
	Username string
}
