package main

import (
	"log"

	"Kevinmajesta/OrderManagementAPI/configs"
	"Kevinmajesta/OrderManagementAPI/internal/builder"
	"Kevinmajesta/OrderManagementAPI/pkg/cache"
	"Kevinmajesta/OrderManagementAPI/pkg/encrypt"
	"Kevinmajesta/OrderManagementAPI/pkg/postgres"
	"Kevinmajesta/OrderManagementAPI/pkg/server"
	"Kevinmajesta/OrderManagementAPI/pkg/token"
	"Kevinmajesta/OrderManagementAPI/worker"
	"Kevinmajesta/OrderManagementAPI/pkg/email"
)

func main() {
	// Load environment variables
	cfg, err := configs.NewConfig(".env")
	checkError(err)

	// Init PostgreSQL DB
	db, err := postgres.InitPostgres(&cfg.Postgres)
	checkError(err)

	// Init Redis
	redisDB := cache.InitCache(&cfg.Redis)

	// Init encryption tools
	encryptTool := encrypt.NewEncryptTool(cfg.Encrypt.SecretKey, cfg.Encrypt.IV)

	// Init JWT generator
	tokenUseCase := token.NewTokenUseCase(cfg.JWT.SecretKey)

	emailSender := email.NewEmailSender(cfg) 
	worker.StartEmailWorker(emailSender)

	// Build Echo route groups
	publicRoutes := builder.BuildPublicRoutes(db, redisDB, tokenUseCase, encryptTool, cfg)
	privateRoutes := builder.BuildPrivateRoutes(db, redisDB, encryptTool, cfg, tokenUseCase)

	// Start server
	srv := server.NewServer("app", publicRoutes, privateRoutes, cfg.JWT.SecretKey)
	srv.Run()
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}
}
