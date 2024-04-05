package models

import (
	"fmt"

	"github.com/go-playground/validator"
)

var validate = validator.New()

type AppError struct {
	Type   ErrorType
	Reason string `json:"reason"`
}

func NewAppError(code ErrorType, err string) AppError {
	return AppError{
		Type:   code,
		Reason: err,
	}
}

func (e AppError) Add(err error) AppError {
	fmt.Println(err)
	e.Reason = fmt.Sprintf("%s : %s", e.Reason, err.Error())
	return e
}

func (e AppError) Error() string {
	return fmt.Sprintf("%s, %s", e.Type, e.Reason)
}

var (
	UserNotFoundError        = NewAppError(ErrorNotFound, "user not found")
	EnvironmentNotFoundError = NewAppError(ErrorNotFound, "environment not found")
	SubmissionNotFoundError  = NewAppError(ErrorNotFound, "submission not found")
	UserInvalidInput         = NewAppError(ErrorBadData, "invalid input")
	InternalError            = NewAppError(ErrorInternal, "internal server error")
	BadRequest               = NewAppError(ErrorBadData, "bad request")
	NoError                  = NewAppError(ErrorNone, "")
	UnauthorizedError        = NewAppError(ErrorUnauthorized, "unauthorized")
)

type ErrorType string

const (
	ErrorNone          ErrorType = ""
	ErrorTimeout       ErrorType = "timeout"
	ErrorCanceled      ErrorType = "canceled"
	ErrorExec          ErrorType = "execution"
	ErrorBadData       ErrorType = "bad_data"
	ErrorInternal      ErrorType = "internal"
	ErrorUnavailable   ErrorType = "unavailable"
	ErrorNotFound      ErrorType = "not_found"
	ErrorNotAcceptable ErrorType = "not_acceptable"
	ErrorUnauthorized  ErrorType = "unauthorized"
)

type AccessLevelModeEnum string

const (
	AccessLevelAdmin AccessLevelModeEnum = "ADMIN"
	AccessLevelUser  AccessLevelModeEnum = "USER"
	AccessLevelGuest AccessLevelModeEnum = "GUEST"
)

type ProgrammingLanguageEnum string

const (
	C          ProgrammingLanguageEnum = "C"
	CPlusPlus  ProgrammingLanguageEnum = "C++"
	Java       ProgrammingLanguageEnum = "Java"
	Python     ProgrammingLanguageEnum = "Python"
	JavaScript ProgrammingLanguageEnum = "JavaScript"
	Go         ProgrammingLanguageEnum = "Go"
	Rust       ProgrammingLanguageEnum = "Rust"
	Text       ProgrammingLanguageEnum = "Text"
	YAML       ProgrammingLanguageEnum = "YAML"
	MYSQL      ProgrammingLanguageEnum = "MYSQL"
	// Add more programming languages as needed
)

func (lang ProgrammingLanguageEnum) GetExtension() string {
	switch lang {
	case C:
		return ".c"
	case CPlusPlus:
		return ".cpp"
	case Java:
		return ".java"
	case Python:
		return ".py"
	case JavaScript:
		return ".js"
	case Go:
		return ".go"
	case Rust:
		return ".rs"
	case Text:
		return ".txt"
	case YAML:
		return ".yml"
	default:
		// Handle unknown language or return a default extension
		return ".txt"
	}
}
