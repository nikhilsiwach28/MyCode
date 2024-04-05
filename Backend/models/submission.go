package models

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	Accepted            Status = "Accepted"              // normal
	MemoryLimitExceeded Status = "Memory Limit Exceeded" // mle
	TimeLimitExceeded   Status = "Time Limit Exceeded"   // tle
	OutputLimitExceeded Status = "Output Limit Exceeded" // ole
	FileError           Status = "File Error"            // fe
	NonzeroExitStatus   Status = "Nonzero Exit Status"
	Signalled           Status = "Signalled"
	InternalErrorStatus Status = "Internal Error" // system error
	Queued              Status = "Queued"
	Running             Status = "Running"
)

type Submission struct {
	ID        uuid.UUID               `json:"id" validate:"required" gorm:"primaryKey"`
	Link      string                  `json:"link"`
	CreatedBy uuid.UUID               `json:"created_by"`
	CreatedAt time.Time               `json:"created_at" validate:"required" gorm:"default:CURRENT_TIMESTAMP"`
	RunTime   string                  `json:"run_time" validate:"required"`
	Lang      ProgrammingLanguageEnum `json:"lang" validate:"required"`
	Status    Status                  `json:"status"`
	Solution  string                  `json:"solution" gorm:"-"`
}

type CreateSubmissionAPIRequest struct {
	Solution  string                  `json:"solution" validate:"required"`
	CreatedBy uuid.UUID               `json:"created_by"`
	Lang      ProgrammingLanguageEnum `json:"lang" validate:"required"`
}

func (r *CreateSubmissionAPIRequest) Parse(req *http.Request) error {
	if err := json.NewDecoder(req.Body).Decode(r); err != nil {
		return err
	}
	return validate.Struct(r)
}

func (r *CreateSubmissionAPIRequest) ToSubmissions() *Submission {
	return &Submission{
		ID:        uuid.New(),
		Solution:  r.Solution,
		CreatedBy: r.CreatedBy,
		CreatedAt: time.Now(),
		Lang:      r.Lang,
		Status:    Queued,
	}
}

type SubmissionAPIResponse struct {
	Message *Submission
}

func NewCreateSubmissionAPIResponse(submission *Submission) *SubmissionAPIResponse {
	return &SubmissionAPIResponse{
		Message: submission, // Pass UUID directly
	}
}

func (sr *SubmissionAPIResponse) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(sr)
}

type UserSubmissionsAPIResponse struct {
	Message []Submission
}

func NewUserSubmissionsAPIResponse(submission []Submission) *UserSubmissionsAPIResponse {
	return &UserSubmissionsAPIResponse{
		Message: submission, // Pass UUID directly
	}
}

func (sr *UserSubmissionsAPIResponse) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(sr)
}
