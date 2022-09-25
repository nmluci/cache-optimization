package echttputil

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nmluci/go-backend/pkg/dto"
	"github.com/nmluci/go-backend/pkg/errs"
)

func WriteSuccessResponse(ec echo.Context, data interface{}) error {
	return ec.JSON(http.StatusOK, dto.BaseResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func WriteErrorResponse(ec echo.Context, err error) error {
	errResp := errs.GetErrorResp(err)
	return ec.JSON(errResp.Status, dto.BaseResponse{
		Code:    errResp.Code,
		Message: errResp.Message,
		Data:    nil,
	})
}
