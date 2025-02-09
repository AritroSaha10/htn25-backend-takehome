package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/AritroSaha10/htn25-backend-takehome/lib"
	"github.com/AritroSaha10/htn25-backend-takehome/model"
	"github.com/AritroSaha10/htn25-backend-takehome/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

type UserController struct{}

func (c UserController) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", c.GetAllUsers)
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", c.GetUserByID)
		r.Put("/", c.UpdateUserByID)
	})
	return r
}

func (c UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := model.GetUsers(lib.GetDB())
	if err != nil {
		render.Render(w, r, util.ErrServerRender(err))
		return
	}

	// Convert users to a slice of render.Renderer
	renderers := make([]render.Renderer, len(users))
	for i := range users {
		renderers[i] = &users[i]
	}

	if err := render.RenderList(w, r, renderers); err != nil {
		render.Render(w, r, util.ErrRender(err))
		return
	}
}

func (c UserController) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not parse id")
		render.Render(w, r, util.ErrBadRequestRender("id is not unsigned int", err))
		return
	}
	user, err := model.GetUserByID(lib.GetDB(), uint(idUint))

	// Check if user was found, handle other errors appropriately
	if errors.Is(err, util.ErrNotFound) {
		log.Error().Str("id", id).Msg("could not find user")
		render.Render(w, r, util.ErrNotFoundRender(err))
		return
	} else if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not fetch user")
		render.Render(w, r, util.ErrServerRender(err))
		return
	}

	if err := render.Render(w, r, &user); err != nil {
		render.Render(w, r, util.ErrRender(err))
		return
	}
}

func (c UserController) UpdateUserByID(w http.ResponseWriter, r *http.Request) {
	// TODO: Needs to be able to perform partial updates on user info
	// but not scans. Make sure to keep edge cases in mind.
	w.Write([]byte("Hello, User!"))
}
