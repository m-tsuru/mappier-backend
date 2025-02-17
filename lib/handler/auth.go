package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/m-tsuru/mappier-backend/lib/auth"
	"github.com/m-tsuru/mappier-backend/lib/db"
)

func HandleAuth(c *fiber.Ctx) error {
	// @Summary Login with Spotify
	// @Description Login with Spotify
	// @Tags login
	// @Accept json
	// @Produce json
	// @Router /auth/login [get]
	// @Success 200
	// @Header 200 {string} location
    url, err := auth.GetSpotifyRedirectUrl()
    if err != nil {
        e := err.Error()
        return c.Status(400).JSON(
            Error{
                Error: e,
            },
        )
    }
    c.Status(fiber.StatusSeeOther)
    c.Redirect(*url)

    return c.JSON(Url{Url: url})
}

func HandleAuthCallback(c *fiber.Ctx) error {
	//@Summary Redirect from Spotify
	//@Description Redirect from Spotify
	//@Tags login
	//@Accept json
	//@Produce json
	//@Router /auth/callback [get]
	//@Success 200
	//@Header 200 {string} location

	d, err := db.Open()
	if err != nil {
		return err
	}

	m := c.Queries()
	res, err := auth.GetSpotifyAccessToken(m["state"], m["code"])
	if err != nil {
		e := err.Error()
		return c.Status(500).JSON(
			Error{
				Error: e,
			},
		)
	}

	user, err := auth.GetSpotifyUser(*res)
	if err != nil {
		e := err.Error()
		return c.Status(500).JSON(
			Error{
				Error: e,
			},
		)
	}

	err = db.SaveSpotifyAccessToken(d, res, user.ID)
	if err != nil {
		e := err.Error()
		return c.Status(500).JSON(
			Error{
				Error: e,
			},
		)
	}

	err = db.SaveSpotifyUser(d, *user)
	if err != nil {
		e := err.Error()
		return c.Status(500).JSON(
			Error{
				Error: e,
			},
		)
	}

	cookie := (&fiber.Cookie{
		Name: "SessionId",
		Value: utils.UUIDv4(),
		Path: "/",
		Expires: time.Now().Add(72 * time.Hour),
		MaxAge: int((72 * time.Hour).Seconds()),
		SameSite: "Lax",
		Domain: "localhost",
		HTTPOnly: true,
	})

	err = db.SaveSessionId(d, user.ID, cookie.Value, cookie.Expires)
	if err != nil {
		e := err.Error()
		return c.Status(500).JSON(
			Error{
				Error: e,
			},
		)
	}

	c.Cookie(cookie)

	return c.Status(200).JSON(
		Message{
			Message: "Logged In Successfully",
		},
	)
}
