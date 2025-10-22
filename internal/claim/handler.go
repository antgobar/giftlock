package claim

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
	mux.HandleFunc("GET /claims", h.viewOwnClaims)
}

func (h *Handler) viewOwnClaims(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(time.Second*3))
	defer cancel()

	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		log.Println("ERROR:", "no user in context")
		http.Error(w, "Error getting user", http.StatusUnauthorized)
		return
	}

	groupGifts, err := h.svc.ViewUserClaims(ctx, user.ID)
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error retrieving claims", http.StatusInternalServerError)
		return
	}
	data := struct {
		User  *model.User
		Gifts []*model.GroupGift
	}{
		User:  user,
		Gifts: groupGifts,
	}

	h.p.Present(w, r, "user_claims", data)
}
