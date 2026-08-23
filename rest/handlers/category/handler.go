package category

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/shohann/golang-ecommerce-api/config"
	middleware "github.com/shohann/golang-ecommerce-api/rest/middlewares"
	"github.com/shohann/golang-ecommerce-api/util"
)

type Handler struct {
	cnf         *config.Config
	svc         Service
	middlewares *middleware.Middlewares
}

func NewHandler(
	cnf *config.Config,
	svc Service,
	middlewares *middleware.Middlewares,
) *Handler {
	return &Handler{
		cnf:         cnf,
		svc:         svc,
		middlewares: middlewares,
	}
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req ReqCreateCategory

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cat, err := h.svc.Create(req.Name)
	if err != nil {
		util.SendAppError(w, err)
		return
	}

	util.SendData(w, http.StatusCreated, ToCategoryResponse(cat))
}

func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	page := int64(1)
	limit := int64(10)

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		parsed, err := strconv.ParseInt(pageStr, 10, 64)
		if err != nil || parsed < 1 {
			util.SendError(w, http.StatusBadRequest, "Invalid page")
			return
		}
		page = parsed
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		parsed, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil || parsed < 1 {
			util.SendError(w, http.StatusBadRequest, "Invalid limit")
			return
		}
		limit = parsed
	}

	if limit > 100 {
		limit = 100
	}

	categories, total, err := h.svc.List(page, limit)
	if err != nil {
		util.SendAppError(w, err)
		return
	}

	util.SendPage(w, ToCategoryResponses(categories), page, limit, total)
}

func (h *Handler) GetCategory(w http.ResponseWriter, r *http.Request) {
	id, ok := parseCategoryID(w, r)
	if !ok {
		return
	}

	cat, err := h.svc.Get(id)
	if err != nil {
		util.SendAppError(w, err)
		return
	}

	util.SendData(w, http.StatusOK, ToCategoryResponse(cat))
}

func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, ok := parseCategoryID(w, r)
	if !ok {
		return
	}

	var req ReqUpdateCategory
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cat, err := h.svc.Update(id, req.Name)
	if err != nil {
		util.SendAppError(w, err)
		return
	}

	util.SendData(w, http.StatusOK, ToCategoryResponse(cat))
}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, ok := parseCategoryID(w, r)
	if !ok {
		return
	}

	if err := h.svc.Delete(id); err != nil {
		util.SendAppError(w, err)
		return
	}

	util.SendData(w, http.StatusOK, "category deleted")
}

func parseCategoryID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		util.SendError(w, http.StatusBadRequest, "Invalid category id")
		return 0, false
	}
	return id, true
}
