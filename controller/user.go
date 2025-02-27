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

// @Summary Get all users
// @Description Get a list of all users in the system
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} model.User
// @Failure 500 {object} util.ErrResponse
// @Router /users [get]
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

// @Summary Get user by ID
// @Description Get a single user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.User
// @Failure 400 {object} util.ErrResponse
// @Failure 404 {object} util.ErrResponse
// @Failure 500 {object} util.ErrResponse
// @Router /users/{id} [get]
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
		render.Render(w, r, util.ErrNotFoundRender("no user with given id exists", err))
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

// @Summary Update user by ID
// @Description Update a user's information by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body model.UserUpdate true "User update information"
// @Success 200 {object} model.User
// @Failure 400 {object} util.ErrResponse
// @Failure 404 {object} util.ErrResponse
// @Failure 500 {object} util.ErrResponse
// @Router /users/{id} [put]
func (c UserController) UpdateUserByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not parse id")
		render.Render(w, r, util.ErrBadRequestRender("id is not unsigned int", err))
		return
	}

	// Parse request body into a UserUpdate struct
	userUpdate := model.UserUpdate{}
	if err := render.Bind(r, &userUpdate); err != nil {
		render.Render(w, r, util.ErrBadRequestRender("invalid request body", err))
		return
	}

	// Update the user and handle all error cases
	user, err := model.UpdateUserByID(lib.GetDB(), uint(idUint), userUpdate)
	if errors.Is(err, util.ErrNotFound) {
		render.Render(w, r, util.ErrNotFoundRender("no user with given id exists", err))
		return
	} else if errors.Is(err, util.ErrBadRequest) {
		render.Render(w, r, util.ErrBadRequestRender(err.Error(), err))
		return
	} else if err != nil {
		render.Render(w, r, util.ErrServerRender(err))
		return
	}

	if err := render.Render(w, r, &user); err != nil {
		render.Render(w, r, util.ErrRender(err))
		return
	}
}
