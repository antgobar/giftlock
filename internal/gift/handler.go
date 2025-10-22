package gift

import (
	"context"
	"fmt"
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
	mux.HandleFunc("POST /groups/{id}/gifts", h.addGiftToWishList)
	mux.HandleFunc("GET /gifts", h.viewOwnGifts)
	mux.HandleFunc("GET /groups/{groupId}/user/{memberId}/gifts", h.viewGroupMemberGifts)
	mux.HandleFunc("DELETE /user/me/gifts/{id}", h.deleteOwnGift)
	mux.HandleFunc("POST /gifts/{id}/claim", h.claim)
	mux.HandleFunc("POST /gift/{id}/unclaim", h.unclaim)
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

	groupGifts, err := h.svc.ViewUserGifts(ctx, user.ID)
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error retrieving gifts", http.StatusInternalServerError)
		return
	}
	data := struct {
		User  *model.User
		Gifts []*model.GroupGift
	}{
		User:  user,
		Gifts: groupGifts,
	}

	h.p.Present(w, r, "user_gifts", data)
}

func (h *Handler) viewGroupMemberGifts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(time.Second*3))
	defer cancel()

	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		log.Println("ERROR:", "no user in context")
		http.Error(w, "Error getting user", http.StatusUnauthorized)
		return
	}

	groupId, err := model.IdFromString[model.GroupId](r.PathValue("groupId"))
	if err != nil {
		log.Println("ERROR: invalid group id", err)
		http.Error(w, "Invalid group id", http.StatusBadRequest)
		return
	}

	memberId, err := model.IdFromString[model.UserId](r.PathValue("memberId"))
	if err != nil {
		log.Println("ERROR: invalid user id", err)
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	gifts, err := h.svc.ViewGroupUserGifts(ctx, groupId, memberId)
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error retrieving gifts", http.StatusInternalServerError)
		return
	}

	for _, gift := range gifts {
		fmt.Println("GIFT:", gift.Title, gift.CreatedBy, gift.ClaimedBy)
	}

	data := struct {
		CurrentUserId model.UserId
		Gifts         []*model.Gift
	}{
		CurrentUserId: user.ID,
		Gifts:         gifts,
	}

	h.p.Present(w, r, "group_user_gifts", data)

}

func (h *Handler) deleteOwnGift(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(time.Second*3))
	defer cancel()

	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		log.Println("ERROR:", "no user in context")
		http.Error(w, "Error getting user", http.StatusUnauthorized)
		return
	}

	giftId, err := model.IdFromString[model.GiftId](r.PathValue("id"))
	if err != nil {
		log.Println("ERROR: invalid gift id", err)
		http.Error(w, "Invalid gift id", http.StatusBadRequest)
		return
	}

	err = h.svc.DeleteGift(ctx, user.ID, giftId)
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error deleting gift", http.StatusInternalServerError)
		return
	}

	log.Println("Gift", giftId, "deleted by", user.ID)

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) claim(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(time.Second*3))
	defer cancel()

	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		log.Println("ERROR:", "no user in context")
		http.Error(w, "Error getting user", http.StatusUnauthorized)
		return
	}

	giftId, err := model.IdFromString[model.GiftId](r.PathValue("id"))
	if err != nil {
		log.Println("ERROR: invalid gift id", err)
		http.Error(w, "Invalid gift id", http.StatusBadRequest)
		return
	}

	err = h.svc.ClaimGift(ctx, user.ID, giftId)
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error claiming gift", http.StatusInternalServerError)
		return
	}

	data := struct {
		GiftId model.GiftId
	}{
		GiftId: giftId,
	}

	h.p.Present(w, r, "claim_result", data)
}

func (h *Handler) unclaim(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(time.Second*3))
	defer cancel()

	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		log.Println("ERROR:", "no user in context")
		http.Error(w, "Error getting user", http.StatusUnauthorized)
		return
	}

	giftId, err := model.IdFromString[model.GiftId](r.PathValue("id"))
	if err != nil {
		log.Println("ERROR: invalid gift id", err)
		http.Error(w, "Invalid gift id", http.StatusBadRequest)
		return
	}

	err = h.svc.UnclaimGift(ctx, user.ID, giftId)
	if err != nil {
		log.Println("ERROR:", err.Error())
		http.Error(w, "Error unclaiming gift", http.StatusInternalServerError)
		return
	}

	data := struct {
		GiftId model.GiftId
	}{
		GiftId: giftId,
	}

	h.p.Present(w, r, "unclaim_result", data)
}
