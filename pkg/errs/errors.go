package errs

import (
	"errors"
	"net/http"

	"github.com/nmluci/go-backend/pkg/constant"
	"github.com/nmluci/go-backend/pkg/dto"
)

var (
	ErrBadRequest = errors.New("bad request")
	ErrUnknown    = errors.New("internal server error")
)

const (
	ErrCodeUndefined  constant.ErrCode = 1
	ErrCodeBadRequest constant.ErrCode = 2
)

const (
	ErrStatusUnknown     = http.StatusInternalServerError
	ErrStatusClient      = http.StatusBadRequest
	ErrStatusNotLoggedIn = http.StatusUnauthorized
	ErrStatusNoAccess    = http.StatusForbidden
)

var errorMap = map[error]dto.ErrorResponse{
	ErrUnknown:    ErrorResponse(ErrStatusUnknown, ErrCodeUndefined, ErrUnknown),
	ErrBadRequest: ErrorResponse(ErrStatusClient, ErrCodeBadRequest, ErrBadRequest),
}

func ErrorResponse(status int, code constant.ErrCode, err error) dto.ErrorResponse {
	return dto.ErrorResponse{
		Status:  status,
		Code:    code,
		Message: err.Error(),
	}
}

func GetErrorResp(err error) (errResponse dto.ErrorResponse) {
	errResponse, ok := errorMap[err]
	if !ok {
		errResponse = errorMap[ErrUnknown]
	}

	return
}
