package model

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

// User represents a user in the database
type User struct {
	// ID the unique ID for the user
	ID uuid.UUID

	// Name the user's first name
	Name string

	// CreatedAt the date the user was created
	CreatedAt time.Time

	// UpdatedAt the date the user was last updated
	UpdatedAt time.Time

	// DeletedAt the date the user was deleted
	DeletedAt *time.Time
}

// CreateUser creates a new user in the database.
func CreateUser(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}
