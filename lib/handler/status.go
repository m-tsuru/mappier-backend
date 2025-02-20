package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m-tsuru/mappier-backend/lib/auth"
	"github.com/m-tsuru/mappier-backend/lib/db"
	"github.com/m-tsuru/mappier-backend/lib/structs"
)

func HandleIsLogin(c *fiber.Ctx) error {
	d, err := db.Open()
	if err != nil {
		return c.Status(500).JSON(
			&Error{
				Error: err.Error(),
			},
		)
	}

	ck := c.Cookies("SessionId")
	_, err = db.GetUserfromSessionId(d, ck)

	if err != nil {
		if err.Error() == "not logged in" {
			return c.Status(200).JSON(&structs.IsLogin{
				IsLogin: false,
			})
		}
		return c.Status(500).JSON(
		&Error{
			Error: err.Error(),
		})
	}

	return c.Status(200).JSON(&structs.IsLogin{
		IsLogin: true,
	})
}

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
		e := err.Error()
		if e == "there is no playing state" {
			return c.Status(200).JSON(
				Message{
					Message: "there is no playing state",
				},
			)
		}
		return c.Status(500).JSON(
			Error{
				Error: e,
			},
		)
	}

	return c.Status(200).JSON(state)
}
