// database.go
package database

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/nikhilsiwach28/MyCode.git/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewPostgres(connString string) *repository {
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		fmt.Println("Error Connecting Psql", err)
		return nil
	}

	_ = db.AutoMigrate(&models.User{}, &models.Submission{})

	return &repository{db: db}
}

func (r *repository) CreateUser(user *models.User) (*models.User, error) {
	result := r.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (r *repository) CreateSubmission(submission *models.Submission) (*models.Submission, error) {
	result := r.db.Create(submission)
	if result.Error != nil {
		return nil, result.Error
	}
	return submission, nil
}

func (r *repository) GetUser(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *repository) GetUserSubmissions(userID uuid.UUID) ([] models.Submission, error) {
	var submissions []models.Submission
	err := r.db.Where("created_by = ?", userID).Find(&submissions).Error
	return submissions, err
}

func (r *repository) GetSubmissionByID(submissionID uuid.UUID) (*models.Submission, error) {
	var submission models.Submission
	err := r.db.First(&submission, "id = ?", submissionID).Error
	if err != nil {
		return nil, err
	}
	return &submission, nil
}
