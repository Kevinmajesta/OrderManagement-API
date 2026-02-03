package builder

import (
	"Kevinmajesta/OrderManagementAPI/configs"
	"Kevinmajesta/OrderManagementAPI/internal/http/handler"
	"Kevinmajesta/OrderManagementAPI/internal/http/router"
	"Kevinmajesta/OrderManagementAPI/internal/repository"
	"Kevinmajesta/OrderManagementAPI/internal/service"
	"Kevinmajesta/OrderManagementAPI/pkg/cache"
	"Kevinmajesta/OrderManagementAPI/pkg/email"
	"Kevinmajesta/OrderManagementAPI/pkg/encrypt"
	"Kevinmajesta/OrderManagementAPI/pkg/midtrans"
	"Kevinmajesta/OrderManagementAPI/pkg/route"
	"Kevinmajesta/OrderManagementAPI/pkg/token"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func BuildPublicRoutes(db *gorm.DB, redisDB *redis.Client, tokenUseCase token.TokenUseCase, encryptTool encrypt.EncryptTool,
	cfg *configs.Config, midtransService *midtrans.MidtransService) []*route.Route {
	EmailSenderService := email.NewEmailSender(cfg)
	userRepository := repository.NewUserRepository(db, nil)
	userService := service.NewUserService(userRepository, tokenUseCase, encryptTool, EmailSenderService)

	userHandler := handler.NewUserHandler(userService)

	adminRepository := repository.NewAdminRepository(db, nil)
	adminService := service.NewAdminService(adminRepository, tokenUseCase, encryptTool, EmailSenderService)
	adminHandler := handler.NewAdminHandler(adminService)

	// Order service for Midtrans webhook
	cacheable := cache.NewCacheable(redisDB)
	orderRepository := repository.NewOrderRepository(db, cacheable)
	orderService := service.NewOrderService(orderRepository, db, midtransService)
	midtransHandler := handler.NewMidtransHandler(orderService)

	return router.PublicRoutes(userHandler, adminHandler, midtransHandler)
}

func BuildPrivateRoutes(db *gorm.DB, redisDB *redis.Client, encryptTool encrypt.EncryptTool, cfg *configs.Config, tokenUseCase token.TokenUseCase, midtransService *midtrans.MidtransService) []*route.Route {
	cacheable := cache.NewCacheable(redisDB)
	userRepository := repository.NewUserRepository(db, cacheable)
	userService := service.NewUserService(userRepository, nil, encryptTool, nil)

	userHandler := handler.NewUserHandler(userService)

	adminRepository := repository.NewAdminRepository(db, cacheable)
	adminService := service.NewAdminService(adminRepository, nil, encryptTool, nil)
	adminHandler := handler.NewAdminHandler(adminService)

	productRepository := repository.NewProductRepository(db, cacheable)
	productService := service.NewProductService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	orderRepository := repository.NewOrderRepository(db, cacheable)
	orderService := service.NewOrderService(orderRepository, db, midtransService)
	orderHandler := handler.NewOrderHandler(orderService)

	cartRepository := repository.NewCartRepository(db)
	cartService := service.NewCartService(cartRepository, orderService, productRepository)
	cartHandler := handler.NewCartHandler(cartService)

	receiptRepository := repository.NewReceiptRepository(db)
	receiptService := service.NewReceiptService(receiptRepository, orderRepository, db)
	receiptHandler := handler.NewReceiptHandler(receiptService)

	return router.PrivateRoutes(userHandler, adminHandler, productHandler, *orderHandler, cartHandler, receiptHandler)
}
