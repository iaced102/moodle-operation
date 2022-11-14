package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Id string `json:"id"`
	Status string `json:"status"`
}
type CreateUserPayload struct {
	Id    string `json:"id"`
}
type DeleteUserPayload struct {
	Id    string `json:"id"`
}

type CreateUserResponse struct {
	Id string `json:"id"`
	Status string `json:"status"`
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

	user.Id = payload.Id
	user.Status = "Active"
	// add user to mongo
	userCollection := h.db.Collection("users")
	_, err = userCollection.InsertOne(context.Background(), user)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	response.Id = payload.Id
	response.Status = "Created"
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}


// delete user from mongo
func (h *Handler) UserDelete(w http.ResponseWriter, r *http.Request) {
	// get the payload
	var payload DeleteUserPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// delete user from mongo
	userCollection := h.db.Collection("users")
	_, err = userCollection.DeleteMany(r.Context(), payload)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	var response CreateUserResponse
	response.Id = payload.Id
	response.Status = "Deleted"
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(response)
}

