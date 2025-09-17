package user

import (
	"context"
	"giftlock/internal/model"
	"giftlock/internal/presentation"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	svc       *Service
	presenter presentation.Presenter
}

func NewHandler(svc *Service, p presentation.Presenter) *Handler {
	if svc == nil {
		log.Fatalln("User handler is nil")
	}
	return &Handler{svc: svc, presenter: p}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /register", h.register)
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	username := r.FormValue("username")
	password := r.FormValue("password")
	err := h.svc.Register(ctx, username, password)

	if err == ErrUsernameTaken {
		data := struct {
			Error string
			User  *model.User
		}{
			Error: "Username is already taken",
			User:  nil,
		}

		if err := h.presenter.Present(w, r, "register", data); err != nil {
			log.Println("ERROR:", err.Error())
			http.Error(w, "template error", http.StatusInternalServerError)
			return
		}
		return
	}

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "error registering user", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
