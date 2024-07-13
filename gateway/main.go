package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"
	client "github.com/rasha-hantash/chariot-takehome/gateway/grpcClient"
)

type Config struct {
	DefaultPort string `env:"GATEWAY_PORT" envDefault:"8080"`
	Env         string `env:"ENVIRONMENT" envDefault:"local"`
	ApiAddr     string `env:"API_SERVICE_ADDR" envDefault:"localhost:9093"`
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

	// Initialize and start the gateway
	router := http.NewServeMux()
	
	// Add your routes here, for example:
	// router.HandleFunc("/users", handleUsers(grpcClient))

	serverAddr := fmt.Sprintf(":%s", cfg.DefaultPort)
	log.Printf("Starting server on %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}