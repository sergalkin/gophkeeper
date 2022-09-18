package model

import "github.com/google/uuid"

type User struct {
	ID       *uuid.UUID `json:"id"`
	Login    string     `json:"login" validate:"gte=3"`
	Password string     `json:"-" validate:"gte=3"`
}
