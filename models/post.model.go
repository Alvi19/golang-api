package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Name      string    `gorm:"uniqueIndex;not null" json:"name,omitempty"`
	Lokasi    string    `gorm:"not null" json:"lokasi,omitempty"`
	Image     string    `gorm:"not null" json:"image,omitempty"`
	Status    bool      `gorm:"default:true" json:"status,omitempty"`
	User      uuid.UUID `gorm:"not null" json:"user,omitempty"`
	CreatedAt time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

type CreatePostRequest struct {
	Name      string    `json:"name"  binding:"required"`
	Lokasi    string    `json:"lokasi" binding:"required"`
	Image     string    `json:"image" binding:"required"`
	Status    bool      `json:"status,omitempty"`
	User      string    `json:"user,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type UpdatePost struct {
	Name      string    `json:"name,omitempty"`
	Lokasi    string    `json:"lokasi,omitempty"`
	Image     string    `json:"image,omitempty"`
	Status    bool      `json:"status,omitempty"`
	User      string    `json:"user,omitempty"`
	CreateAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
