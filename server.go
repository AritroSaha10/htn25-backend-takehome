package main

import (
	"github.com/AritroSaha10/htn25-backend-takehome/controller"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

var (
	serv *Server
)

type Server struct {
	Router      *chi.Mux
	DB          *gorm.DB
	Port        string
	Environment string
}

func CreateNewServer(db *gorm.DB, port string, environment string) *Server {
	return &Server{
		Router:      chi.NewRouter(),
		DB:          db,
		Port:        port,
		Environment: environment,
	}
}

func (s *Server) MountHandlers() {
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)

	s.Router.Mount("/users", controller.UserController{}.Routes())
	s.Router.Mount("/scans", controller.ScanController{}.Routes())
}
