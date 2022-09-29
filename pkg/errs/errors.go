package errs

import (
	"errors"
	"net/http"

	"github.com/nmluci/cache-optimization/pkg/constant"
	"github.com/nmluci/cache-optimization/pkg/dto"
)

var (
	ErrBadRequest   = errors.New("bad request")
	ErrUnknown      = errors.New("internal server error")
	ErrUnauthorized = errors.New("unauthorized")
	ErrDuplicated   = errors.New("duplicated data found")
	ErrNotFound     = errors.New("data not found")
	// ErrDummyError   = errors.New("catched an dummy")
)

const (
	ErrCodeUndefined    constant.ErrCode = 1
	ErrCodeBadRequest   constant.ErrCode = 2
	ErrCodeUnauthorized constant.ErrCode = 3
	ErrCodeDuplicated   constant.ErrCode = 4
	ErrCodeNotFound     constant.ErrCode = 5
	ErrCodeDummyError   constant.ErrCode = 99
)

const (
	ErrStatusUnknown     = http.StatusInternalServerError
	ErrStatusClient      = http.StatusBadRequest
	ErrStatusNotLoggedIn = http.StatusUnauthorized
	ErrStatusNoAccess    = http.StatusForbidden
	ErrStatusDuplicated  = http.StatusConflict
	ErrStatusNotFound    = http.StatusNotFound
)

var errorMap = map[error]dto.ErrorResponse{
	ErrUnknown:      ErrorResponse(ErrStatusUnknown, ErrCodeUndefined, ErrUnknown),
	ErrBadRequest:   ErrorResponse(ErrStatusClient, ErrCodeBadRequest, ErrBadRequest),
	ErrUnauthorized: ErrorResponse(ErrStatusNoAccess, ErrCodeUnauthorized, ErrUnauthorized),
	ErrDuplicated:   ErrorResponse(ErrStatusDuplicated, ErrCodeDuplicated, ErrDuplicated),
	ErrNotFound:     ErrorResponse(ErrStatusNotFound, ErrCodeNotFound, ErrNotFound),
	// ErrDummyError:   ErrorResponse(ErrStatusClient, ErrCodeDummyError, ErrDummyError),
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
