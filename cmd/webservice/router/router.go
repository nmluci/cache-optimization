package router

import (
	"github.com/labstack/echo/v4"
	"github.com/nmluci/cache-optimization/cmd/webservice/handler"
	"github.com/nmluci/cache-optimization/internal/config"
	"github.com/nmluci/cache-optimization/internal/middleware"
	"github.com/nmluci/cache-optimization/internal/repository"
	"github.com/nmluci/cache-optimization/internal/service"
	"github.com/sirupsen/logrus"

	e_middleware "github.com/labstack/echo/v4/middleware"
)

type InitRouterParams struct {
	Logger  *logrus.Entry
	Service service.Service
	Repo    repository.Repository
	Ec      *echo.Echo
	Conf    *config.Config
}

func Init(params *InitRouterParams) {
	params.Ec.GET(PingPath, handler.HandlePing(params.Service.Ping))

	params.Ec.Use(e_middleware.CORS())

	// Cacheable
	params.Ec.POST(AuthRegisterPath, handler.HandleRegisterUser(params.Service.Register))
	params.Ec.POST(AuthLoginPath, handler.HandleLoginUser(params.Service.Login))

	params.Ec.PUT(UserIDPath, handler.HandleEditUser(params.Service.EditUser), middleware.SessionAuthenticator(params.Repo, 1, 2))
	params.Ec.DELETE(UserIDPath, handler.HandleDeleteUser(params.Service.DeleteUser), middleware.SessionAuthenticator(params.Repo, 1, 2))

	params.Ec.GET(ProductsPath, handler.HandleAllProduct(params.Service.FindProducts), middleware.SessionAuthenticator(params.Repo, 1, 2))
	params.Ec.GET(ProductsIDPath, handler.HandleProductDetail(params.Service.FindProductByID), middleware.SessionAuthenticator(params.Repo, 1, 2))
	params.Ec.POST(ProductsPath, handler.HandleStoreProduct(params.Service.InsertProduct), middleware.SessionAuthenticator(params.Repo, 2))
	params.Ec.PUT(ProductsIDPath, handler.HandleEditProduct(params.Service.UpdateProduct), middleware.SessionAuthenticator(params.Repo, 2))

	params.Ec.POST(OrderCheckoutPath, handler.HandleCheckout(params.Service.Checkout), middleware.SessionAuthenticator(params.Repo, 1, 2))

	// Non Cacheable
	params.Ec.POST(NCAuthLoginPath, handler.HandleNCLoginUser(params.Service.ForceLogin))
	params.Ec.PUT(NCUserIDPath, handler.HandleNCEditUser(params.Service.ForceEditUser), middleware.SessionAuthenticatorNoCache(params.Repo, 1, 2))
	params.Ec.GET(NCProductsPath, handler.HandleNCAllProduct(params.Service.ForceFindProducts), middleware.SessionAuthenticatorNoCache(params.Repo, 1, 2))
	params.Ec.GET(NCProductsIDPath, handler.HandleNCProductDetail(params.Service.ForceFindProductByID), middleware.SessionAuthenticatorNoCache(params.Repo, 1, 2))
}
