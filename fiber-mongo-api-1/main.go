package main

import (
	"fiber-mongo-api/mongodb"
	"fiber-mongo-api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Resource clean up on scope exit
	defer mongodb.MongoClient.Close()

	app := fiber.New()

	//enable cors
	app.Use(cors.New(cors.Config{AllowOrigins: "*"}))

	// Handle api routes
	routes.UserRoute(app)

	// Listen on port 8080
	app.Listen(":8080")
}
