package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"go-django/internal/database"
	"go-django/internal/grpcclient"
	"go-django/internal/models"
	"go-django/internal/pb"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	grpClient grpcclient.GrpClientInterface
}

func NewUserController(grpClient grpcclient.GrpClientInterface) *UserController {
	return &UserController{grpClient: grpClient}
}

func (uc *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, email FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []models.User{}

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	response := models.ResultResponse{
		RESULTS: users,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user models.User
	err = database.DB.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword)
	result, err := database.DB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", user.Name, user.Email, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.ID = int(id)
	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)

	// Call gRPC to notify Django service
	grpcUser := &pb.User{
		Id:    int32(user.ID),
		Name:  user.Name,
		Email: user.Email,
	}
	log.Println(grpcUser)
	_, err = uc.grpClient.CreateUser(context.Background(), grpcUser)
	if err != nil {
		log.Println(err.Error(), http.StatusInternalServerError)
	}
}

func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user models.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		user.Password = string(hashedPassword)
		_, err = database.DB.Exec(
			`UPDATE users
            SET name = ?, email = ?, password = ?
            WHERE id = ?`,
			user.Name,
			user.Email,
			user.Password,
			id,
		)
	} else {
		_, err = database.DB.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", user.Name, user.Email, id)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	// Call gRPC to notify Django service
	grpcUser := &pb.User{
		Id:       int32(id),
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
	log.Println(grpcUser)
	_, err = uc.grpClient.UpdateUser(context.Background(), grpcUser)
	if err != nil {
		log.Println(err.Error(), http.StatusInternalServerError)
	}

}

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	// Call gRPC to notify Django service
	_, err = uc.grpClient.DeleteUser(context.Background(), int32(id))
	if err != nil {
		log.Println(err.Error(), http.StatusInternalServerError)
	}

}
