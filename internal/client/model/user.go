package model

import (
	"github.com/google/uuid"
)

type User struct {
	ID       *uuid.UUID
	Login    string
	Password string `json:"-"`
}
