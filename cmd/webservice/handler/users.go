package handler

import (
	"context"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/nmluci/cache-optimization/internal/util/echttputil"
	"github.com/nmluci/cache-optimization/pkg/dto"
	"github.com/nmluci/cache-optimization/pkg/errs"
)

type RegisterUserHandler func(ctx context.Context, payload *dto.PublicUserPayload) (err error)
type LoginUserHandler func(ctx context.Context, payload *dto.PublicUserLoginPayload) (sessionKey string, err error)
type EditUserHandler func(ctx context.Context, id uint64, payload *dto.PublicUserPayload) (err error)
type DeleteUserHandler func(ctx context.Context, id uint64, sessionKey string) (err error)

func HandleRegisterUser(handler RegisterUserHandler) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		var req *dto.PublicUserPayload
		if err = c.Bind(req); err != nil {
			err = errs.ErrBadRequest
			return echttputil.WriteErrorResponse(c, err)
		}

		masterKey := c.Request().Header.Get("X-Misaki")
		if masterKey != "" {
			req.MasterKey = &masterKey
		}

		err = handler(c.Request().Context(), req)
		if err != nil {
			return echttputil.WriteErrorResponse(c, nil)
		}

		return echttputil.WriteSuccessResponse(c, nil)
	}
}

func HandleLoginUser(handler LoginUserHandler) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		var req *dto.PublicUserLoginPayload
		if err = c.Bind(req); err != nil {
			err = errs.ErrBadRequest
			return echttputil.WriteErrorResponse(c, err)
		}

		sessionKey, err := handler(c.Request().Context(), req)
		if err != nil {
			return echttputil.WriteErrorResponse(c, nil)
		}

		c.Response().Header().Set("Session-Id", sessionKey)

		return echttputil.WriteSuccessResponse(c, nil)
	}
}

func HandleEditUser(handler EditUserHandler) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		var req *dto.PublicUserPayload
		if err = c.Bind(req); err != nil {
			err = errs.ErrBadRequest
			return echttputil.WriteErrorResponse(c, err)
		}

		id := c.QueryParam("id")
		parsedId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return echttputil.WriteErrorResponse(c, nil)
		}

		sessionKey := c.Request().Header.Get("Session-Id")
		if sessionKey != "" {
			req.SessionKey = sessionKey
		}

		err = handler(c.Request().Context(), parsedId, req)
		if err != nil {
			return echttputil.WriteErrorResponse(c, nil)
		}

		return echttputil.WriteSuccessResponse(c, nil)
	}
}

func HandleNCLoginUser(handler LoginUserHandler) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		var req *dto.PublicUserLoginPayload
		if err = c.Bind(req); err != nil {
			err = errs.ErrBadRequest
			return echttputil.WriteErrorResponse(c, err)
		}

		sessionKey, err := handler(c.Request().Context(), req)
		if err != nil {
			return echttputil.WriteErrorResponse(c, nil)
		}

		c.Response().Header().Set("Session-Id", sessionKey)

		return echttputil.WriteSuccessResponse(c, nil)
	}
}

func HandleNCEditUser(handler EditUserHandler) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		var req *dto.PublicUserPayload
		if err = c.Bind(req); err != nil {
			err = errs.ErrBadRequest
			return echttputil.WriteErrorResponse(c, err)
		}

		id := c.QueryParam("id")
		parsedId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return echttputil.WriteErrorResponse(c, nil)
		}

		sessionKey := c.Request().Header.Get("Session-Id")
		if sessionKey != "" {
			req.SessionKey = sessionKey
		}

		err = handler(c.Request().Context(), parsedId, req)
		if err != nil {
			return echttputil.WriteErrorResponse(c, nil)
		}

		return echttputil.WriteSuccessResponse(c, nil)
	}
}

func HandleDeleteUser(handler DeleteUserHandler) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		id := c.QueryParam("id")
		parsedId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return echttputil.WriteErrorResponse(c, nil)
		}

		sessionKey := c.Request().Header.Get("Session-Id")

		err = handler(c.Request().Context(), parsedId, sessionKey)
		if err != nil {
			return echttputil.WriteErrorResponse(c, nil)
		}

		return echttputil.WriteSuccessResponse(c, nil)

	}
}
