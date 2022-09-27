package handler

import (
	"context"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/nmluci/cache-optimization/internal/model"
	"github.com/nmluci/cache-optimization/internal/util/echttputil"
	"github.com/nmluci/cache-optimization/pkg/errs"
)

type ProductDetailHandler func(ctx context.Context, id uint64) (res *model.Product, err error)
type AllProductHandler func(ctx context.Context) (res []*model.Product, err error)
type StoreProductHandler func(ctx context.Context, payload *model.Product) (err error)
type UpdateProductHandler func(ctx context.Context, id uint64, payload *model.Product) (err error)
type DeleteProductHandler func(ctx context.Context, id uint64) (err error)

func HandleProductDetail(handler ProductDetailHandler) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		id := c.QueryParam("id")
		parsedId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return echttputil.WriteErrorResponse(c, nil)
		}

		data, err := handler(c.Request().Context(), parsedId)
		if err != nil {
			return echttputil.WriteErrorResponse(c, nil)
		}

		return echttputil.WriteSuccessResponse(c, data)
	}
}

func HandleAllProduct(handler AllProductHandler) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		data, err := handler(c.Request().Context())
		if err != nil {
			return echttputil.WriteErrorResponse(c, nil)
		}

		return echttputil.WriteSuccessResponse(c, data)
	}
}

func HandleNCProductDetail(handler ProductDetailHandler) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		id := c.QueryParam("id")
		parsedId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return echttputil.WriteErrorResponse(c, nil)
		}

		data, err := handler(c.Request().Context(), parsedId)
		if err != nil {
			return echttputil.WriteErrorResponse(c, nil)
		}

		return echttputil.WriteSuccessResponse(c, data)
	}
}

func HandleNCAllProduct(handler AllProductHandler) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		data, err := handler(c.Request().Context())
		if err != nil {
			return echttputil.WriteErrorResponse(c, nil)
		}

		return echttputil.WriteSuccessResponse(c, data)
	}
}

func HandleStoreProduct(handler StoreProductHandler) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		var req *model.Product
		if err = c.Bind(req); err != nil {
			err = errs.ErrBadRequest
			return echttputil.WriteErrorResponse(c, err)
		}

		err = handler(c.Request().Context(), req)
		if err != nil {
			return echttputil.WriteErrorResponse(c, nil)
		}

		return echttputil.WriteSuccessResponse(c, nil)
	}
}

func HandleEditProduct(handler UpdateProductHandler) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		var req *model.Product
		if err = c.Bind(req); err != nil {
			err = errs.ErrBadRequest
			return echttputil.WriteErrorResponse(c, err)
		}

		id := c.QueryParam("id")
		parsedId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return echttputil.WriteErrorResponse(c, nil)
		}

		err = handler(c.Request().Context(), parsedId, req)
		if err != nil {
			return echttputil.WriteErrorResponse(c, nil)
		}

		return echttputil.WriteSuccessResponse(c, nil)
	}
}

func HandleDeleteProduct(handler DeleteProductHandler) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		id := c.QueryParam("id")
		parsedId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return echttputil.WriteErrorResponse(c, nil)
		}

		err = handler(c.Request().Context(), parsedId)
		if err != nil {
			return echttputil.WriteErrorResponse(c, nil)
		}

		return echttputil.WriteSuccessResponse(c, nil)
	}
}
