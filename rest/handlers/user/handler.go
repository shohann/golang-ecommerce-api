package user

import (
	"encoding/json"
	"net/http"

	"github.com/shohann/golang-ecommerce-api/config"
	"github.com/shohann/golang-ecommerce-api/domain"
	"github.com/shohann/golang-ecommerce-api/util"
)

type Handler struct {
	cnf *config.Config
	svc Service
}

func NewHandler(
	cnf *config.Config,
	svc Service,
) *Handler {
	return &Handler{
		cnf: cnf,
		svc: svc,
	}
}

type ReqCreateUser struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req ReqCreateUser

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)

	if err != nil {
		util.SendError(w, http.StatusBadRequest, "Invalid request body")

	}

	createdUser, err := h.svc.Create(domain.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		util.SendError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	util.SendData(w, http.StatusCreated, createdUser)
}
