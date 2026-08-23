package product

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/shohann/golang-ecommerce-api/apperr"
	"github.com/shohann/golang-ecommerce-api/config"
	"github.com/shohann/golang-ecommerce-api/domain"
	middleware "github.com/shohann/golang-ecommerce-api/rest/middlewares"
	"github.com/shohann/golang-ecommerce-api/util"
)

type Handler struct {
	cnf         *config.Config
	middlewares *middleware.Middlewares
	svc         Service
}

func NewHandler(
	cnf *config.Config,
	middlewares *middleware.Middlewares,
	svc Service,
) *Handler {
	return &Handler{
		cnf:         cnf,
		middlewares: middlewares,
		svc:         svc,
	}
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req ReqCreateProduct

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)

	if err != nil {
		util.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	createdProduct, err := h.svc.Create(domain.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
		IsActive:    req.IsActive,
	})

	util.SendData(w, http.StatusCreated, createdProduct)
}

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
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

	products, total, err := h.svc.GetPublicProducts(page, limit)
	if err != nil {
		util.SendAppError(w, err)
		return
	}

	util.SendPage(w, products, page, limit, total)
}

func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil || id <= 0 {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	err = h.svc.Delete(id)

	if err != nil {
		if apperr.IsKind(err, apperr.KindNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
