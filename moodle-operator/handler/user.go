package handler

import (
	"encoding/json"
	"fmt"
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
}
type DeleteUserResponse struct {
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

	userCollection := h.db.Collection("users")
	user.Email = payload.Email
	// count user by email
	count, err := userCollection.CountDocuments(r.Context(), bson.M{"email": payload.Email})
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	if count == 0 {
		user.Id = uuid.New().String()
		user.Status = "Active"
		_, err = userCollection.InsertOne(r.Context(), user)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		response.Status = "User created"
	} else {
		response.Status = "User already exist"
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// get all users from mongo



// 	user.Id = uuid.New().String()
// 	user.Status = "Active"
// 	// add user to mongo
// 	userCollection := h.db.Collection("users")
// 	_, err = userCollection.InsertOne(context.Background(), user)
// 	if err != nil {
// 		w.Write([]byte(err.Error()))
// 	}
// 	response.Status = "Active"
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(response)
// }


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
	var response DeleteUserResponse
	response.Id = payload.Id
	response.Status = "Deleted"
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(response)
}

