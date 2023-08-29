package app

import (
	"avito-backend-internship/internal/pkg/db"
	"avito-backend-internship/internal/pkg/service"
	"context"
	"fmt"
	"net/http"
)

type Server struct {
	cfg     *Config
	ctx     context.Context
	db      db.DBops
	service service.Service
}

func NewServer(ctx context.Context, cfg *Config, db db.DBops, service service.Service) *Server {
	return &Server{
		cfg:     cfg,
		ctx:     ctx,
		db:      db,
		service: service,
	}
}

func (s *Server) Run() {
	mux := http.NewServeMux()

	s.initHandlers(mux)

	fmt.Println("Server started on port:", s.cfg.Port)

	err := http.ListenAndServe(":"+s.cfg.Port, mux)

	if err != nil {
		fmt.Println("Failed to start server on port:", s.cfg.Port, "ERROR:", err.Error())
	}
}

func (s *Server) initHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/", s.mainHandler)
	mux.HandleFunc("/create-segment", s.createSegmentHandler)
	mux.HandleFunc("/delete-segment", s.deleteSegmentHandler)
	mux.HandleFunc("/user-in-segment", s.userInSegmentHandler)
	mux.HandleFunc("/get-user-segments", s.getUserSegmentsHandler)
}
