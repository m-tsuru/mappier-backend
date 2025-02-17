package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m-tsuru/mappier-backend/lib/db"
)

func HandleMe(c *fiber.Ctx) error {
	d, err := db.Open()
	if err != nil {
		return c.Status(500).JSON(
			Error{
				Error: err.Error(),
			},
		)
	}

	ck := c.Cookies("SessionId")
	userId, err := db.GetUserfromSessionId(d, ck)
	if err != nil {
		return c.Status(500).JSON(
			Error{
				Error: err.Error(),
			},
		)
	}

	user, err := db.GetUser(d, *userId)
	if err != nil {
		return c.Status(500).JSON(
			Error{
				Error: err.Error(),
			},
		)
	}
	return c.Status(200).JSON(user)
}
