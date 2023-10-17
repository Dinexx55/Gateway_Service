package mapper

import (
	"GatewayService/internal/service"
	"net/http"
)

type ErrorMapper struct {
	mapper ErrorMap
}

func NewAuthErrorMapper() ErrorMapper {
	authErrMap := NewAuthErrMap()
	mapper := ErrorMapper{mapper: authErrMap}
	return mapper
}

type ErrorInfo struct {
	StatusCode int
	Message    string
}

type ErrorMap map[error]ErrorInfo

func (m ErrorMapper) MapError(err error) ErrorInfo {
	if value, ok := m.mapper[err]; ok {
		return value
	}

	inf := ErrorInfo{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal server error",
	}
	return inf
}

func NewAuthErrMap() ErrorMap {
	return ErrorMap{
		service.ErrUserNotFound:    {StatusCode: http.StatusBadRequest, Message: "User with provided login does not exist"},
		service.ErrInvalidPassword: {StatusCode: http.StatusBadRequest, Message: "Wrong password provided"},
	}
}
