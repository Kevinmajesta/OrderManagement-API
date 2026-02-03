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
	adminHandler handler.AdminHandler,
	midtransHandler *handler.MidtransHandler) []*route.Route {
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
		{
			Method:  http.MethodPost,
			Path:    "/midtrans/notification",
			Handler: midtransHandler.HandleNotification,
		},
	}
}

func PrivateRoutes(userHandler handler.UserHandler,
	adminHandler handler.AdminHandler, productHandler handler.ProductHandler,
	orderHandler handler.OrderHandler, cartHandler *handler.CartHandler) []*route.Route {
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
		{
			Method:  http.MethodPost,
			Path:    "/orders",
			Handler: orderHandler.CreateOrder,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPost,
			Path:    "/cart/items",
			Handler: cartHandler.AddItem,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPut,
			Path:    "/cart/items/:cart_item_id",
			Handler: cartHandler.UpdateItem,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/cart/items/:cart_item_id",
			Handler: cartHandler.RemoveItem,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/cart",
			Handler: cartHandler.GetCart,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPost,
			Path:    "/cart/checkout",
			Handler: cartHandler.Checkout,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/orders/:order_id/status",
			Handler: orderHandler.UpdateOrderStatus,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodGet,
			Path:    "/orders/history",
			Handler: orderHandler.GetOrderHistory,
			Roles:   allRoles,
		},
	}
}
