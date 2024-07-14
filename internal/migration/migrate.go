package migration

import (
	user_entities "GoFiber-API/app/user/entities"
	database "GoFiber-API/external/database/postgres"
	"fmt"
	"log"
)

func Migrate() {
	err := database.Connection.AutoMigrate(
		&user_entities.User{},
		&user_entities.UserToken{},
	)

	if err != nil {
		log.Fatal("Error Auto Migrating the Database: ", err)
	}

	fmt.Println("Database Auto Migrated Successfully!")
}
