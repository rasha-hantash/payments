package handlers

import (
	"encoding/json"
	"net/http"
	pb "github.com/rasha-hantash/chariot-takehome/api/grpc/proto"
	client "github.com/rasha-hantash/chariot-takehome/gateway/grpcClient"
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