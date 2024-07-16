package main

import (
	// "fmt"
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
	client "github.com/rasha-hantash/chariot-takehome/gateway/grpcClient"
	h "github.com/rasha-hantash/chariot-takehome/gateway/handlers"
)

type Config struct {
	DefaultPort string `env:"GATEWAY_PORT" envDefault:"8080"`
	Env         string `env:"ENVIRONMENT" envDefault:"local"`
	ApiAddr     string `env:"API_SERVICE_ADDR" envDefault:"localhost:9093"`
}

func main() {
	ctx := context.Background()
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

	// Initialize rate limiter todo: look more into bursts 
	limiter := rate.NewLimiter(rate.Limit(2), 10) // 2 requests per second, allow bursts up to 10

	// Add rate limiter middleware
	router.Use(rateLimiterMiddleware(limiter))

	// Add health check endpoint
	// Tools like Pingdom, UptimeRobot, or StatusCake can periodically send HTTP requests to your /health endpoint.
	router.HandleFunc("/health", healthCheckHandler).Methods("GET")

	// Add your routes here, for example:
	router.HandleFunc("/create_user", h.CreateUserHandler(ctx, grpcClient)).Methods("POST")
	router.HandleFunc("/create_account", h.CreateAccountHandler(ctx, grpcClient)).Methods("POST")
	router.HandleFunc("/deposit_funds", h.DepositFundsHandler(ctx, grpcClient)).Methods("POST")
	router.HandleFunc("/withdraw_funds", h.WithdrawFundsHandler(ctx, grpcClient)).Methods("POST")
	router.HandleFunc("/transfer_funds", h.TransferFundsHandler(ctx, grpcClient)).Methods("POST")
	router.HandleFunc("/list_transactions", h.ListTransactionsHandler(ctx, grpcClient)).Methods("GET")
	router.HandleFunc("/get_account_balance", h.GetAccountBalanceHandler(ctx, grpcClient)).Methods("GET")

	log.Println("Gateway server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}


func rateLimiterMiddleware(limiter *rate.Limiter) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
