package repository_test

import (
	"errors"
	"go-jwt-api/model"
	"go-jwt-api/repository"
	"testing"

	"gorm.io/gorm"
)

func TestRepositoryCRUD_first(t *testing.T) {
	// Create item
	item := &model.Item{Name: "TestItem", Price: 100}
	err := repository.CreateItem(item)
	if err != nil {
		t.Fatalf("Failed to create item: %v", err)
	}

	// Try duplicate name
	dup := &model.Item{Name: "TestItem", Price: 200}
	err = repository.CreateItem(dup)
	if err == nil || !errors.Is(err, gorm.ErrDuplicatedKey) && err.Error() != "item with the same name already exists" {
		t.Fatalf("Expected duplicate name error, got: %v", err)
	}

	// Get item
	found, err := repository.GetItemByID(item.ID)
	if err != nil {
		t.Fatalf("Failed to get item: %v", err)
	}
	if found.Name != item.Name || found.Price != item.Price {
		t.Errorf("Mismatch on fetch. Want: %+v, Got: %+v", item, found)
	}

	// Update item
	found.Price = 999
	err = repository.UpdateItem(found)
	if err != nil {
		t.Fatalf("Failed to update item: %v", err)
	}
	updated, _ := repository.GetItemByID(found.ID)
	if updated.Price != 999 {
		t.Errorf("Expected updated price 999, got %d", updated.Price)
	}

	// Delete item
	err = repository.DeleteItem(found.ID)
	if err != nil {
		t.Fatalf("Failed to delete item: %v", err)
	}

	// Verify deletion
	_, err = repository.GetItemByID(found.ID)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Expected ErrRecordNotFound after deletion, got: %v", err)
	}
}
