package assets

import (
	"net/http"
	"strings"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.home)
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.HandlerFunc(frontendAssets)))
	mux.HandleFunc("/favicon.ico", favicon)
	mux.HandleFunc("/bulma.min.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/dist/bulma.min.css")
	})
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" || (!strings.HasPrefix(r.URL.Path, "/api/") && !strings.HasPrefix(r.URL.Path, "/assets/") && r.URL.Path != "/favicon.ico") {
		http.ServeFile(w, r, "frontend/dist/index.html")
		return
	}
	http.NotFound(w, r)
}

func favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/dist/favicon.ico")
}

func frontendAssets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	fs := http.FileServer(http.Dir("frontend/dist/assets"))
	fs.ServeHTTP(w, r)
}
