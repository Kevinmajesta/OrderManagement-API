package main

import (
	"Kevinmajesta/OrderManagement-API/configs" 
	"Kevinmajesta/OrderManagement-API/internal/builder"
	"Kevinmajesta/OrderManagement-API/pkg/cache"
	"Kevinmajesta/OrderManagement-API/pkg/encrypt"
	"Kevinmajesta/OrderManagement-API/pkg/postgres"
	"Kevinmajesta/OrderManagement-API/pkg/server"
	"Kevinmajesta/OrderManagement-API/pkg/token"
	"log" 
)

func main() {
	// Load configurations from .env file
	cfg, err := configs.NewConfig(".env")
	checkError(err)

	// Initialize PostgreSQL database connection
	db, err := postgres.InitPostgres(&cfg.Postgres)
	checkError(err)

	// Initialize Redis cache connection
	redisDB := cache.InitCache(&cfg.Redis)

	// Initialize encryption tool
	encryptTool := encrypt.NewEncryptTool(cfg.Encrypt.SecretKey, cfg.Encrypt.IV)

	// Initialize JWT token use case
	tokenUseCase := token.NewTokenUseCase(cfg.JWT.SecretKey)

	// untuk menerima *configs.Config, BUKAN *entity.Config.
	publicRoutes := builder.BuildPublicRoutes(db, redisDB, tokenUseCase, encryptTool, cfg) // Gunakan cfg
	privateRoutes := builder.BuildPrivateRoutes(db, redisDB, encryptTool, cfg, tokenUseCase) // Gunakan cfg

	// Initialize and run the server
	srv := server.NewServer("app", publicRoutes, privateRoutes, cfg.JWT.SecretKey)
	srv.Run()
}

func checkError(err error) {
	if err != nil {
		// Gunakan log.Fatalf untuk keluar dengan pesan error yang jelas
		log.Fatalf("Fatal error: %v", err)
	}
}
