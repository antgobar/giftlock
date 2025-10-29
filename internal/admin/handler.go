package admin

import (
	"context"
	"giftlock/internal/auth"
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
	mux.HandleFunc("GET /admin/users", h.listUsers)
	mux.HandleFunc("GET /admin", h.adminPage)
}

func (h *Handler) adminPage(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		log.Println("ERROR:", "no user in context")
		http.Error(w, "Error getting user", http.StatusUnauthorized)
		return
	}
	data := struct {
		User *model.User
	}{
		User: user,
	}
	h.p.Present(w, r, "admin", data)
}

func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(time.Second*3))
	defer cancel()

	userList, err := h.svc.users.List(ctx)
	if err != nil {
		log.Printf("ERROR: listing users %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Users []*model.User
	}{
		Users: userList,
	}
	h.p.Present(w, r, "users", data)
}
