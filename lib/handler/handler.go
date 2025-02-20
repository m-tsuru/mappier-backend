package handler

import (
	_ "time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/redirect"
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

	// app.Use(limiter.New(limiter.Config{
	// 	Max:        1000,
    //     Expiration: 10 * time.Second,
    // }))

	app.Use(redirect.New(redirect.Config{
		Rules: map[string]string{
		  "/app/v1/": "/app/v1/set.html",
		  "/app/v1": "/app/v1/set.html",
		  "/app/":   "/app/v1/set.html",
		  "/app":   "/app/v1/set.html",
		  "/": "/app/v1/set.html",
		},
		StatusCode: 301,
	}))

	auth := app.Group("/auth")
	auth.Get("/login", HandleAuth)
	auth.Get("/callback", HandleAuthCallback)
	auth.Get("/me", HandleMe)

	status := app.Group("/status")
	status.Get("/playing", HandlePlaying)
	status.Get("/Islogin", HandleIsLogin)

	get := app.Group("/get")
	get.Get("/location", HandleLocate)

	app.Static("/app", "./app")

	app.Listen(":8080")
}
