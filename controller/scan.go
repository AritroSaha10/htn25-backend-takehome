package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ScanController struct{}

func (c ScanController) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", c.GetAllScans)
	return r
}

func (c ScanController) GetAllScans(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Scans!"))
}
