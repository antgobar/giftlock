package user

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	if svc == nil {
		log.Fatalln("User handler is nil")
	}
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /register", h.register)
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(time.Second*3))
	defer cancel()

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	err := h.svc.Register(ctx, username, password)

	if errors.Is(err, ErrUsernameTaken) {
		http.Error(w, "username taken", http.StatusConflict)
		return
	}

	if err != nil {
		log.Printf("ERROR: registering user '%s': %v", username, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
