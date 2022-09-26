package dto

import "github.com/nmluci/cache-optimization/pkg/constant"

type ErrorResponse struct {
	Status  int
	Code    constant.ErrCode
	Message string
}

type BaseResponse struct {
	Code    constant.ErrCode `json:"code"`
	Message string           `json:"message"`
	Data    interface{}      `json:"data"`
}
