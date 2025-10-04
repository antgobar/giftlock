package group

import (
	"context"
	"encoding/json"
	"giftlock/internal/auth"
	"giftlock/internal/model"
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
	mux.HandleFunc("DELETE /api/groups/{id}", h.deleteOwnGroup)
	mux.HandleFunc("GET /api/groups/{id}", h.viewGroup)
}

type groupCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
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

	var req groupCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("ERROR: unable to parse JSON:", err)
		http.Error(w, "invalid JSON data", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		log.Println("ERROR: missing group name")
		http.Error(w, "group name required", http.StatusBadRequest)
		return
	}

	if len(req.Name) > 255 {
		log.Println("ERROR: group name too long")
		http.Error(w, "group name must be 255 characters or less", http.StatusBadRequest)
		return
	}

	if len(req.Description) > 1000 {
		log.Println("ERROR: group description too long")
		http.Error(w, "group description must be 1000 characters or less", http.StatusBadRequest)
		return
	}

	group, err := h.svc.CreateAndJoinGroup(ctx, user.ID, req.Name, req.Description)
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

func (h *Handler) deleteOwnGroup(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(time.Second*3))
	defer cancel()

	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		log.Println("ERROR:", "no user in context")
		http.Error(w, "Error getting user", http.StatusUnauthorized)
		return
	}

	groupId, err := model.IdFromString[model.GroupId](r.PathValue("id"))
	if err != nil {
		log.Println("ERROR: invalid user id", err)
		http.Error(w, "Invalid group id", http.StatusBadRequest)
		return
	}

	if err := h.svc.DeleteGroup(ctx, user.ID, groupId); err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error deleting group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) viewGroup(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(time.Second*3))
	defer cancel()

	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		log.Println("ERROR:", "no user in context")
		http.Error(w, "Error getting user", http.StatusUnauthorized)
		return
	}

	groupId, err := model.IdFromString[model.GroupId](r.PathValue("id"))
	if err != nil {
		log.Println("ERROR: invalid group id", err)
		http.Error(w, "Invalid group id", http.StatusBadRequest)
		return
	}

	groupDetails, err := h.svc.ViewGroup(ctx, user.ID, groupId)
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error fetching group", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(groupDetails); err != nil {
		log.Println("ERROR: unable to encode response:", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
