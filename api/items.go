package api

import (
	"encoding/json"
	"go-jwt-api/config"
	"go-jwt-api/model"
	"go-jwt-api/repository"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func getItems(w http.ResponseWriter, r *http.Request) {
	var items []model.Item
	if err := config.DB.Find(&items).Error; err != nil {
		http.Error(w, "Failed to fetch items", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(items)
}

func getItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	item, err := repository.GetItemByID(uint(id))
	if err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(item)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	var item model.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := repository.CreateItem(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(item)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updated model.Item
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updated.ID = uint(id)

	if err := repository.UpdateItem(&updated); err != nil {
		http.Error(w, "Failed to update item", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(updated)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := repository.DeleteItem(uint(id)); err != nil {
		http.Error(w, "Failed to delete item", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
