package handler

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/nmluci/cache-optimization/internal/util/echttputil"
	"github.com/nmluci/cache-optimization/pkg/dto"
	"github.com/nmluci/cache-optimization/pkg/errs"
)

type CheckoutHandler func(ctx context.Context, payload *dto.PublicCheckout) (err error)

func HandleCheckout(handler CheckoutHandler) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		req := &dto.PublicCheckout{}
		if err = c.Bind(req); err != nil {
			err = errs.ErrBadRequest
			return echttputil.WriteErrorResponse(c, err)
		}

		req.SessionKey = c.Request().Header.Get("Session-Id")
		err = handler(c.Request().Context(), req)
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, nil)
	}
}

func HandleNCCheckout(handler CheckoutHandler) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		req := &dto.PublicCheckout{}
		if err = c.Bind(req); err != nil {
			err = errs.ErrBadRequest
			return echttputil.WriteErrorResponse(c, err)
		}

		req.SessionKey = c.Request().Header.Get("Session-Id")
		err = handler(c.Request().Context(), req)
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, nil)
	}
}
