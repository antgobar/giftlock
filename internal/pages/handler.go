package pages

import (
	"giftlock/internal/auth"
	"giftlock/internal/model"
	"giftlock/internal/presentation"
	"log"
	"net/http"
)

type Handler struct {
	presenter presentation.Presenter
}

func NewHandler(p presentation.Presenter) *Handler {
	return &Handler{presenter: p}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.home)
	mux.HandleFunc("GET /register", h.register)
	mux.HandleFunc("GET /login", h.login)
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	user, _ := auth.UserFromContext(r.Context())

	data := struct {
		User *model.User
	}{
		User: user,
	}

	if err := h.presenter.Present(w, r, "home", data); err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error loading home page", http.StatusInternalServerError)
	}
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	if err := h.presenter.Present(w, r, "login", nil); err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error loading login page", http.StatusInternalServerError)
	}
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	if err := h.presenter.Present(w, r, "register", nil); err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error loading registration page", http.StatusInternalServerError)
	}
}
