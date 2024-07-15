package handlers

import (
	"encoding/json"
	"net/http"
	"sync"

	pb "github.com/rasha-hantash/chariot-takehome/api/grpc/proto"
	client "github.com/rasha-hantash/chariot-takehome/gateway/grpcClient"
)

var (
	idempotencyKeys = make(map[string]bool)
	mu              sync.Mutex
)

func CreateUserHandler(grpcClient *client.ApiClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pb.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := grpcClient.CreateUser(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

func CreateAccountHandler(grpcClient *client.ApiClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pb.CreateAccountRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		account, err := grpcClient.CreateAccount(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(account)
	}
}

func DepositFundsHandler(grpcClient *client.ApiClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pb.DepositFundsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		// Check for idempotency key
		if idempotencyKeys[req.IdempotencyKey] {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "Duplicate request detected"}`))
			return
		}

		transaction, err := grpcClient.DepositFunds(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		idempotencyKeys[req.IdempotencyKey] = true

		json.NewEncoder(w).Encode(transaction)
	}
}

func WithdrawFundsHandler(grpcClient *client.ApiClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pb.WithdrawFundsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		// Check for idempotency key
		if idempotencyKeys[req.IdempotencyKey] {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "Duplicate request detected"}`))
			return
		}

		// Call the WithdrawFunds method
		transaction, err := grpcClient.WithdrawFunds(&req)
		if err != nil {
			// Send error response
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Mark the idempotency key as used
		idempotencyKeys[req.IdempotencyKey] = true

		// Encode and send the transaction with status 200
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // Explicitly setting the status code to 200
		if err := json.NewEncoder(w).Encode(transaction); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func TransferFundsHandler(grpcClient *client.ApiClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pb.TransferFundsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		// Check for idempotency key
		if idempotencyKeys[req.IdempotencyKey] {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "Duplicate request detected"}`))
			return
		}

		transaction, err := grpcClient.TransferFunds(&req)
		if err != nil {
			// Send error response
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Mark the idempotency key as used
		idempotencyKeys[req.IdempotencyKey] = true

		// Encode and send the transaction with status 200
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // Explicitly setting the status code to 200
		if err := json.NewEncoder(w).Encode(transaction); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func ListTransactionsHandler(grpcClient *client.ApiClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pb.ListTransactionsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		transactions, err := grpcClient.ListTransactions(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // Explicitly setting the status code to 200
		if err := json.NewEncoder(w).Encode(transactions); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetAccountBalanceHandler(grpcClient *client.ApiClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pb.GetAccountBalanceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		balance, err := grpcClient.GetAccountBalance(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(balance)
	}
}
