package models

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/wijaksanapanji/go-expense-tracker/config"
)

type User struct {
	CommonFields
	Name        string        `json:"name"`
	Email       string        `json:"email"`
	Password    string        `json:"password"`
	Transaction []Transaction `json:"transactions"`
}

type UserResponse struct {
	CommonFields
	Name  string `json:"name"`
	Email string `json:"email"`
}

// DEV - For Development Purposes
func GetAllUser(w http.ResponseWriter, r *http.Request) {
	var users []User
	config.DBConnection.Find(&users)
	render.JSON(w, r, users)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	var userExist User
	exist := config.DBConnection.Where(&User{Email: user.Email}).First(&userExist)
	if exist.RowsAffected > 0 {
		render.JSON(w, r, struct {
			Message string `json:"message"`
		}{
			Message: "User with that email already exist!",
		})
		return
	}

	hashedPassword, err := config.HashPassword(user.Password)
	if err != nil {
		render.JSON(w, r, struct {
			Message string `json:"message"`
		}{
			Message: "Failed to create user",
		})
		return
	}

	user.Password = hashedPassword
	config.DBConnection.Create(&user)
	render.JSON(w, r, &UserResponse{
		Name:         user.Name,
		Email:        user.Email,
		CommonFields: user.CommonFields,
	})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var request map[string]string
	json.NewDecoder(r.Body).Decode(&request)

	var user User
	result := config.DBConnection.Where(&User{Email: request["email"]}).First(&user)
	if result.Error != nil {
		render.JSON(w, r, struct {
			Message string `json:"message"`
		}{
			Message: "User not found",
		})
		return
	}

	match := config.CheckHashPassword(request["password"], user.Password)
	if match {
		_, token, _ := config.TokenAuth.Encode(map[string]interface{}{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		})

		render.JSON(w, r, struct {
			Token string `json:"token"`
		}{
			Token: token,
		})
		return
	}

	render.JSON(w, r, struct {
		Message string `json:"message"`
	}{
		Message: "Wrong password",
	})

}

func ProfileUser(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	var user User
	user.ID = uint(claims["id"].(float64))
	config.DBConnection.First(&user)
	render.JSON(w, r, &UserResponse{
		Name:         user.Name,
		Email:        user.Email,
		CommonFields: user.CommonFields,
	})
}
