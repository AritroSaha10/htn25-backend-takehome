package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserController struct{}

func (c UserController) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", c.GetAllUsers)
	return r
}

func (c UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Users!"))
}
