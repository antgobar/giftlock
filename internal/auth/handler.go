package auth

import (
	"context"
	"encoding/json"
	"giftlock/internal/session"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/login", h.logIn)
	mux.HandleFunc("POST /api/logout", h.logOut)
	mux.HandleFunc("GET /api/logout", h.logOut)
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) logIn(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	var req loginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	sesh, err := h.svc.LogIn(ctx, req.Username, req.Password)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "invalid login", http.StatusUnauthorized)
		return
	}

	session.SetCookie(w, string(sesh.Token))
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) logOut(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	user, ok := UserFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusOK)
		return
	}

	err := h.svc.LogOut(ctx, user.ID)
	if err != nil {
		log.Println("Error logging out:", err.Error())
	}
	w.WriteHeader(http.StatusOK)
}
