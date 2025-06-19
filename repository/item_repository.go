package repository

import (
	"go-jwt-api/config"
	"go-jwt-api/model"
)

func CreateItem(item *model.Item) error {
	return config.DB.Create(item).Error
}

func GetItemByID(id uint) (*model.Item, error) {
	var item model.Item
	err := config.DB.First(&item, id).Error
	return &item, err
}

func UpdateItem(item *model.Item) error {
	return config.DB.Save(item).Error
}

func DeleteItem(id uint) error {
	return config.DB.Delete(&model.Item{}, id).Error
}
