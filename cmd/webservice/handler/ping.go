package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/nmluci/cache-optimization/internal/util/echttputil"
	"github.com/nmluci/cache-optimization/pkg/dto"
)

type PingHandler func() (pingResponse dto.PublicPingResponse)

func HandlePing(handler PingHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		resp := handler()
		return echttputil.WriteSuccessResponse(c, resp)
	}
}
