package pages

import (
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
	mux.HandleFunc("/", h.home)
	mux.HandleFunc("GET /register", h.register)
	mux.HandleFunc("GET /login", h.login)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.HandlerFunc(staticResources)))
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.HandlerFunc(frontendAssets)))
	mux.HandleFunc("/favicon.ico", favicon)
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "frontend/dist/index.html")
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	if err := h.presenter.Present(w, r, "login", nil); err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	if err := h.presenter.Present(w, r, "register", nil); err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func staticResources(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	fs := http.FileServer(http.Dir("static"))
	fs.ServeHTTP(w, r)
}

func favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/dist/favicon.ico")
}

func frontendAssets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	fs := http.FileServer(http.Dir("frontend/dist/assets"))
	fs.ServeHTTP(w, r)
}
