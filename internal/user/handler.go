package user

import (
	"context"
	"errors"
	"giftlock/internal/model"
	"giftlock/internal/presentation"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	svc *Service
	p   presentation.Presenter
}

func NewHandler(svc *Service, p presentation.Presenter) *Handler {
	if svc == nil {
		log.Fatalln("User handler is nil")
	}
	return &Handler{svc: svc, p: p}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /register", h.register)
	mux.HandleFunc("POST /users/search/exclude-group/{groupId}", h.searchUserNotInGroup)
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

func (h *Handler) searchUserNotInGroup(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(time.Second*3))
	defer cancel()

	groupId, err := model.IdFromString[model.GroupId](r.PathValue("groupId"))
	if err != nil {
		log.Println("ERROR: invalid group id", err)
		http.Error(w, "Invalid group id", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	usernameSearchTerm := r.FormValue("usernameSearchTerm")

	users, err := h.svc.repo.SearchUserNotInGroup(ctx, groupId, usernameSearchTerm)
	if err != nil {
		log.Println("ERROR:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Users   []*model.User
		GroupID model.GroupId
	}{
		Users:   users,
		GroupID: groupId,
	}
	if err := h.p.Present(w, r, "users", data); err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error loading users page", http.StatusInternalServerError)
	}
}
