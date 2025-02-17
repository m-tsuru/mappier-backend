package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m-tsuru/mappier-backend/lib/auth"
	"github.com/m-tsuru/mappier-backend/lib/db"
)

func HandlePlaying(c *fiber.Ctx) error {
	d, err := db.Open()
	if err != nil {
		return c.Status(500).JSON(
			Error{
				Error: err.Error(),
			},
		)
	}

	sessionId := c.Cookies("SessionId")

	accessToken, err := db.GetAccessTokenfromSessionId(d, sessionId)
	if err != nil {
		return c.Status(500).JSON(
			Error{
				Error: err.Error(),
			},
		)
	}

	state, err := auth.GetSpotifyPlayingState(d, *accessToken)
	if err != nil {
		return c.Status(500).JSON(
			Error{
				Error: err.Error(),
			},
		)
	}

	return c.Status(200).JSON(state)
}
