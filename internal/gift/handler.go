package gift

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
	mux.HandleFunc("POST /gifts", h.addGiftToWishList)
	mux.HandleFunc("GET /gifts/me", h.viewOwnGifts)
	mux.HandleFunc("GET /gifts/user/{id}", h.viewUserGifts)
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

	if err := r.ParseForm(); err != nil {
		log.Println("ERROR: unable to parse form:", err)
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	giftTitle := r.FormValue("title")
	if giftTitle == "" {
		log.Println("ERROR: missing gift title")
		http.Error(w, "gift title required", http.StatusBadRequest)
		return
	}
	giftDescription := r.FormValue("description")
	giftLink := r.FormValue("link")

	gift, err := h.svc.CreateOwnGift(ctx, user.ID, giftTitle, giftDescription, giftLink)
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error creating gift", http.StatusInternalServerError)
		return
	}

	data := struct {
		Gift *model.Gift
	}{
		Gift: gift,
	}
	if err := h.presenter.Present(w, r, "gift", data); err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "resource error", http.StatusInternalServerError)
	}
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
	data := struct {
		User    *model.User
		Devices []*model.Gift
	}{
		User:    user,
		Devices: gifts,
	}
	if err := h.presenter.Present(w, r, "gifts", data); err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "resource error", http.StatusInternalServerError)
	}
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
	data := struct {
		User    *model.User
		Devices []*model.Gift
	}{
		User:    nil,
		Devices: gifts,
	}
	if err := h.presenter.Present(w, r, "gifts", data); err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "resource error", http.StatusInternalServerError)
	}

}
