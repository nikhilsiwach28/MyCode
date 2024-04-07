package models

type RequestMessage struct {
	ID         		string `json:"id"`
	InputFileS3Key  string `json:"input_file_s3_key"`
	CreatedBy  		string `json:"created_by"`
	CreatedAt  		string `json:"created_at"`
	RunTime    		string `json:"run_time"`
	Language   		string `json:"lang"`
	Status     		string `json:"status"`
	OutputFile 		string `json:"output_file"`
}

type ResponseMessage struct {
	Key   string `json:"submissionId"`
	Value string `json:"outputKey"`
}
