package models

type RequestMessage struct {
	ID         string `json:"id"`
	InputFile  string `json:"input_file"`
	CreatedBy  string `json:"created_by"`
	CreatedAt  string `json:"created_at"`
	RunTime    string `json:"run_time"`
	Language   string `json:"lang"`
	Status     string `json:"status"`
	OutputFile string `json:"output_file"`
}

type ResponseMessage struct {
	Key   string `json:"submissionId"`
	Value string `json:"outputKey"`
}
