package handler

import (
	"github.com/gofiber/fiber/v2"
)

type Error struct {
	Error string `json:"error"`
}

type Message struct {
	Message string `json:"message"`
}

type Url struct {
	Url *string
}

func Setup() {
	app := fiber.New()

	auth := app.Group("/auth")
	auth.Get("/login", HandleAuth)
	auth.Get("/callback", HandleAuthCallback)
	auth.Get("/me", HandleMe)

	status := app.Group("/status")
	status.Get("/playing", HandlePlaying)

	app.Listen(":8080")
}
