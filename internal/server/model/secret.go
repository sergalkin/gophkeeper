package model

import (
	"time"

	"github.com/google/uuid"
)

type Secret struct {
	ID        int        `json:"ID"`
	UserID    uuid.UUID  `json:"user_id"`
	TypeID    int        `json:"type_id"`
	Title     string     `json:"title"`
	Content   []byte     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
