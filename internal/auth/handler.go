package auth

import (
	"context"
	"giftlock/internal/model"
	"giftlock/internal/presentation"
	"giftlock/internal/session"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	svc       *Service
	presenter presentation.Presenter
}

func NewHandler(svc *Service, p presentation.Presenter) *Handler {
	return &Handler{svc: svc, presenter: p}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /login", h.logIn)
	mux.HandleFunc("/logout", h.logOut)
}

func (h *Handler) logIn(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	username := r.FormValue("username")
	password := r.FormValue("password")

	sesh, err := h.svc.LogIn(ctx, username, password)
	if err != nil {
		log.Println(err.Error())
		data := struct {
			Error string
			User  *model.User
		}{
			Error: "Incorrect username or password",
			User:  nil,
		}

		if err := h.presenter.Present(w, r, "login", data); err != nil {
			log.Println("ERROR:", err.Error())
			http.Error(w, "template error", http.StatusInternalServerError)
			return
		}
		return
	}

	session.SetCookie(w, string(sesh.Token))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) logOut(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	user, ok := UserFromContext(r.Context())
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := h.svc.LogOut(ctx, user.ID)
	if err != nil {
		log.Println("ERROR:", err.Error())
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
