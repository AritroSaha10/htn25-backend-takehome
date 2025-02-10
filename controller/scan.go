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
	r.Put("/{badge_code}", c.ScanUser)
	return r
}

// @Summary Get aggregate scan statistics
// @Description Get aggregated scan statistics with optional filtering
// @Tags scans
// @Accept json
// @Produce json
// @Param min_frequency query int false "Minimum frequency filter"
// @Param max_frequency query int false "Maximum frequency filter"
// @Param activity_category query string false "Activity category filter"
// @Success 200 {array} model.ScanAggregate
// @Failure 500 {object} util.ErrResponse
// @Router /scans [get]
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

// @Summary Record a new scan for a user
// @Description Record a new activity scan for a user using their badge code
// @Tags scans
// @Accept json
// @Produce json
// @Param badge_code path string true "User's badge code"
// @Param scan body model.Scan true "Scan information"
// @Success 200 {object} model.Scan
// @Failure 400 {object} util.ErrResponse
// @Failure 404 {object} util.ErrResponse
// @Failure 500 {object} util.ErrResponse
// @Router /scans/{badge_code} [put]
func (c ScanController) ScanUser(w http.ResponseWriter, r *http.Request) {
	// Get the badge code from the URL parameter
	badgeCode := chi.URLParam(r, "badge_code")
	if badgeCode == "" {
		render.Render(w, r, util.ErrBadRequestRender("badge code is required", nil))
		return
	}

	// Find user with given badge code
	user := model.User{}
	tx := lib.
		GetDB().
		Where("badge_code = ?", badgeCode).
		Limit(1).
		Find(&user)
	if tx.RowsAffected == 0 {
		render.Render(w, r, util.ErrNotFoundRender("no user with given badge code exists", nil))
		return
	}

	// Parse request body into a Scan struct
	scan := model.Scan{
		UserID: user.ID,
	}
	if err := render.Bind(r, &scan); err != nil {
		render.Render(w, r, util.ErrBadRequestRender("invalid request body", err))
		return
	}

	// Add scan to database
	err := model.CreateScan(lib.GetDB(), &scan)
	if err != nil {
		render.Render(w, r, util.ErrServerRender(err))
		return
	}

	if err := render.Render(w, r, &scan); err != nil {
		render.Render(w, r, util.ErrRender(err))
		return
	}
}
