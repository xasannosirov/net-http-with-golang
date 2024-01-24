package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"test/netHttp/models"
	"test/netHttp/storage"

	"github.com/google/uuid"
)

func main() {
	http.HandleFunc("/user/create", CreateUser)
	http.HandleFunc("/user/all", GetAllUsers)
	http.HandleFunc("/user/update", UpdateUser)
	http.HandleFunc("/user/delete", DeleteUser)
	http.HandleFunc("/user/get", GetUser)
	log.Println("Server is running...")
	if err := http.ListenAndServe("localhost:8088", nil); err != nil {
		fmt.Println("Error while running server", err)
	}
}

// this function creates user body of user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error while getting body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user *models.User
	if err = json.Unmarshal(bodyByte, &user); err != nil {
		log.Println("error while unmarshalling body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := uuid.NewString()
	user.Id = id

	respUser, err := storage.CreateUser(user)
	if err != nil {
		log.Println("error while creating user", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	respBody, err := json.Marshal(respUser)
	if err != nil {
		log.Println("error while marshalling to response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(respBody)
}

// this function updates user, gets params id and body of user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error while getting body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user *models.User
	if err = json.Unmarshal(bodyByte, &user); err != nil {
		log.Println("error while unmarshalling body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user_id := r.URL.Query().Get("id")

	respUser, err := storage.UpdateUser(user_id, user)
	if err != nil {
		log.Println("error while updating user", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	respBody, err := json.Marshal(respUser)
	if err != nil {
		log.Println("error while marshalling to response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

// this function deleted user with params id
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	user_id := r.URL.Query().Get("id")

	if err := storage.DeleteUser(user_id); err != nil {
		log.Println("error while deleting user", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deleted User"))
}

// this function gets user with params id
func GetUser(w http.ResponseWriter, r *http.Request) {
	user_id := r.URL.Query().Get("id")

	respUser, err := storage.GetUser(user_id)
	if err != nil {
		log.Println("Error while getting user", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respBody, err := json.Marshal(respUser)
	if err != nil {
		log.Println("error while marshalling body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

// this function get users with params page and params limit
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")

	intPage, err := strconv.Atoi(page)
	if err != nil {
		log.Println("Error while converting page")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	limit := r.URL.Query().Get("limit")

	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		log.Println("Error while converting limit")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := storage.GetAllUsers(intPage, intLimit)
	if err != nil {
		log.Println("Error while getting all users", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respBody, err := json.Marshal(users)
	if err != nil {
		log.Println("error while marshalling body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}
