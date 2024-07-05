package seeds

import (
	"GoFiber-API/app/user"
	database "GoFiber-API/external/database/postgres"
	internal_log "GoFiber-API/internal/log"
	"GoFiber-API/internal/utils"

	"log"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
)

func UserSeeder(count int) {
	log.Println("User Seeder is running for count:", count)

	insertedCount := 0

	hashedPassword, _ := utils.HashPassword("admin")

	// Admin Seeder
	adminUser := &user.User{
		Name:      "Admin",
		Email:     "admin@admin.com",
		Password:  hashedPassword,
		UID:       uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.Connection.FirstOrCreate(adminUser, "email = ?", adminUser.Email).Error; err != nil {
		internal_log.Logger.Error("Error creating Admin: " + err.Error())
	}

	for i := 0; i < count; i++ {

		newUser := &user.User{
			Name:      faker.FirstName(),
			Email:     faker.Email(),
			Password:  hashedPassword,
			UID:       uuid.New().String(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := database.Connection.FirstOrCreate(newUser, "email = ?", newUser.Email).Error; err != nil {
			internal_log.Logger.Error("Error creating User: " + err.Error())
		} else {
			insertedCount++
		}
	}

	log.Println("User Seeder is done. inserted count:", insertedCount)
}
