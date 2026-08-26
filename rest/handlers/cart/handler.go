package cart

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

func (h *Handler) AddCartItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := authUserID(w, r)
	if !ok {
		return
	}

	var req ReqAddCartItem
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	createdItem, err := h.svc.AddItem(userID, req.ProductID, req.Quantity)
	if err != nil {
		util.SendAppError(w, err)
		return
	}

	util.SendData(w, http.StatusCreated, createdItem)
}

func (h *Handler) GetUserCart(w http.ResponseWriter, r *http.Request) {
	payload, ok := middleware.GetUserPayload(r)

	if !ok || payload.Sub == "" {
		util.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := strconv.Atoi(payload.Sub)
	if err != nil {
		util.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	cartItems, err := h.svc.GetCartByUserId(int64(id))
	if err != nil {
		util.SendAppError(w, err)
		return
	}

	util.SendData(w, http.StatusOK, cartItems)
}

func (h *Handler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := authUserID(w, r)
	if !ok {
		return
	}

	itemID, ok := parseCartItemID(w, r)
	if !ok {
		return
	}

	var req ReqUpdateCartItem
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	item, err := h.svc.UpdateItem(userID, itemID, req.Quantity)
	if err != nil {
		util.SendAppError(w, err)
		return
	}

	util.SendData(w, http.StatusOK, item)
}

func (h *Handler) DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := authUserID(w, r)
	if !ok {
		return
	}

	itemID, ok := parseCartItemID(w, r)
	if !ok {
		return
	}

	if err := h.svc.DeleteItem(userID, itemID); err != nil {
		util.SendAppError(w, err)
		return
	}

	util.SendData(w, http.StatusOK, "cart item deleted")
}

func authUserID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	payload, ok := middleware.GetUserPayload(r)
	if !ok || payload.Sub == "" {
		util.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return 0, false
	}

	id, err := strconv.ParseInt(payload.Sub, 10, 64)
	if err != nil {
		util.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return 0, false
	}

	return id, true
}

func parseCartItemID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		util.SendError(w, http.StatusBadRequest, "Invalid cart item id")
		return 0, false
	}
	return id, true
}
