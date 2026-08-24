package cart

import (
	"encoding/json"
	"net/http"

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

func (h *Handler) AddCartItem(w http.ResponseWriter, r *http.Request) {
	var req ReqAddCartItem

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)

	if err != nil {
		util.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// For now use hard coded cart id
	cartId := 1
	// userId := 6

	createdItem, err := h.svc.AddItem(domain.CartItem{
		CartID:    int64(cartId),
		Quantity:  req.Quantity,
		ProductID: req.ProductID,
	})

	util.SendData(w, http.StatusCreated, createdItem)

}
