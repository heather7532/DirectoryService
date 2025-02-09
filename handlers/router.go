package handlers

import (
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Server struct with DB connection pool
type Server struct {
	DB *pgxpool.Pool
}

// create function to create new server struct
func NewServer(pool *pgxpool.Pool) *Server {
	return &Server{
		pool,
	}
}

// NewRouter sets up the REST routes.
func (s *Server) NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/services", s.RegisterServiceHandler).Methods("POST")
	r.HandleFunc("/service-instances", s.RegisterServiceInstanceHandler).Methods("POST")
	r.HandleFunc("/service-instances/{id}", s.RemoveServiceInstanceHandler).Methods("DELETE")

	return r
}
