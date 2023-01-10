package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	Id string `json:"id"`
	Email string `json:"email"`
	Status string `json:"status"`
}
type CreateUserPayload struct {
	Email string `json:"email"`
}
type DeleteUserPayload struct {
	Id    string `json:"id"`
}

type CreateUserResponse struct {
	Status string `json:"status"`
	Error string `json:"error"`
}
type DeleteUserResponse struct {
	Id string `json:"id"`
	Status string `json:"status"`
	Error string `json:"error"`
}

// add user to mongo
func (h *Handler) UserAdd(w http.ResponseWriter, r *http.Request) {
	// get the payload
	var payload CreateUserPayload
	var response CreateUserResponse
	var user User
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(payload)

	userCollection := h.db.Collection("users")
	user.Email = payload.Email
	// count user by email
	count, err := userCollection.CountDocuments(r.Context(), bson.M{"email": payload.Email})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotImplemented)
		response.Error = "Not implemented"
		json.NewEncoder(w).Encode(response)
		return
	}
	if count == 0 {
		user.Id = uuid.New().String()
		user.Status = "Active"
		_, err = userCollection.InsertOne(r.Context(), user)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotImplemented)
			response.Error = "Not implemented"
			json.NewEncoder(w).Encode(response)
			return
		}
		response.Status = "User created"
		w.WriteHeader(http.StatusCreated)
	} else {
		response.Status = "User already exist"
		w.WriteHeader(http.StatusConflict)
		return
	}
	json.NewEncoder(w).Encode(response)
}


// delete user from mongo
func (h *Handler) UserDelete(w http.ResponseWriter, r *http.Request) {
	// get the payload
	var payload DeleteUserPayload
	var response DeleteUserResponse
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		response.Error = "Bad request"
		json.NewEncoder(w).Encode(response)
		return
	}
	// delete user from mongo
	userCollection := h.db.Collection("users")
	_, err = userCollection.DeleteMany(r.Context(), payload)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotImplemented)
		response.Error = "Not implemented"
		json.NewEncoder(w).Encode(response)
		return
	}
	response.Id = payload.Id
	response.Status = "Deleted"
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(response)
}

