package controller

import (
	"net/http"
	"strconv"

	"github.com/AritroSaha10/htn25-backend-takehome/lib"
	"github.com/AritroSaha10/htn25-backend-takehome/model"
	"github.com/AritroSaha10/htn25-backend-takehome/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type ScanController struct{}

func (c ScanController) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", c.GetAggregateScans)
	r.Put("/{id}", c.ScanUser)
	return r
}

func (c ScanController) GetAggregateScans(w http.ResponseWriter, r *http.Request) {
	minFreqRaw, minErr := strconv.Atoi(r.URL.Query().Get("min_frequency"))
	maxFreqRaw, maxErr := strconv.Atoi(r.URL.Query().Get("max_frequency"))
	activityCategoryRaw := r.URL.Query().Get("activity_category")

	// Change to nil if the query param is not provided
	activityCategory := &activityCategoryRaw
	if activityCategoryRaw == "" {
		activityCategory = nil
	}
	minFreq := &minFreqRaw
	if minErr != nil {
		minFreq = nil
	}
	maxFreq := &maxFreqRaw
	if maxErr != nil {
		maxFreq = nil
	}

	scans, err := model.GetScanAggregates(lib.GetDB(), activityCategory, minFreq, maxFreq)
	if err != nil {
		render.Render(w, r, util.ErrServerRender(err))
		return
	}

	// Convert to list of renderers
	renderers := make([]render.Renderer, len(scans))
	for i, scan := range scans {
		renderers[i] = &scan
	}
	if err := render.RenderList(w, r, renderers); err != nil {
		render.Render(w, r, util.ErrRender(err))
		return
	}
}

func (c ScanController) ScanUser(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the URL parameter
	id := chi.URLParam(r, "id")
	idUint64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		render.Render(w, r, util.ErrBadRequestRender("id is not unsigned int", err))
		return
	}
	idUint := uint(idUint64)

	// Check if the user exists
	if lib.GetDB().Where("id = ?", idUint).Limit(1).Find(&model.User{}).RowsAffected == 0 {
		render.Render(w, r, util.ErrBadRequestRender("user does not exist", nil))
		return
	}

	// Parse request body into a Scan struct
	scan := model.Scan{
		UserID: idUint,
	}
	if err := render.Bind(r, &scan); err != nil {
		render.Render(w, r, util.ErrBadRequestRender("invalid request body", err))
		return
	}

	// Add scan to database
	err = model.CreateScan(lib.GetDB(), &scan)
	if err != nil {
		render.Render(w, r, util.ErrServerRender(err))
		return
	}

	if err := render.Render(w, r, &scan); err != nil {
		render.Render(w, r, util.ErrRender(err))
		return
	}
}
