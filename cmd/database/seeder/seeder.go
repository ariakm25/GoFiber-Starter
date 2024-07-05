package main

import (
	"GoFiber-API/cmd/database/seeder/seeds"
	database "GoFiber-API/external/database/postgres"
	"GoFiber-API/internal/config"

	"flag"
	"fmt"
	"log"
	"strconv"
)

func main() {

	config.InitConfig(".")

	err := database.ConnectDB(config.GetConfig.DB_HOST, config.GetConfig.DB_PORT, config.GetConfig.DB_USER, config.GetConfig.DB_PASSWORD, config.GetConfig.DB_NAME, config.GetConfig.DB_SSL_MODE)
	if err != nil {
		log.Fatalf("Error connect to Database: %s", err)
		panic(err)
	}

	tableFlag := flag.String("table", "all_table", "specify the table")
	countFlag := flag.String("count", "1", "specify the count")

	flag.Parse()

	countStr := *countFlag
	count, err := strconv.Atoi(countStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	table := *tableFlag

	var allTableSeeder = func(count int) {
		seeds.UserSeeder(count)
	}

	// map table name to seeder function
	seederFuncs := map[string]func(int){
		"all":   allTableSeeder,
		"users": seeds.UserSeeder,
	}

	seederFunc, ok := seederFuncs[table]
	if !ok {
		fmt.Println("Invalid table name. Please specify a valid table name.")
		fmt.Println("Example: go run cmd/database/seeder/seeder.go -table=users -count=10")
		return
	}

	seederFunc(count)

}
