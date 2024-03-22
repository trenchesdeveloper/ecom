package types

import (
	go_ecom "github.com/trenchesdeveloper/go-ecom"
	"regexp"
	"strings"
)

type UserStore interface {
	CreateUser(user User) error
	GetUserByEmail(email string) (*User, error)
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"createdAt"`
}

type RegisterInput struct {
	FirstName       string `json:"firstName" validate:"required"`
	LastName        string `json:"lastName" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=8`
}

func (r *RegisterInput) Sanitize() {
	r.FirstName = strings.Trim(r.FirstName, " ")
	r.LastName = strings.Trim(r.LastName, " ")
	r.Email = strings.Trim(r.Email, " ")
	r.Email = strings.ToLower(r.Email)
}

func (r *RegisterInput) Validate() error {
	if r.Password != r.ConfirmPassword {
		return go_ecom.ErrPasswordMismatch
	}

	if len(r.Password) < 8 {
		return go_ecom.ErrPasswordTooShort
	}

	// Validate email
	if !validateEmail(r.Email) {
		return go_ecom.ErrInvalidEmail

	}

	return nil
}

func validateEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(regex, email)
	return match
}
