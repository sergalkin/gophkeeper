package model

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id"`
	Login    string    `json:"login"`
	Password string    `json:"password"`
}

func (u *User) Identity() string {
	return u.ID.String()
}
