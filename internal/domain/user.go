package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" validate:"required,min=5"`
	Mobile    string    `json:"mobile" validate:"required,min=10,max=11"`
	Password  string    `json:"password" validate:"required,min=6"`
	Role      string    `json:"role" validate:"required"`
	Timestamp time.Time `json:"timestamp"`
}

type LoginData struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Name  string `json:"name"`
	Token string `json:"token"`
	Role  string `json:"role"`
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
