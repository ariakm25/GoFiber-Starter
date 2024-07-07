package migration

import (
	"GoFiber-API/app/user"
	database "GoFiber-API/external/database/postgres"
	"fmt"
	"log"
)

func Migrate() {
	err := database.Connection.AutoMigrate(
		&user.User{},
		&user.UserToken{},
	)

	if err != nil {
		log.Fatal("Error Auto Migrating the Database: ", err)
	}

	fmt.Println("Database Auto Migrated Successfully!")
}
