package main

import (
	"hospital-system/config"
	"hospital-system/routes"
	"log"
	"os"

	"hospital-system/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	config.ConnectDB()

	app := fiber.New(fiber.Config{
		AppName: "Hospital System API v1.0.0",
	})

	// Middleware de validación de respuestas
	app.Use(middleware.ResponseValidator()) // Para logging de validación
	// app.Use(middleware.StrictResponseValidator()) // Versión estricta que devuelve error

	// Agregar después de las configuraciones de CORS y antes de las rutas
	app.Use(middleware.LoggerMiddleware())
	routes.SetupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf(" Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
