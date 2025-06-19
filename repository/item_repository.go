package repository

import (
	"errors"
	"go-jwt-api/config"
	"go-jwt-api/model"

	"gorm.io/gorm"
)

// CreateItem adds a new item with transaction and optional name uniqueness check
func CreateItem(item *model.Item) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// Optional: Check for duplicate item name
		var existing model.Item
		if err := tx.Where("name = ?", item.Name).First(&existing).Error; err == nil {
			return errors.New("item with the same name already exists")
		}

		return tx.Create(item).Error
	})
}

// GetItemByID retrieves an item by its ID
func GetItemByID(id uint) (*model.Item, error) {
	var item model.Item
	err := config.DB.First(&item, id).Error
	return &item, err
}

// UpdateItem updates an item in a transaction
func UpdateItem(item *model.Item) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Save(item).Error
	})
}

// DeleteItem deletes an item by ID in a transaction
func DeleteItem(id uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&model.Item{}, id).Error
	})
}
