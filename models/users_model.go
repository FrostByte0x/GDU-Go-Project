package models

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Create a go-compatible enum for the roles of the restaurant
//
// Roles are checked by the Authorization Middleware
type Role string

const (
	Administrator Role = "administrator"
	Preparator    Role = "preparator"
	Reception     Role = "reception"
)

// Users perform data operations on the resources of the restaurant
type User struct {
	ID        uuid.UUID `gorm:"primaryKey;type:char(36)"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
	UpdatedAt time.Time
	Username  string `gorm:"size:32" json:"username"`
	Password  string `gorm:"size:255" json:"password"`
	Role      Role   `gorm:"type:enum('administrator','preparator','reception')" json:"role"`
}

// Testing Gorm hooks
//
// https://gorm.io/docs/hooks.html
//
// BeforeCreate will be called before every User Creation
func (u *User) BeforeCreate(tx *gorm.DB) error {
	UserUUID, err := uuid.NewV7()
	if err != nil {
		return err
	}
	u.ID = UserUUID
	return nil
}

// After a user is created, log a message and the corresponding UUID
func (u *User) AfterCreate(tx *gorm.DB) error {
	slog.Info("New user created, UUID automatically created", "UUID", u.ID)
	return nil
}

// Work in progress - struct to return to administrator when changing a user
type UserReturn struct {
	ID        uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `gorm:"size:32"`
	Role      Role      `json:"role"`
}
