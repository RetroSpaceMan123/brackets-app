package handler

import (
	"encoding/json"
	"net/http"

	"example.com/api/app/model"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetAllUsers(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	users := []model.User{}
	db.Find(&users)
	respondJSON(w, http.StatusOK, users)
}

func SignUp(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	// Parse request into User instance
	user := model.User{}
	err := json.NewDecoder(r.Body).Decode(&user) // These lines changed slightly from tutorial
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		//w.WriteHeader(http.StatusBadRequest)
		respondError(w, http.StatusBadRequest, err.Error()) // The above line replaced with our respondError function created in common.go
		return
	}

	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)

	user.Password = string(hashedPassword) // Is changing the user object's password to the hashed version the best way to pass the data to the db?

	user.UUID = uuid.New().String() // This creates a unique UUID tied to the user account, allowing for changes to usernames and emails without loosing user data

	// These lines changed from tutorial, insert new user into database
	if err := db.Save(&user).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Respond with JSON
	respondJSON(w, http.StatusCreated, user)
}

func SignIn(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	// Parse request into User instance
	user := model.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		//w.WriteHeader(http.StatusBadRequest)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	//Get first user from database matching username
	vars := mux.Vars(r)
	username := vars["username"]
	storedUser := getUserOr404(db, username, w, r)
	if storedUser == nil {
		//respondError(w, http.StatusInternalServerError, err.Error())
		//Error called in getUserOr404 helper function I think
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		respondError(w, http.StatusUnauthorized, err.Error())
	}

	respondJSON(w, http.StatusOK, user)
}

func GetUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	username := vars["username"]
	user := getUserOr404(db, username, w, r)
	if user == nil {
		return
	}
	respondJSON(w, http.StatusOK, user)
}

func getUserOr404(db *gorm.DB, username string, w http.ResponseWriter, r *http.Request) *model.User {
	user := model.User{}
	if err := db.First(&user, model.User{Username: username}).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &user
}

func UpdateUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	username := vars["username"]
	user := getUserOr404(db, username, w, r)
	if user == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&user).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, user)
}

func DeleteUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	username := vars["username"]
	user := getUserOr404(db, username, w, r)
	if user == nil {
		return
	}
	if err := db.Delete(&user).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}
