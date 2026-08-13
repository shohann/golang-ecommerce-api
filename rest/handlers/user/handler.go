package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/shohann/golang-ecommerce-api/config"
	"github.com/shohann/golang-ecommerce-api/domain"
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

type ReqCreateUser struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ReqqLoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req ReqCreateUser

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	createdUser, err := h.svc.Create(domain.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		util.SendAppError(w, err)
		return
	}

	util.SendData(w, http.StatusCreated, createdUser)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req ReqqLoginUser

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode((&req))

	if err != nil {
		util.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	loginResponse, err := h.svc.Login(
		req.Email,
		req.Password,
	)

	if err != nil {
		util.SendAppError(w, err)
		return
	}

	util.SendData(w, http.StatusOK, loginResponse)
}

func (h *Handler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.svc.GetProfile(id)
	if err != nil {
		util.SendAppError(w, err)
		return
	}

	util.SendData(w, http.StatusOK, user)
}
