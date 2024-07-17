package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"sync"

	pb "github.com/rasha-hantash/chariot-takehome/api/grpc/proto"
	client "github.com/rasha-hantash/chariot-takehome/gateway/grpcClient"
)

var (
	idempotencyKeys = make(map[string]bool)
	mu              sync.Mutex
)

func CreateUserHandler(ctx context.Context, grpcClient *client.ApiClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pb.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := grpcClient.CreateUser(ctx, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

func CreateAccountHandler(ctx context.Context, grpcClient *client.ApiClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req pb.CreateAccountRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		account, err := grpcClient.CreateAccount(ctx, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(account)
	}
}

func DepositFundsHandler(ctx context.Context, grpcClient *client.ApiClient) http.HandlerFunc {
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

		transaction, err := grpcClient.DepositFunds(ctx, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		idempotencyKeys[req.IdempotencyKey] = true

		json.NewEncoder(w).Encode(transaction)
	}
}

func WithdrawFundsHandler(ctx context.Context, grpcClient *client.ApiClient) http.HandlerFunc {
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
		transaction, err := grpcClient.WithdrawFunds(ctx, &req)
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

func TransferFundsHandler(ctx context.Context, grpcClient *client.ApiClient) http.HandlerFunc {
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

		transaction, err := grpcClient.TransferFunds(ctx, &req)
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

func ListTransactionsHandler(ctx context.Context, grpcClient *client.ApiClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := r.URL.Query().Get("account_id")
		cursor := r.URL.Query().Get("cursor")
		limit := r.URL.Query().Get("limit")
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			slog.ErrorContext(ctx, "error parsing limit", "error", err, "limit", limit)
			http.Error(w, "invalid limit value", http.StatusBadRequest)
			return
		}

		// Validate required query parameters
		if accountID == "" {
			http.Error(w, "missing required query parameters: account_id and at_time", http.StatusBadRequest)
			return
		}
		req := pb.ListTransactionsRequest{
			AccountId: accountID,
			Cursor:    cursor,
			Limit:     int32(limitInt),
		}

		transactions, err := grpcClient.ListTransactions(ctx, &req)
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

// GetAccountBalanceHandler handles the request for getting account balance
func GetAccountBalanceHandler(ctx context.Context, grpcClient *client.ApiClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract query parameters
		accountID := r.URL.Query().Get("account_id")
		slog.Info("query param", "accountID", accountID)

		// Validate required query parameters
		if accountID == "" {
			http.Error(w, "missing required query parameters: account_id and at_time", http.StatusBadRequest)
			return
		}

		// Create the request object
		req := pb.GetAccountBalanceRequest{
			AccountId: accountID,
			// AtTime:    atTime,
		}

		// Call the gRPC client
		balance, err := grpcClient.GetAccountBalance(ctx, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Encode the response as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // Explicitly setting the status code to 200
		if err := json.NewEncoder(w).Encode(balance); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
