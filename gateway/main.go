package main

import (
	// "fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/gorilla/mux"
	client "github.com/rasha-hantash/chariot-takehome/gateway/grpcClient"
	h "github.com/rasha-hantash/chariot-takehome/gateway/handlers"
)

type Config struct {
	DefaultPort string `env:"GATEWAY_PORT" envDefault:"8080"`
	Env         string `env:"ENVIRONMENT" envDefault:"local"`
	ApiAddr     string `env:"API_SERVICE_ADDR" envDefault:"localhost:9093"`
}

// GatewayHandler wraps the ApiClient and implements http.Handler
type GatewayHandler struct {
	apiClient *client.ApiClient
}

func main() {
	// Parse environment variables
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}

	// Initialize gRPC client
	grpcClient, err := client.NewApiClient(cfg.ApiAddr)
	if err != nil {
		log.Fatalf("Failed to create API gRPC client: %v", err)
	}
	defer grpcClient.Conn.Close()

	slog.Info("Gateway server started", "port", cfg.DefaultPort, "env", cfg.Env, "api_addr", cfg.ApiAddr)

	// Initialize and start the gateway
	router := mux.NewRouter()

	// Add your routes here, for example:
	router.HandleFunc("/create_user", h.CreateUserHandler(grpcClient)).Methods("POST")
	router.HandleFunc("/create_account", h.CreateAccountHandler(grpcClient)).Methods("POST")
	router.HandleFunc("/deposit_funds", h.DepositFundsHandler(grpcClient)).Methods("POST")
	router.HandleFunc("/withdraw_funds", h.WithdrawFundsHandler(grpcClient)).Methods("POST")
	router.HandleFunc("/transfer_funds", h.TransferFundsHandler(grpcClient)).Methods("POST")
	router.HandleFunc("/list_transactions", h.ListTransactionsHandler(grpcClient)).Methods("GET")
	router.HandleFunc("/get_account_balance", h.GetAccountBalanceHandler(grpcClient)).Methods("GET")

	log.Println("Gateway server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
