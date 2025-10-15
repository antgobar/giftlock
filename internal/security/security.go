package security

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const SessionName string = "giftlock_app_session"

func GenerateUUID() string {
	return uuid.New().String()
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ExpireInOneYear() time.Time {
	return time.Now().UTC().Add(time.Hour * 24 * 365)

}
