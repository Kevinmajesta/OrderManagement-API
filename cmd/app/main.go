package main

import (
	"log"

	"Kevinmajesta/OrderManagementAPI/configs"
	seeder "Kevinmajesta/OrderManagementAPI/db/seed"
	"Kevinmajesta/OrderManagementAPI/internal/builder"
	"Kevinmajesta/OrderManagementAPI/pkg/cache"
	"Kevinmajesta/OrderManagementAPI/pkg/email"
	"Kevinmajesta/OrderManagementAPI/pkg/encrypt"
	"Kevinmajesta/OrderManagementAPI/pkg/midtrans"
	"Kevinmajesta/OrderManagementAPI/pkg/postgres"
	"Kevinmajesta/OrderManagementAPI/pkg/server"
	"Kevinmajesta/OrderManagementAPI/pkg/token"
	"Kevinmajesta/OrderManagementAPI/worker"
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

	// Init Midtrans
	midtransService := midtrans.NewMidtransService(&cfg.Midtrans)

	emailSender := email.NewEmailSender(cfg)
	worker.StartEmailWorker(emailSender)
	worker.StartPhotoWorker()
	seeder.SeedAdmin(db)
	seeder.SeedUser(db)
	seeder.SeedProducts(db)

	// Build Echo route groups
	publicRoutes := builder.BuildPublicRoutes(db, redisDB, tokenUseCase, encryptTool, cfg, midtransService)
	privateRoutes := builder.BuildPrivateRoutes(db, redisDB, encryptTool, cfg, tokenUseCase, midtransService)

	// Start server
	srv := server.NewServer("app", publicRoutes, privateRoutes, cfg.JWT.SecretKey)
	srv.Run()
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}
}
