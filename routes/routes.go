package routes

import (
	"hospital-system/handlers"
	"hospital-system/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	api := app.Group("/api/v1", middleware.GeneralAPIRateLimit())

	auth := api.Group("/auth")
	auth.Post("/register", middleware.RegisterRateLimit(), handlers.Register)
	auth.Post("/login", middleware.AuthRateLimit(), handlers.Login)
	auth.Post("/refresh", middleware.AuthRateLimit(), handlers.RefreshToken)

	auth.Get("/mfa/setup/:user_id", handlers.InitialMFASetup)
	auth.Post("/mfa/setup/:user_id/verify", handlers.VerifyInitialMFASetup)

	protected := api.Group("/", middleware.JWTMiddleware(), middleware.MedicalRateLimit())

	usuarios := protected.Group("/usuarios")
	usuarios.Get("/", middleware.RequirePermission("usuarios", "read"), handlers.GetUsuarios)
	usuarios.Get("/:id", middleware.RequirePermission("usuarios", "read"), handlers.GetUsuario)
	usuarios.Post("/", middleware.RequirePermission("usuarios", "create"), middleware.AdminRateLimit(), handlers.CreateUsuario)
	usuarios.Put("/:id", middleware.RequirePermission("usuarios", "update"), handlers.UpdateUsuario)
	usuarios.Delete("/:id", middleware.RequirePermission("usuarios", "delete"), middleware.AdminRateLimit(), handlers.DeleteUsuario)

	consultorios := protected.Group("/consultorios")
	consultorios.Get("/", middleware.RequirePermission("consultorios", "read"), handlers.GetConsultorios)
	consultorios.Get("/:id", middleware.RequirePermission("consultorios", "read"), handlers.GetConsultorio)
	consultorios.Post("/", middleware.RequirePermission("consultorios", "create"), middleware.AdminRateLimit(), handlers.CreateConsultorio)
	consultorios.Put("/:id", middleware.RequirePermission("consultorios", "update"), middleware.AdminRateLimit(), handlers.UpdateConsultorio)
	consultorios.Delete("/:id", middleware.RequirePermission("consultorios", "delete"), middleware.AdminRateLimit(), handlers.DeleteConsultorio)

	consultas := protected.Group("/consultas")
	consultas.Get("/", middleware.RequirePermission("consultas", "read"), handlers.GetConsultas)
	consultas.Post("/", middleware.RequirePermission("consultas", "create"), handlers.CreateConsulta)

	expedientes := protected.Group("/expedientes")
	expedientes.Get("/", middleware.RequirePermission("expedientes", "read"), handlers.GetExpedientes)
	expedientes.Get("/:id", middleware.RequirePermission("expedientes", "read"), handlers.GetExpediente)
	expedientes.Post("/", middleware.RequirePermission("expedientes", "create"), handlers.CreateExpediente)
	expedientes.Put("/:id", middleware.RequirePermission("expedientes", "update"), handlers.UpdateExpediente)
	expedientes.Delete("/:id", middleware.RequirePermission("expedientes", "delete"), middleware.AdminRateLimit(), handlers.DeleteExpediente)

	horarios := protected.Group("/horarios")
	horarios.Get("/", middleware.RequirePermission("horarios", "read"), handlers.GetHorarios)
	horarios.Get("/:id", middleware.RequirePermission("horarios", "read"), handlers.GetHorario)
	horarios.Post("/", middleware.RequirePermission("horarios", "create"), handlers.CreateHorario)
	horarios.Put("/:id", middleware.RequirePermission("horarios", "update"), handlers.UpdateHorario)
	horarios.Delete("/:id", middleware.RequirePermission("horarios", "delete"), middleware.AdminRateLimit(), handlers.DeleteHorario)

	recetas := protected.Group("/recetas")
	recetas.Get("/", middleware.RequirePermission("recetas", "read"), handlers.GetRecetas)
	recetas.Get("/:id", middleware.RequirePermission("recetas", "read"), handlers.GetReceta)
	recetas.Post("/", middleware.RequirePermission("recetas", "create"), handlers.CreateReceta)
	recetas.Put("/:id", middleware.RequirePermission("recetas", "update"), handlers.UpdateReceta)
	recetas.Delete("/:id", middleware.RequirePermission("recetas", "delete"), middleware.AdminRateLimit(), handlers.DeleteReceta)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Hospital System API is running",
		})
	})

	api.Get("/rate-limit/status", middleware.JWTMiddleware(), middleware.RequireRole("admin"), func(c *fiber.Ctx) error {
		clientIP := c.Query("ip")
		if clientIP == "" {
			clientIP = c.IP()
		}

		status := middleware.GetRateLimitStatus(clientIP)
		return c.JSON(fiber.Map{
			"ip":     clientIP,
			"status": status,
		})
	})

	api.Get("/consultas/:id", handlers.GetConsulta)
	api.Put("/consultas/:id", handlers.UpdateConsulta)
	api.Delete("/consultas/:id", handlers.DeleteConsulta)

	logs := api.Group("/logs", middleware.JWTMiddleware(), middleware.RequireRole("admin"))
	logs.Get("/", handlers.GetLogs)
	logs.Delete("/cleanup", handlers.DeleteOldLogs)
}
