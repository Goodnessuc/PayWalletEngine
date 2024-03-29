package http

import (
	"PayWalletEngine/internal/users"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// CreateUser decodes a User object from the HTTP request body and then tries to create a new user in the database using the CreateUser method of the UserService interface. If the user is successfully created, it encodes and sends the created user as a response.
func (h *Handler) CreateUser(writer http.ResponseWriter, request *http.Request) {
	var u users.User
	if err := json.NewDecoder(request.Body).Decode(&u); err != nil {
		http.Error(writer, "Failed to decode request body", http.StatusBadRequest)
		log.Println("Failed to decode request body:", err)
		return
	}

	err := h.Users.CreateUser(request.Context(), &u)
	if err != nil {
		http.Error(writer, "Failed to create user", http.StatusInternalServerError)
		log.Println("Failed to create user:", err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(writer).Encode(u); err != nil {
		log.Panicln("Failed to encode response:", err)
	}
}

// GetUserByID extracts the id from the URL parameters and then fetches the user with that id from the database using the GetUserByID method of the UserService interface. If the user is found, it encodes and sends the user as a response.
func (h *Handler) GetUserByID(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	u, err := h.Users.GetUserByID(request.Context(), id)
	if err != nil {
		log.Println(err)
		return
	}
	if err := json.NewEncoder(writer).Encode(u); err != nil {
		log.Panicln(err)
	}
}

func (h *Handler) GetByEmail(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	email := vars["email"]

	// Check if email is valid
	if !isValidEmail(email) {
		http.Error(writer, "Invalid email format", http.StatusBadRequest)
		return
	}

	u, err := h.Users.GetByEmail(request.Context(), email)
	if err != nil {
		log.Println(err)
		http.Error(writer, "Failed to fetch user by email", http.StatusInternalServerError)
		return
	}

	// Check if user exists
	if u == nil {
		http.Error(writer, "User not found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(writer).Encode(u); err != nil {
		log.Println("Failed to encode user data: ", err)
		http.Error(writer, "Failed to process user data", http.StatusInternalServerError)
	}
}

// GetByUsername extracts the username from the URL parameters and then fetches the user with that username from the database using the GetByUsername method of the UserService interface. If the user is found, it encodes and sends the user as a response.
func (h *Handler) GetByUsername(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	username := vars["username"]
	u, err := h.Users.GetByUsername(request.Context(), username)
	if err != nil {
		log.Println(err)
		return
	}
	if err := json.NewEncoder(writer).Encode(u); err != nil {
		log.Panicln(err)
	}
}

// UpdateUser updates a user by TransactionID.
func (h *Handler) UpdateUser(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	// Decode the request body to get the updated user information
	var u users.User
	if err := json.NewDecoder(request.Body).Decode(&u); err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update user
	err = h.Users.UpdateUser(request.Context(), u, uint(id))
	if err != nil {
		http.Error(writer, "Failed to update user", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Encode and send response
	if err := json.NewEncoder(writer).Encode(u); err != nil {
		http.Error(writer, "Failed to encode response", http.StatusInternalServerError)
		log.Panicln(err)
	}
}

// ChangeUserStatus extracts the id from the URL parameters and then deletes the user with that id from the database using the ChangeUserStatus method of the UserService interface. If the user is successfully deleted, it sends a No Content status code as a response.
func (h *Handler) ChangeUserStatus(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	// Decode the request body to get the updated user information
	var u users.User
	if err := json.NewDecoder(request.Body).Decode(&u); err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update user
	err = h.Users.ChangeUserStatus(request.Context(), u, uint(id))
	if err != nil {
		http.Error(writer, "Failed to update user", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if err := json.NewEncoder(writer).Encode(map[string]string{"status": "OK"}); err != nil {
		log.Panicln(err)
	}

}

func (h *Handler) Ping(writer http.ResponseWriter, request *http.Request) {
	err := h.Users.ReadyCheck(request.Context())
	if err != nil {
		log.Println(err)
		return
	}
	if err := json.NewEncoder(writer).Encode(map[string]string{"status": "OK"}); err != nil {
		log.Panicln(err)
	}
}

// ResetPassword decodes a User object from the HTTP request body, then attempts to reset the user's password in the database using the ResetPassword method of the UserService interface. If the password is successfully reset, it sends an OK response.
func (h *Handler) ResetPassword(writer http.ResponseWriter, request *http.Request) {

	// Decode the request body to get the updated user information
	var u users.User
	if err := json.NewDecoder(request.Body).Decode(&u); err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Reset password
	err := h.Users.ResetPassword(request.Context(), u)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)

		return
	}
	// Send response
	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(map[string]string{"status": "Password reset successful"}); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)

	}
}
