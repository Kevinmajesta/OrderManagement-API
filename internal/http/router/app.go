package router

import (
	"net/http"

	"Kevinmajesta/OrderManagementAPI/internal/http/handler"
	"Kevinmajesta/OrderManagementAPI/pkg/route"
)

const (
	Admin = "admin"
	User  = "user"
)

var (
	allRoles  = []string{Admin, User}
	onlyAdmin = []string{Admin}
	onlyUser  = []string{User}
)

func PublicRoutes(userHandler handler.UserHandler,
	adminHandler handler.AdminHandler) []*route.Route {
	return []*route.Route{
		{
			Method:  http.MethodPost,
			Path:    "/login",
			Handler: userHandler.LoginUser,
		},
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: userHandler.CreateUser,
		},
		{
			Method:  http.MethodPost,
			Path:    "/login/admin",
			Handler: adminHandler.LoginAdmin,
		},
		{
			Method:  http.MethodPost,
			Path:    "/admins",
			Handler: adminHandler.CreateAdmin,
		},
		{
			Method:  http.MethodPost,
			Path:    "/password-reset-request",
			Handler: userHandler.RequestPasswordReset,
		},
		{
			Method:  http.MethodPost,
			Path:    "/verification-account",
			Handler: userHandler.VerifUser,
		},
		{
			Method:  http.MethodPost,
			Path:    "/password-reset",
			Handler: userHandler.ResetPassword,
		},
	}
}

func PrivateRoutes(userHandler handler.UserHandler,
	adminHandler handler.AdminHandler, productHandler handler.ProductHandler) []*route.Route {
	return []*route.Route{

		{
			Method:  http.MethodPut,
			Path:    "/users/:user_id",
			Handler: userHandler.UpdateUser,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/users/:user_id",
			Handler: userHandler.DeleteUser,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: adminHandler.FindAllUser,
			Roles:   onlyAdmin,
		},

		{
			Method:  http.MethodPut,
			Path:    "/admins/:user_id",
			Handler: adminHandler.UpdateAdmin,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/admins/:user_id",
			Handler: adminHandler.DeleteAdmin,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/:user_id",
			Handler: userHandler.GetUserProfile,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPost,
			Path:    "/products",
			Handler: productHandler.CreateProduct,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodPut,
			Path:    "/products/:product_id",
			Handler: productHandler.UpdateProduct,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/products/:product_id",
			Handler: productHandler.DeleteProduct,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodGet,
			Path:    "/products",
			Handler: productHandler.FindAllProduct,
			Roles:   allRoles,
		},
	}
}
