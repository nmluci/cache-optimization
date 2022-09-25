package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/nmluci/go-backend/internal/util/echttputil"
	"github.com/nmluci/go-backend/pkg/dto"
)

type PingHandler func() (pingResponse dto.PublicPingResponse)

func HandlePing(handler PingHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		resp := handler()
		return echttputil.WriteSuccessResponse(c, resp)
	}
}
