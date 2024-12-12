package main

import (
	"log"
	"flag"
	"github.com/gofiber/fiber/v2"
)

var IsDebug bool

func init() {
	flag.BoolVar(&IsDebug, "D", false, "ログを出力する")
}

func debugAccessHandler(c *fiber.Ctx) error {
	return c.Status(200).JSON(
		fiber.Map{
			"success": true,
			"version": "0.0.1",
		},
	)
}

func defaultAccessHandler(c *fiber.Ctx) error {
	return c.Status(200).JSON(
		fiber.Map{
			"success": true,
		},
	)
}

func main() {
	flag.Parse()

	app := fiber.New()

	if IsDebug {
		app.Get("/", debugAccessHandler)
		log.Fatal(app.Listen(":3000"))
	} else {
		app.Get("/", defaultAccessHandler)
		app.Listen(":3000")
	}

}
