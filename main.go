package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type CommonFields struct {
	ID        uint       `gorm:"primaryKey, index" json:"id"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
}

type User struct {
	CommonFields
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"-"`
	Transaction []Transaction
}

type TransactionType string

const (
	income = iota
	expense
)

type Transaction struct {
	CommonFields
	Type        TransactionType `json:"type"`
	Description string          `json:"description"`
	Date        time.Time       `json:"date"`
	Amount      int             `json:"amount"`
	CategoryID  int             `json:"category_id"`
	UserID      int             `json:"user_id"`
}

type Category struct {
	CommonFields
	Name        string        `json:"name"`
	Transaction []Transaction `json:"transactions"`
}

var db *gorm.DB
var dbError error
var tokenAuth *jwtauth.JWTAuth

const (
	secretKey string = "SECRET_KEY"
)

func main() {
	tokenAuth = jwtauth.New("HS256", []byte(secretKey), nil)

	db, dbError = gorm.Open(sqlite.Open("development.db"), &gorm.Config{})
	if dbError != nil {
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

	r.Route("/users", func(r chi.Router) {
		r.Get("/", allUser)
		r.Post("/register", registerUser)
		r.Post("/login", loginUser)

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator)

			r.Get("/profile", profileUser)
		})
	})

	// DEV - For Development Purposes
	r.Post("/reset", func(w http.ResponseWriter, r *http.Request) {
		db.Exec("DELETE FROM transactions")
		db.Exec("DELETE FROM categories")
		db.Exec("DELETE FROM users")
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

func allCategory(w http.ResponseWriter, r *http.Request) {
	var categories []Category
	db.Find(&categories)
	render.JSON(w, r, categories)
}

func addCategory(w http.ResponseWriter, r *http.Request) {
	var category Category
	json.NewDecoder(r.Body).Decode(&category)
	db.Create(&category)
	render.JSON(w, r, category)
}

// DEV - For Development Purposes
func allUser(w http.ResponseWriter, r *http.Request) {
	var users []User
	db.Find(&users)
	render.JSON(w, r, users)
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	// TODO - Hash password
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	db.Create(&user)
	render.JSON(w, r, user)
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	// TODO - Implement proper login method
	var request map[string]string
	json.NewDecoder(r.Body).Decode(&request)

	var user User
	result := db.Where(&User{Email: request["email"], Password: request["password"]}).First(&user)
	if result.Error != nil {
		render.JSON(w, r, struct {
			Message string `json:"message"`
		}{
			Message: "User not found",
		})
		return
	}

	_, token, _ := tokenAuth.Encode(map[string]interface{}{"id": user.ID})

	render.JSON(w, r, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}

func profileUser(w http.ResponseWriter, r *http.Request) {
	// TODO - Get user by ID
	_, claims, _ := jwtauth.FromContext(r.Context())

	var user User
	user.ID = uint(claims["id"].(float64))
	db.First(&user)

	render.JSON(w, r, user)
}
