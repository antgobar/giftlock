package gift

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
	mux.HandleFunc("POST /groups/{id}/gifts", h.addGiftToWishList)
	mux.HandleFunc("GET /user/me/gifts", h.viewOwnGifts)
	mux.HandleFunc("GET /user/{id}/gifts", h.viewUserGifts)
}

func (h *Handler) addGiftToWishList(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	giftTitle := r.FormValue("title")
	giftDescription := r.FormValue("description")
	giftLink := r.FormValue("link")

	gift, err := h.svc.CreateOwnGift(ctx, user.ID, groupId, giftTitle, giftDescription, giftLink)
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error creating gift", http.StatusInternalServerError)
		return
	}

	log.Println("Gift", gift, "created by", user.ID)

	http.Redirect(w, r, "/groups/"+r.PathValue("id"), http.StatusSeeOther)
}

func (h *Handler) viewOwnGifts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(time.Second*3))
	defer cancel()

	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		log.Println("ERROR:", "no user in context")
		http.Error(w, "Error getting user", http.StatusUnauthorized)
		return
	}

	gifts, err := h.svc.ViewUserGifts(ctx, user.ID)
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error retrieving gifts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gifts)
}

func (h *Handler) viewUserGifts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(time.Second*3))
	defer cancel()

	userId, err := model.IdFromString[model.UserId](r.PathValue("id"))
	if err != nil {
		log.Println("ERROR: invalid user id", err)
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	gifts, err := h.svc.ViewUserGifts(ctx, userId)
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error retrieving gifts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gifts)

}
