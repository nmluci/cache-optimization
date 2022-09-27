package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/nmluci/cache-optimization/internal/repository"
	"github.com/nmluci/cache-optimization/internal/util/echttputil"
	"github.com/nmluci/cache-optimization/pkg/errs"
)

func SessionAuthenticator(r repository.Repository, authPriv ...uint64) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			sessionKey := c.Request().Header.Get("Session-Id")
			if sessionKey == "" {
				return echttputil.WriteErrorResponse(c, errs.ErrUnauthorized)
			}

			user, err := r.FindUserBySession(context.TODO(), sessionKey)
			if err != nil {
				return echttputil.WriteErrorResponse(c, errs.ErrUnauthorized)
			}

			for _, priv := range authPriv {
				if user.Priv == priv {
					return next(c)
				}
			}

			return echttputil.WriteErrorResponse(c, errs.ErrUnauthorized)
		}
	}
}

func SessionAuthenticatorNoCache(r repository.Repository, authPriv ...uint64) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			sessionKey := c.Request().Header.Get("Session-Id")
			if sessionKey == "" {
				return echttputil.WriteErrorResponse(c, errs.ErrUnauthorized)
			}

			user, err := r.FindUserSessionByKey(context.TODO(), sessionKey)
			if err != nil {
				return echttputil.WriteErrorResponse(c, errs.ErrUnauthorized)
			}

			for _, priv := range authPriv {
				if user.Priv == priv {
					return next(c)
				}
			}

			return echttputil.WriteErrorResponse(c, errs.ErrUnauthorized)
		}
	}
}
