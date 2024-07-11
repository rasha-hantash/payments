package gateway

import (
	"github.com/gorilla/mux"
	"myproject/internal/grpc/client"
	"myproject/internal/handlers"
)

func NewRouter(apiClient *client.APIClient) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/users", handlers.CreateUser(apiClient)).Methods("POST")
	return r
}
