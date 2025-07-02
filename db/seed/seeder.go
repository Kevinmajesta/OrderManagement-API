package seeder

import (
	"log"
	"time"

	"Kevinmajesta/OrderManagementAPI/internal/entity"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) {
	var count int64
	db.Model(&entity.User{}).Where("email = ?", "admin@gmail.com").Count(&count)
	if count > 0 {
		log.Println("✅ Admin already exists")
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := &entity.User{
		UserId:       uuid.New(),
		Email:        "admin@gmail.com",
		Password:     string(hashedPassword),
		Fullname:     "Admin Niaga",
		Phone:        "-",
		Role:         "admin",
		Verification: true,
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Fatalf("❌ Failed to seed admin: %v", err)
	}
	log.Println("✅ Admin seeded successfully")
}

func SeedUser(db *gorm.DB) {
	var count int64
	db.Model(&entity.User{}).Where("email = ?", "user@gmail.com").Count(&count)
	if count > 0 {
		log.Println("✅ User already exists")
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	user := &entity.User{
		UserId:       uuid.New(),
		Email:        "user@gmail.com",
		Password:     string(hashedPassword),
		Fullname:     "User Biasa",
		Phone:        "08123456789",
		Role:         "user",
		Verification: true,
	}

	if err := db.Create(&user).Error; err != nil {
		log.Fatalf("❌ Failed to seed user: %v", err)
	}
	log.Println("✅ User seeded successfully")
}

func SeedProducts(db *gorm.DB) {
	var count int64
	db.Model(&entity.Products{}).Count(&count)
	if count > 0 {
		log.Println("✅ Products already seeded")
		return
	}

	products := []entity.Products{
		{
			ProductID:   uuid.New(),
			Name:        "Kemeja Lengan Panjang",
			Description: "Kemeja bahan katun premium",
			Price:       150000,
			Stock:       20,
			PhotoURL:    "/assets/photos/kemeja1.jpg",
		},
		{
			ProductID:   uuid.New(),
			Name:        "Celana Jeans Slim Fit",
			Description: "Celana jeans biru cocok untuk sehari-hari",
			Price:       200000,
			Stock:       15,
			PhotoURL:    "/assets/photos/jeans1.jpg",
		},
		{
			ProductID:   uuid.New(),
			Name:        "Sepatu Sneakers",
			Description: "Sneakers kekinian warna putih",
			Price:       300000,
			Stock:       10,
			PhotoURL:    "/assets/photos/sneakers1.jpg",
		},
	}

	for _, p := range products {
		p.CreatedAt = time.Now()
		p.UpdatedAt = time.Now()
		if err := db.Create(&p).Error; err != nil {
			log.Printf("❌ Failed to seed product %s: %v", p.Name, err)
		}
	}

	log.Println("✅ Products seeded successfully")
}
