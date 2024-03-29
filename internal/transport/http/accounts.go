package http

import (
	"PayWalletEngine/internal/accounts"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// CreateAccount decodes an Account object from the HTTP request body and then tries to create a new account in the database using the CreateAccount method of the AccountService interface. If the account is successfully created, it encodes and sends the created account as a response.
func (h *Handler) CreateAccount(writer http.ResponseWriter, request *http.Request) {
	// Decode the account from the request body
	var acct accounts.Account
	if err := json.NewDecoder(request.Body).Decode(&acct); err != nil {
		http.Error(writer, "Failed to decode request body", http.StatusBadRequest)
		log.Println("Failed to decode request body:", err)
		return
	}

	// Create the account in the database
	if err := h.Accounts.CreateAccount(request.Context(), &acct); err != nil {
		http.Error(writer, fmt.Sprintf("Failed to create account: %v", err), http.StatusInternalServerError)
		log.Println("Failed to create account:", err)
		return
	}

	// Encode and send the created account as a response
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(writer).Encode(acct); err != nil {
		http.Error(writer, "Failed to encode response", http.StatusInternalServerError)
		log.Println("Failed to encode response:", err)
	}
}

// GetAccountByID extracts the id from the URL parameters and then fetches the account with that id from the database using the GetAccountByID method of the AccountService interface. If the account is found, it encodes and sends the account as a response.
func (h *Handler) GetAccountByID(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	account, err := h.Accounts.GetAccountByID(request.Context(), uint(id))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(writer).Encode(account); err != nil {
		log.Panicln(err)
	}

}

// GetAccountByNumber extracts the number from the URL parameters and then fetches the account with that number from the database using the GetAccountByNumber method of the AccountService interface. If the account is found, it encodes and sends the account as a response.
func (h *Handler) GetAccountByNumber(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	number, err := strconv.ParseInt(vars["number"], 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	account, err := h.Accounts.GetAccountByNumber(request.Context(), uint(number))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(writer).Encode(account); err != nil {
		log.Panicln(err)
	}
}

// UpdateAccountDetails decodes an AccountNumber object from the HTTP request body and then updates the account in the database using the UpdateAccountDetails method of the AccountService interface. If the account is successfully updated, it responds with a status code 200 OK.
func (h *Handler) UpdateAccountDetails(writer http.ResponseWriter, request *http.Request) {
	var acct accounts.Account
	if err := json.NewDecoder(request.Body).Decode(&acct); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.Accounts.UpdateAccountDetails(request.Context(), acct)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(writer).Encode(acct); err != nil {
		log.Panicln(err)
	}

	writer.WriteHeader(http.StatusOK)

}

func (h *Handler) GetUserDetailsByAccountNumber(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	accountNumber, err := strconv.ParseInt(vars["account_number"], 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch the user details by account number
	user, err := h.Accounts.GetUserByAccountNumber(request.Context(), uint(accountNumber))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode and send the user details as a response
	if err := json.NewEncoder(writer).Encode(user); err != nil {
		log.Panicln(err)
	}
	writer.WriteHeader(http.StatusOK)
}

func (h *Handler) GetAccountsByUserID(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	userID, err := strconv.ParseInt(vars["user_id"], 10, 64)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	// Fetch the accounts by user TransactionID
	accounts, err := h.Accounts.GetAccountsByUserID(request.Context(), uint(userID))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode and send the accounts as a response
	if err := json.NewEncoder(writer).Encode(accounts); err != nil {
		http.Error(writer, "Failed to encode accounts", http.StatusInternalServerError)
		return
	}
}
