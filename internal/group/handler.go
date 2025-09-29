package group

import (
	"context"
	"encoding/json"
	"giftlock/internal/auth"
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
	mux.HandleFunc("POST /api/groups", h.createGroup)
	mux.HandleFunc("GET /api/groups", h.getCreatedGroups)
}

func (h *Handler) createGroup(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(time.Second*3))
	defer cancel()

	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		log.Println("ERROR:", "no user in context")
		http.Error(w, "Error getting user", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Println("ERROR: unable to parse form:", err)
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	groupName := r.FormValue("name")
	if groupName == "" {
		log.Println("ERROR: missing group name")
		http.Error(w, "group name required", http.StatusBadRequest)
		return
	}
	groupDescription := r.FormValue("description")

	group, err := h.svc.CreateAndJoinGroup(ctx, user.ID, groupName, groupDescription)
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error creating group", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(group); err != nil {
		log.Println("ERROR: unable to encode response:", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getCreatedGroups(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(time.Second*3))
	defer cancel()

	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		log.Println("ERROR:", "no user in context")
		http.Error(w, "Error getting user", http.StatusUnauthorized)
		return
	}

	groups, err := h.svc.GetCreatedGroups(ctx, user.ID)
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error fetching groups", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(groups); err != nil {
		log.Println("ERROR: unable to encode response:", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
