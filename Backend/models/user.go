package models

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	UserName string    `json:"username"`
	Email    string    `json:"email" gorm:"unique"`
}

type CreateUserAPIRequest struct {
	UserName string `json:"user_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

func (r *CreateUserAPIRequest) Parse(req *http.Request) error {
	if err := json.NewDecoder(req.Body).Decode(r); err != nil {
		return err
	}
	return validate.Struct(r)
}

func (r *CreateUserAPIRequest) ToUser() *User {
	return &User{
		ID:       uuid.New(),
		UserName: r.UserName,
		Email:    r.Email,
	}
}

type UserAPIResponse struct {
	Message *User
}

func NewCreateUserAPIResponse(user *User) *UserAPIResponse {
	return &UserAPIResponse{
		Message: user,
	}
}

func (ur *UserAPIResponse) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(ur)
}
