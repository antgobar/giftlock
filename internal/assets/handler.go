package assets

import (
	"net/http"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.HandlerFunc(staticResources)))
	mux.HandleFunc("GET /favicon.ico", favicon)
}
func staticResources(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	fs := http.FileServer(http.Dir("static"))
	fs.ServeHTTP(w, r)
}

func favicon(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/static/favicon.ico", http.StatusMovedPermanently)
}
