package main

import (
	"github.com/AritroSaha10/htn25-backend-takehome/controller"
	_ "github.com/AritroSaha10/htn25-backend-takehome/docs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"gorm.io/gorm"
)

var (
	serv *Server
)

// @title           HTN25 Backend API
// @version         0.1
// @description     Backend API for Hack the North 2025 Backend Challenge

// @host      localhost:8080
// @BasePath  /
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

	// Mount Swagger endpoint only in development mode
	if s.Environment == "development" {
		s.Router.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		))
	}

	s.Router.Mount("/users", controller.UserController{}.Routes())
	s.Router.Mount("/scans", controller.ScanController{}.Routes())
}
