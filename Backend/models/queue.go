package models

type QueueMessage struct {
	SubmissionID string `json:"submissionId"`
	OutputKey    string `json:"outputKey"`
}