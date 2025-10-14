package group

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
	mux.HandleFunc("POST /groups", h.createGroup)
	mux.HandleFunc("DELETE /groups/{id}", h.deleteOwnGroup)
	mux.HandleFunc("GET /groups/{id}", h.viewGroup)
	mux.HandleFunc("GET /groups", h.getCreatedGroups)
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
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	groupName := r.FormValue("name")
	groupDescription := r.FormValue("description")

	if groupName == "" {
		log.Println("ERROR: missing group name")
		http.Error(w, "group name required", http.StatusBadRequest)
		return
	}

	if len(groupName) > 255 {
		log.Println("ERROR: group name too long")
		http.Error(w, "group name must be 255 characters or less", http.StatusBadRequest)
		return
	}

	if len(groupDescription) > 1000 {
		log.Println("ERROR: group description too long")
		http.Error(w, "group description must be 1000 characters or less", http.StatusBadRequest)
		return
	}

	_, err := h.svc.CreateAndJoinGroup(ctx, user.ID, groupName, groupDescription)
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error creating group", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/groups", http.StatusSeeOther)
}

func (h *Handler) getCreatedGroups(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(time.Second*3))
	defer cancel()

	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		log.Println("ERROR:", "no user in context")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	groups, err := h.svc.GetCreatedGroups(ctx, user.ID)
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error fetching groups", http.StatusInternalServerError)
		return
	}

	data := struct {
		User   *model.User
		Groups []*model.Group
	}{
		User:   user,
		Groups: groups,
	}
	if err := h.p.Present(w, r, "groups", data); err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error loading groups page", http.StatusInternalServerError)
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

	log.Println("Group", groupDetails)

	data := struct {
		User  *model.User
		Group *model.Group
	}{
		User:  user,
		Group: &model.Group{Name: "Some Group", Description: "Some desc", CreatedAt: time.Now()},
	}

	if err := h.p.Present(w, r, "group", data); err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error loading groups page", http.StatusInternalServerError)
	}
}
