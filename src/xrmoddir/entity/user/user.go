package user

import (
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
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

func (u *User) SetPassword(password string) error {
	enc, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Error generating password hash: %v", err)
	}
	u.password = enc
	return nil
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
