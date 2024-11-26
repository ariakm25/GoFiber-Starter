package seeds

import (
	"GoFiber-API/entities"
	database "GoFiber-API/external/database/postgres"
	internal_casbin "GoFiber-API/internal/casbin"
	internal_log "GoFiber-API/internal/log"
	"GoFiber-API/internal/utils"

	"log"
	"time"

	"github.com/bxcodec/faker/v3"
)

func UserSeeder(count int) {
	log.Println("User Seeder is running for count:", count)

	insertedCount := 0

	hashedPassword, _ := utils.HashPassword("password")

	// Admin Seeder
	adminUser := &entities.User{
		Name:      "Admin",
		Email:     "admin@admin.com",
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	adminQ := database.Connection.FirstOrCreate(adminUser, "email = ?", adminUser.Email)

	if adminQ.Error != nil {
		internal_log.Logger.Error("Error creating Admin: " + adminQ.Error.Error())
	}

	if adminQ.RowsAffected > 0 {
		insertedCount++
	}

	internal_casbin.CasbinEnforcer.AddPermissionForUser("admin", "user", "create")
	internal_casbin.CasbinEnforcer.AddPermissionForUser("admin", "user", "read")
	internal_casbin.CasbinEnforcer.AddPermissionForUser("admin", "user", "update")
	internal_casbin.CasbinEnforcer.AddPermissionForUser("admin", "user", "delete")
	internal_casbin.CasbinEnforcer.AddPermissionForUser("admin", "user", "restore")

	internal_casbin.CasbinEnforcer.AddRoleForUser(adminUser.UID, "admin")

	for i := 0; i < count; i++ {
		newUser := &entities.User{
			Name:      faker.FirstName(),
			Email:     faker.Email(),
			Password:  hashedPassword,
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
