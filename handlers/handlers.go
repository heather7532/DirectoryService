package handlers

import (
	"DirectoryService/db"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"

	"DirectoryService/models"
)

// Handler function for registering a service
func (s *Server) RegisterServiceHandler(w http.ResponseWriter, r *http.Request) {
	var service models.Service
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	dbCtx := db.NewDbCtx(s.DB)

	ctx := context.Background()
	newService, err := dbCtx.RegisterService(ctx, service)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newService); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// Handler to register and instance of a service
func (s *Server) RegisterServiceInstanceHandler(w http.ResponseWriter, r *http.Request) {
	var serviceInstance models.ServiceInstance
	if err := json.NewDecoder(r.Body).Decode(&serviceInstance); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	dbCtx := db.NewDbCtx(s.DB)

	ctx := context.Background()
	newInstance, err := dbCtx.CreateServiceInstance(ctx, serviceInstance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newInstance); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// Handler to remove a service instance
func (s *Server) RemoveServiceInstanceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instanceID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid instance ID", http.StatusBadRequest)
		return
	}

	dbCtx := db.NewDbCtx(s.DB)

	ctx := context.Background()
	err = dbCtx.RemoveServiceInstance(ctx, instanceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
