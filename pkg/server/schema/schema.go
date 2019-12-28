package schema

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/gofrs/uuid"
)

// JSONTime wrapper to correctly structure time stamps on the request
type JSONTime struct {
	time.Time
}

// MarshalJSON adhere to interface
func (jt *JSONTime) MarshalJSON() ([]byte, error) {
	t := fmt.Sprintf("\"%s\"", jt.Format("2006-01-02T15:04:05Z"))
	return []byte(t), nil
}

// TimeToJSONTime converts time struct to JSONTime
func TimeToJSONTime(t time.Time) JSONTime {
	return JSONTime{t}
}

// RegisterUserRequest request to create a new user.
type RegisterUserRequest struct {
	Email           string  `json:"email,omitempty"`
	Password        string  `json:"password,omitempty"`
	Name            string  `json:"name"`
	Gender          *string `json:"gender,omitempty"`
	ProfileImageURL *string `json:"profile_image_url,omitempty"`
}

// Bind conform to Binder interface
func (r *RegisterUserRequest) Bind(req *http.Request) error {
	return nil
}

// RegisterUserResponse response with newly registered user
type RegisterUserResponse struct {
	User *User `json:"user"`
}

// Render conforms to Renderer interface
func (r *RegisterUserResponse) Render(w http.ResponseWriter, req *http.Request) error {
	render.Status(req, http.StatusCreated)
	return nil
}

// User is the public struct representing a user model
type User struct {
	ID              uuid.UUID `json:"id"`
	Email           string    `json:"email"`
	Name            string    `json:"name"`
	Gender          *string   `json:"gender,omitempty"`
	ProfileImageURL *string   `json:"profile_image_url,omitempty"`
	Status          string    `json:"status"`
	CreatedAt       JSONTime  `json:"created_at,omitempty"`
	UpdatedAt       JSONTime  `json:"updated_at,omitempty"`
}

// GetUserResponse response with newly registered user
type GetUserResponse struct {
	User *User `json:"user"`
}

// Render conforms to Renderer interface
func (r *GetUserResponse) Render(w http.ResponseWriter, req *http.Request) error {
	render.Status(req, http.StatusOK)
	return nil
}
