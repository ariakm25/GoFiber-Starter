package migration

import (
	"GoFiber-API/entities"
	database "GoFiber-API/external/database/postgres"
	"fmt"
	"log"
)

func Migrate() {
	err := database.Connection.AutoMigrate(
		&entities.User{},
		&entities.UserToken{},
		&entities.UserSession{},
	)

	if err != nil {
		log.Fatal("Error Auto Migrating the Database: ", err)
	}

	fmt.Println("Database Auto Migrated Successfully!")
}
