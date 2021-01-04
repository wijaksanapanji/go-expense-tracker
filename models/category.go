package models

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/wijaksanapanji/go-expense-tracker/config"
)

type Category struct {
	CommonFields
	Name        string        `json:"name"`
	Transaction []Transaction `json:"transactions"`
}

func GetAllCategories(w http.ResponseWriter, r *http.Request) {
	var categories []Category
	config.DBConnection.Find(&categories)
	render.JSON(w, r, categories)
}

func AddCategory(w http.ResponseWriter, r *http.Request) {
	var category Category
	json.NewDecoder(r.Body).Decode(&category)
	config.DBConnection.Create(&category)
	render.JSON(w, r, category)
}
