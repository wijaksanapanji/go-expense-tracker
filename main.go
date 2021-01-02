package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Transaction []Transaction
}

type TransactionType string

const (
	income = iota
	expense
)

type Transaction struct {
	gorm.Model
	Type        TransactionType `json:"type"`
	Description string          `json:"description"`
	Date        time.Time       `json:"date"`
	Amount      int             `json:"amount"`
	CategoryID  int             `json:"category_id"`
	UserID      int             `json:"user_id"`
}

type Category struct {
	gorm.Model
	Name        string `json:"name"`
	Transaction []Transaction
}

var db gorm.DB

func main() {
	db, err := gorm.Open(sqlite.Open("development.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connecting to database!")
	}

	db.AutoMigrate(&Transaction{})
	db.AutoMigrate(&Category{})
	db.AutoMigrate(&User{})

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, struct {
			Author string `json:"author"`
			About  string `json:"about"`
		}{
			Author: "wijaksanapanji",
			About:  "Rest Expense tracker written in golang",
		})
	})

	r.Route("/categories", func(r chi.Router) {
		r.Get("/", allCategory)
		r.Post("/", addCategory)
	})

	http.ListenAndServe(":8000", r)
}

func allCategory(w http.ResponseWriter, r *http.Request) {
	var categories []Category
}

func addCategory(w http.ResponseWriter, r *http.Request) {

}
