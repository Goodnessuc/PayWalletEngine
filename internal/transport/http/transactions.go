package http

import (
	"PayWalletEngine/internal/transactions"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type TransactionService interface {
	GetTransactionByTransactionID(ctx context.Context, transactionID int64) (*transactions.Transaction, error)
	GetTransactionsFromAccount(ctx context.Context, accountNumber int64) ([]transactions.Transaction, error)
	GetTransactionByReference(ctx context.Context, reference string) (*transactions.Transaction, error)
}

// GetTransactionByTransactionID handles the retrieval of a single transaction.
func (h *Handler) GetTransactionByTransactionID(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	stringID := vars["transaction_id"]
	if stringID == "" {
		http.Error(writer, "Transaction ID is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(stringID, 10, 64)
	if err != nil {
		http.Error(writer, "Invalid transaction ID format", http.StatusBadRequest)
		return
	}

	txn, err := h.Transaction.GetTransactionByTransactionID(request.Context(), id)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(txn)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

// GetTransactionsFromAccount handles the retrieval of all transactions made by a specific sender.
func (h *Handler) GetTransactionsFromAccount(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	senderAccountNumberStr := vars["account_number"]
	if senderAccountNumberStr == "" {
		http.Error(writer, "Account number is required", http.StatusBadRequest)
		return
	}

	senderAccountNumber, err := strconv.ParseInt(senderAccountNumberStr, 10, 64)
	if err != nil {
		http.Error(writer, "Invalid account number format", http.StatusBadRequest)
		return
	}

	txns, err := h.Transaction.GetTransactionsFromAccount(request.Context(), senderAccountNumber)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(txns)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

// GetTransactionByReference handles the retrieval of a single transaction by its reference number.
func (h *Handler) GetTransactionByReference(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	stringReference := vars["transaction_reference"]
	if stringReference == "" {
		http.Error(writer, "Reference number is required", http.StatusBadRequest)
		return
	}

	reference, err := strconv.ParseInt(stringReference, 10, 64)
	if err != nil {
		http.Error(writer, "Invalid reference number format", http.StatusBadRequest)
		return
	}

	txn, err := h.Transaction.GetTransactionByReference(request.Context(), strconv.FormatInt(reference, 10))
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(txn)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

// CreditAccount handles crediting an account for a transaction.
func (h *Handler) CreditAccount(writer http.ResponseWriter, request *http.Request) {
	var creditRequest struct {
		ReceiverAccountNumber int64   `json:"receiver_account_number"`
		Amount                float64 `json:"amount"`
		Description           string  `json:"description"`
		PaymentMethod         string  `json:"payment_method"`
	}

	err := json.NewDecoder(request.Body).Decode(&creditRequest)
	if err != nil {
		http.Error(writer, "Invalid request format", http.StatusBadRequest)
		return
	}

	txn, err := h.Transaction.CreditAccount(request.Context(), creditRequest.ReceiverAccountNumber, creditRequest.Amount, creditRequest.Description, creditRequest.PaymentMethod)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(txn)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

// DebitAccount handles debiting the specified account.
func (h *Handler) DebitAccount(writer http.ResponseWriter, request *http.Request) {
	var debitRequest struct {
		SenderAccountNumber int64   `json:"sender_account_number"`
		Amount              float64 `json:"amount"`
		Description         string  `json:"description"`
		PaymentMethod       string  `json:"payment_method"`
	}

	err := json.NewDecoder(request.Body).Decode(&debitRequest)
	if err != nil {
		http.Error(writer, "Invalid request format", http.StatusBadRequest)
		return
	}

	txn, err := h.Transaction.DebitAccount(request.Context(), debitRequest.SenderAccountNumber, debitRequest.Amount, debitRequest.Description, debitRequest.PaymentMethod)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(txn)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

// TransferFunds handles transferring funds by crediting and debiting specified users.
func (h *Handler) TransferFunds(writer http.ResponseWriter, request *http.Request) {
	var transferRequest struct {
		SenderAccountNumber   int64   `json:"sender_account_number"`
		ReceiverAccountNumber int64   `json:"receiver_account_number"`
		Amount                float64 `json:"amount"`
		Description           string  `json:"description"`
		PaymentMethod         string  `json:"payment_method"`
	}

	err := json.NewDecoder(request.Body).Decode(&transferRequest)
	if err != nil {
		http.Error(writer, "Invalid request format", http.StatusBadRequest)
		return
	}

	txn, err := h.Transaction.TransferFunds(request.Context(), transferRequest.SenderAccountNumber, transferRequest.ReceiverAccountNumber, transferRequest.Amount, transferRequest.Description, transferRequest.PaymentMethod)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(txn)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}
