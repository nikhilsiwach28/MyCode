package database

import (
	"github.com/google/uuid"
	"github.com/nikhilsiwach28/MyCode.git/models"
)

type Repository interface {
	CreateUser(*models.User) (*models.User, error)
	CreateSubmission(*models.Submission) (*models.Submission, error)
	GetUser(uuid.UUID) (*models.User, error)
	GetUserSubmissions(uuid.UUID) ([]models.Submission, error)
	GetSubmissionByID(uuid.UUID) (*models.Submission, error)
}
