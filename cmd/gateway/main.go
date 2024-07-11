package main

import (
	"log"
	"myproject/internal/gateway"
	"myproject/internal/grpc/client"
	"net/http"
)

type Config struct {
	DefaultPort string `env:"GATEWAY_PORT" envDefault:"8080"`
	Env         string `env:"ENVIRONMENT" envDefault:"local"`
}

type ApiAddrs struct {
	ApiAddr string `env:"API_SERVICE_ADDR" envDefault:"localhost:9093"`
}

func main() {
	// Initialize gRPC clients
	grpcClient, err := client.NewGrpcClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to create user client: %v", err)
	}

	// Initialize and start the gateway
	router := gateway.NewRouter(grpcClient)
	log.Fatal(http.ListenAndServe(":8080", router))
}

// File: internal/gateway/router.go
