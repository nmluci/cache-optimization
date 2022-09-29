package middleware

import (
	"fmt"

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
				fmt.Println("session key is empty")
				return echttputil.WriteErrorResponse(c, errs.ErrUnauthorized)
			}

			user, err := r.FindUserBySession(c.Request().Context(), sessionKey)
			if err != nil {
				fmt.Printf("an error occured while fetching user session: %+v\n", err)
				return echttputil.WriteErrorResponse(c, errs.ErrUnauthorized)
			}

			// revalidate
			val, err := r.FindUserByID(c.Request().Context(), user.ID)
			if err != nil {
				fmt.Println("an error occured while fetching userdata from db")
				return echttputil.WriteErrorResponse(c, errs.ErrUnauthorized)
			} else if val == nil {
				fmt.Println("user not existed anymore")
				r.InvalidateUserSession(c.Request().Context(), sessionKey)
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
				fmt.Println("empty session id")
				return echttputil.WriteErrorResponse(c, errs.ErrUnauthorized)
			}

			user, err := r.FindUserSessionByKey(c.Request().Context(), sessionKey)
			if err != nil {
				fmt.Printf("an error occured while fetching user session: %+v\n", err)
				return echttputil.WriteErrorResponse(c, errs.ErrUnauthorized)
			} else if user == nil {
				r.InvalidateSessionKey(c.Request().Context(), sessionKey)
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
