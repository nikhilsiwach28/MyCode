package models

type RequestMessage struct {
	Code     string
	Language string
}

type ResponseMessage struct {
	Key   string
	Value string
}
