package router

const (
	PingPath = "/v1/ping"

	UsersPath  = "/v1/users"
	UserIDPath = "/v1/users/{id}"

	ProductsPath   = "/v1/products"
	ProductsIDPath = "/v1/products/{id}"

	OrderCheckoutPath = "/v1/checkout"

	AuthLoginPath    = "/v1/auth/login"
	AuthRegisterPath = "/v1/auth/register"

	// NON CACHE
	NCAuthLoginPath  = "/v1/nc/auth/login"
	NCUserIDPath     = "/v1/nc/users/{id}"
	NCProductsPath   = "/v1/nc/products"
	NCProductsIDPath = "/v1/nc/products/{id}"
)
