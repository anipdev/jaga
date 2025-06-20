package main

import (
	"log"
	"net/http"

	"jaga/config"
	"jaga/repositories"
	"jaga/routes"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, proceeding with system env")
	}
}

func main() {

	db := config.InitDB()
	config.AutoMigrate(db)

	userRepository := repositories.NewUserRepository(db)
	config.SeedSuperUser(userRepository)

	router := routes.RegisterRoutes()

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", router)

}
