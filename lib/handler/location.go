package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m-tsuru/mappier-backend/lib/tools"
)

func HandleLocate(c *fiber.Ctx) (error) {
	ck := c.Cookies("SessionId")
	if ck == "" {
		return c.Status(403).JSON(
			Error{
				Error: "not logged in",
			},
		)
	}

	t, err := tools.GetBuildingFromLatLon(c.QueryFloat("lat"), c.QueryFloat("lon"))
	if err != nil {
		return c.Status(500).JSON(
			Error{
				Error: err.Error(),
			},
		)
	}
	res, err := tools.GetBuildingOfLatLon(*t)
	if err != nil {
		return c.Status(500).JSON(
			Error{
				Error: err.Error(),
			},
		)
	}
	return c.Status(200).JSON(res)
}
