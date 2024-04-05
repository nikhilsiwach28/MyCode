package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikhilsiwach28/MyCode.git/database"
	"github.com/nikhilsiwach28/MyCode.git/models"
)

type Service interface {
	CreateUser(context.Context, *models.User) (*models.User, models.AppError)
	GetUser(context.Context, uuid.UUID) (*models.User, models.AppError)
}
type service struct {
	repo database.Repository
}

func New(repo database.Repository) *service {
	return &service{repo: repo}
}

func (s service) CreateUser(ctx context.Context, user *models.User) (*models.User, models.AppError) {
	user, err := s.repo.CreateUser(user)
	if err != nil {
		fmt.Println("Error Creating User", err)
		return nil, models.InternalError
	}
	return user, models.NoError
}

func (s service) GetUser(ctx context.Context, userID uuid.UUID) (*models.User, models.AppError) {

	user, err := s.repo.GetUser(userID)
	if err != nil {
		fmt.Println("Error Fetching User", err)
		return nil, models.UserNotFoundError
	}
	return user, models.NoError
}
