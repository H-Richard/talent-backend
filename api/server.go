package api

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/H-Richard/talent/api/controllers"
	"github.com/H-Richard/talent/api/seeds"
)

var server = controllers.Server{}

func Start() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting environment, %v", err)
	}
	server.Initialize(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	seeds.Seed(server.DB)
	server.Run(":8080")
}