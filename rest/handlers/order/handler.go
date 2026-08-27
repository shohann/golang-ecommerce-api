package order

import (
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

func (h *Handler) OrderCheckOut(w http.ResponseWriter, r *http.Request) {
	userID, ok := authUserID(w, r)
	if !ok {
		return
	}

	cartItems, err := h.svc.OrderCheckOut(userID)
	if err != nil {
		util.SendAppError(w, err)
		return
	}

	util.SendData(w, http.StatusCreated, cartItems)
	// util.SendData(w, http.StatusCreated, userID)
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
