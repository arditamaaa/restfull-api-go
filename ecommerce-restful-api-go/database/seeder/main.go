package main

import (
	"fmt"
	"os"
	"simple-api-go/config"
	"simple-api-go/database"
	"simple-api-go/database/seeder/seeds"
	"simple-api-go/util"
)

func main() {
	db := database.Connect(config.DBHost, config.DBName)
	args := os.Args
	availableSeeders := []string{"user_seeder", "product_seeder"}
	initSeeds := []string{}

	if len(args) >= 2 {
		seedArgs := args[1:]
		for _, v := range seedArgs {
			if util.StringInSlice(v, availableSeeders) {
				initSeeds = append(initSeeds, v)
			}
		}
	} else {
		initSeeds = availableSeeders
	}

	for _, v := range initSeeds {
		switch v {
		case "user_seeder":
			seeds.UserSeeder(db)
		case "product_seeder":
			seeds.ProductSeeder(db)
		}
	}

	fmt.Println("seed.go", initSeeds)
}
