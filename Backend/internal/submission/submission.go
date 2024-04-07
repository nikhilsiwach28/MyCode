package submission

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikhilsiwach28/MyCode.git/database"
	"github.com/nikhilsiwach28/MyCode.git/models"
)

type Service interface {
	CreateSubmission(context.Context, *models.Submission, []byte) (*models.Submission, models.AppError)
	GetSubmission(context.Context, uuid.UUID) (*models.Submission, models.AppError)
	GetUserSubmissions(context.Context, uuid.UUID) ([]models.Submission, models.AppError)
}
type service struct {
	repo database.Repository
	blob database.BlobService
}

func New(repo database.Repository, blob database.BlobService) *service {
	return &service{repo: repo, blob: blob}
}

func (s service) CreateSubmission(ctx context.Context, submission *models.Submission, file []byte) (*models.Submission, models.AppError) {

	// Inserting in S3
	if err := s.blob.InsertObject(submission.InputFileS3Key, file); err != nil {
		fmt.Println("Error Inserting into S3 ", err)
	}

	// Inserting in DB
	submission, err := s.repo.CreateSubmission(submission)
	if err != nil {
		fmt.Println("Error Creating Submission", err)
		return nil, models.InternalError
	}
	return submission, models.NoError
}

func (s service) GetSubmission(ctx context.Context, submissionID uuid.UUID) (*models.Submission, models.AppError) {

	// Getting from DB
	submission, err := s.repo.GetSubmissionByID(submissionID)
	if err != nil {
		fmt.Println("Error Fetching Submission", err)
		return nil, models.SubmissionNotFoundError
	}

	// Getting from S3
	// TODO: Send this file to Frontend
	file, _ := s.blob.GetObject(submission.InputFileS3Key)
	fmt.Println("file from S3 =", file)

	return submission, models.NoError
}

func (s service) GetUserSubmissions(ctx context.Context, userId uuid.UUID) ([]models.Submission, models.AppError) {

	submission, err := s.repo.GetUserSubmissions(userId)
	if err != nil {
		fmt.Println("Error Fetching Submission", err)
		return nil, models.SubmissionNotFoundError
	}
	return submission, models.NoError
}
