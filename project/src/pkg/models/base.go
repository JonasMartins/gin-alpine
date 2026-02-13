// Package models ...
package models

import (
	"time"

	"github.com/google/uuid"
)

type Base struct {
	ID        int        `json:"id"`
	UUID      uuid.UUID  `json:"uuid,omitempty"`
	Version   int        `json:"version,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"-"`
}

type BaseFeilds struct {
	ID        int        `json:"id"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"-"`
}

type BaseWithoutID struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"-"`
}
