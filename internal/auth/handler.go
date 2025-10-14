package auth

import (
	"context"
	"giftlock/internal/presentation"
	"giftlock/internal/session"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	svc *Service
	p   presentation.Presenter
}

func NewHandler(svc *Service, p presentation.Presenter) *Handler {
	return &Handler{svc: svc, p: p}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /login", h.logIn)
	mux.HandleFunc("GET /logout", h.logOut)
	mux.HandleFunc("POST /logout", h.logOut)
}

func (h *Handler) logIn(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	sesh, err := h.svc.LogIn(ctx, username, password)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "invalid login", http.StatusUnauthorized)
		return
	}

	session.SetCookie(w, string(sesh.Token))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) logOut(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(time.Second*3))
	defer cancel()

	session.ClearCookie(w)

	user, ok := UserFromContext(r.Context())
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err := h.svc.LogOut(ctx, user.ID)
	if err != nil {
		log.Println("Error logging out:", err.Error())
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)

}
