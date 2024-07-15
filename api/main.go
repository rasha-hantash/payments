package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rasha-hantash/chariot-takehome/api/pkgs/logger"
	"google.golang.org/grpc"
	"log"
	"log/slog"
	"net"
	"os"

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

type Client struct {
	Conn *sql.DB
}

func newDBClient(psqlConnStr string) *Client {
	conn, err := sql.Open("postgres", psqlConnStr)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Ping()
	if err != nil {
		log.Println(err.Error())
		log.Fatal(err)
	}
	slog.Info("postgres connection success")
	return &Client{Conn: conn}
}

func main() {
	var c Config

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}

	h := &logger.ContextHandler{Handler: slog.NewJSONHandler(os.Stdout, opts)}
	slogLogger := slog.New(h)
	slog.SetDefault(slogLogger)
	err := env.Parse(&c)
	if err != nil {
		slog.Error("failed to parse default config", "error", err)
		os.Exit(1)
	}
	slog.Info("starting grpc", "port", c.ServerPort)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", c.ServerPort))
	if err != nil {
		slog.Error("failed to start listener for grpc", "error", err)
		os.Exit(1)
	}
	db := newDBClient(c.Database.ConnString)

	grpcOpts := []grpc.ServerOption{
		grpc.UnaryInterceptor(logger.ContextPropagationUnaryServerInterceptor()),
	}
	// Create a gRPC with an interceptor that uses the logger
	s := grpc.NewServer(grpcOpts...)

	t := repository.NewTransactionRepository(db.Conn, "txn_")
	a := repository.NewAccountRepository(db.Conn, "acct_")
	u := repository.NewUserRepository(db.Conn, a, "usr_")

	pb.RegisterApiServiceServer(s, &service.GrpcService{UserRepo: u, AccountRepo: a, TransactionRepo: t})
	if err := s.Serve(listener); err != nil {
		slog.Error("failed to serve", "error", err)
		os.Exit(1)
	}
}
