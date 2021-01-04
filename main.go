package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/wijaksanapanji/go-expense-tracker/config"
	"github.com/wijaksanapanji/go-expense-tracker/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	config.TokenAuth = jwtauth.New("HS256", []byte(config.SecretKey), nil)

	config.DBConnection, config.DBError = gorm.Open(sqlite.Open("development.db"), &gorm.Config{})
	if config.DBError != nil {
		panic("Failed to connecting to database!")
	}

	config.DBConnection.AutoMigrate(&models.Transaction{})
	config.DBConnection.AutoMigrate(&models.Category{})
	config.DBConnection.AutoMigrate(&models.User{})

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
		r.Get("/", models.GetAllCategories)
		r.Post("/", models.AddCategory)
	})

	r.Route("/users", func(r chi.Router) {
		r.Get("/", models.GetAllUser)
		r.Post("/register", models.RegisterUser)
		r.Post("/login", models.LoginUser)

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(config.TokenAuth))
			r.Use(jwtauth.Authenticator)

			r.Get("/profile", models.ProfileUser)
		})
	})

	// DEV - For Development Purposes
	r.Post("/reset", func(w http.ResponseWriter, r *http.Request) {
		config.DBConnection.Exec("DELETE FROM transactions")
		config.DBConnection.Exec("DELETE FROM categories")
		config.DBConnection.Exec("DELETE FROM users")
		render.JSON(w, r, struct {
			Message string `json:"message"`
		}{
			Message: "Succesfully reset database!",
		})
	})

	port := ":8000"
	fmt.Println("Server listening on http://localhost" + port)
	http.ListenAndServe(port, r)
}
