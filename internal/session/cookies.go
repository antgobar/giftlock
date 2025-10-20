package session

import (
	"giftlock/internal/security"
	"net/http"
	"time"
)

func GetCookieValue(request *http.Request) (string, error) {
	cookie, err := request.Cookie(security.SessionName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func SetCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     security.SessionName,
		Value:    token,
		Expires:  security.ExpireInOneYear(),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}

func ClearCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     security.SessionName,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}
