package builder

import (
	"Kevinmajesta/OrderManagement-API/configs"
	"Kevinmajesta/OrderManagement-API/internal/http/handler"
	"Kevinmajesta/OrderManagement-API/internal/http/router"
	"Kevinmajesta/OrderManagement-API/internal/repository"
	"Kevinmajesta/OrderManagement-API/internal/service"
	"Kevinmajesta/OrderManagement-API/pkg/cache"
	"Kevinmajesta/OrderManagement-API/pkg/email"
	"Kevinmajesta/OrderManagement-API/pkg/encrypt"
	"Kevinmajesta/OrderManagement-API/pkg/route"
	"Kevinmajesta/OrderManagement-API/pkg/token"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func BuildPublicRoutes(db *gorm.DB, redisDB *redis.Client, tokenUseCase token.TokenUseCase, encryptTool encrypt.EncryptTool,
	cfg *configs.Config) []*route.Route {
	emailService := email.NewEmailSender(cfg)
	userRepository := repository.NewUserRepository(db, nil)
	userService := service.NewUserService(userRepository, tokenUseCase, encryptTool, emailService)

	userHandler := handler.NewUserHandler(userService)

	adminRepository := repository.NewAdminRepository(db, nil)
	adminService := service.NewAdminService(adminRepository, tokenUseCase, encryptTool, emailService)
	adminHandler := handler.NewAdminHandler(adminService)

	return router.PublicRoutes(userHandler, adminHandler)
}

func BuildPrivateRoutes(db *gorm.DB, redisDB *redis.Client, encryptTool encrypt.EncryptTool, cfg *configs.Config, tokenUseCase token.TokenUseCase) []*route.Route {
	cacheable := cache.NewCacheable(redisDB)
	userRepository := repository.NewUserRepository(db, cacheable)
	userService := service.NewUserService(userRepository, nil, encryptTool, nil)

	userHandler := handler.NewUserHandler(userService)

	adminRepository := repository.NewAdminRepository(db, cacheable)
	adminService := service.NewAdminService(adminRepository, nil, encryptTool, nil)
	adminHandler := handler.NewAdminHandler(adminService)

	return router.PrivateRoutes(userHandler, adminHandler)
}
