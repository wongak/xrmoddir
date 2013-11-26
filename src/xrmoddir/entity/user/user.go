package user

import (
	"code.google.com/p/go.crypto/bcrypt"
	"time"
)

type User struct {
	Id       int64
	Username string
	Created  time.Time

	password []byte
	Email    string
	Active   bool
}

func (u User) ValidatePassword(inputPassword []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(u.password, inputPassword)
	if err == nil {
		return true, nil
	}
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	return false, err
}
