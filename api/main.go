package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	_ "github.com/lib/pq"
	"github.com/rasha-hantash/chariot-takehome/api/pkgs/logger"
	"github.com/rasha-hantash/chariot-takehome/api/pkgs/postgres"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/caarlos0/env/v6"

	service "github.com/rasha-hantash/chariot-takehome/api/grpc"
	pb "github.com/rasha-hantash/chariot-takehome/api/grpc/proto"
	"github.com/rasha-hantash/chariot-takehome/api/grpc/repository"
)

type DatabaseConfig struct {
	ConnString string `env:"CONN_STRING" envDefault:"postgresql://postgres:postgres@localhost:5438/?sslmode=disable"`
	User       string `env:"DB_USER" envDefault:""`
	Port       string `env:"DB_PORT" envDefault:""`
	Host       string `env:"DB_HOST" envDefault:""`
	Region     string `env:"DB_REGION" envDefault:""`
	DBName     string `env:"DB_NAME" envDefault:""`
}

type Config struct {
	ServerPort         string `env:"PORT" envDefault:"9093"`
	Database           DatabaseConfig
	Mode               string `env:"MODE" envDefault:"local"`
	AuthorizedAgentUrl string `env:"AUTHORIZED_AGENT_URL" envDefault:""`
}

func main() {
	var c Config
	err := env.Parse(&c)
	if err != nil {
		slog.Error("failed to parse default config", "error", err)
		os.Exit(1)
	}

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}

	h := &logger.ContextHandler{Handler: slog.NewJSONHandler(os.Stdout, opts)}
	slogLogger := slog.New(h)
	slog.SetDefault(slogLogger)
	slog.Info("starting grpc", "port", c.ServerPort)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", c.ServerPort))
	if err != nil {
		slog.Error("failed to start listener for grpc", "error", err)
		os.Exit(1)
	}

	db, err := postgres.NewDBClient(c.Database.ConnString)
	if err != nil {
		slog.Error("failed to connect to db", "error", err)
		os.Exit(1)
	}

	grpcOpts := []grpc.ServerOption{
		grpc.UnaryInterceptor(logger.ContextPropagationUnaryServerInterceptor()),
	}
	
	// Create a gRPC server with an interceptor that uses the logger
	s := grpc.NewServer(grpcOpts...)
	
	// Initialize repositories
	t := repository.NewTransactionRepository(db, "txn_", "le_")
	a := repository.NewAccountRepository(db, "acct_")
	u := repository.NewUserRepository(db, a, "usr_")

	// Register your service
	pb.RegisterApiServiceServer(s, &service.GrpcService{UserRepo: u, AccountRepo: a, TransactionRepo: t})

	// Create and register the health server
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(s, healthServer)

	// Set the health status
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	slog.Info("gRPC server is running with health check enabled", "port", c.ServerPort)

	if err := s.Serve(listener); err != nil {
		slog.Error("failed to serve", "error", err)
		os.Exit(1)
	}
}