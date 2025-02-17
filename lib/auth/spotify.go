package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/m-tsuru/mappier-backend/lib/structs"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
	"gopkg.in/ini.v1"
	"gorm.io/gorm"
)

func RefreshAccessToken(refreshToken string) (*string, *int, error) {
	fmt.Println("refresh access token")
	raw, err := ini.Load("config.ini")
	if err != nil {
		return nil, nil, err
	}

	u := "https://accounts.spotify.com/api/token"

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("", raw.Section("spotify").Key("CLIENT_ID").String())

	basic := raw.Section("spotify").Key("CLIENT_ID").String() + ":" + raw.Section("spotify").Key("CLIENT_SECRET").String()

    req, err := http.NewRequest("POST", u, bytes.NewBufferString(data.Encode()))
    if err != nil {
        return nil, nil, fmt.Errorf("error creating request: %w", err)
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", basic))

	c := &http.Client{
		Timeout: 3 * time.Second,
	}

	res, err := c.Do(req)
    if err != nil {
        return nil, nil, fmt.Errorf("error sending request: %w", err)
    }
    defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get access token (refreshed) from spotify: %w", err)
	}
	fmt.Println(body)

	var result structs.SpotifyRefreshingToken

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get access token (refreshed) from spotify: %w", err)
	}
	return &result.AccessToken, &result.ExpiresIn, nil
}

func NewSpotifyAuthConf() (*oauth2.Config, error) {
	raw, err := ini.Load("config.ini")
	if err != nil {
		return nil, err
	}
	conf := oauth2.Config{
		ClientID: raw.Section("spotify").Key("CLIENT_ID").String(),
		ClientSecret: raw.Section("spotify").Key("CLIENT_SECRET").String(),
		Scopes: []string{"user-read-private", "user-read-email", "user-read-currently-playing", "user-read-playback-state"},
		Endpoint: oauth2.Endpoint{
			AuthURL: spotify.Endpoint.AuthURL,
			TokenURL: spotify.Endpoint.TokenURL,
		},
		RedirectURL: raw.Section("spotify").Key("REDIRECT_URI").String(),
	}
	return &conf, nil
}

func GetSpotifyRedirectUrl() (*string, error) {
	conf, err := NewSpotifyAuthConf()
	state := "listen"
	if err != nil {
		return nil, fmt.Errorf("get spotify redirect url: %w", err)
	}
	redirectURL := conf.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	return &redirectURL, nil
}

func GetSpotifyAccessToken(state string, code string) (*oauth2.Token, error) {
	if state != "listen" {
		return nil, fmt.Errorf("state is not valid: %v", state)
	} else if code == "" {
		return nil, fmt.Errorf("required param code is nil")
	}

	ctx := context.Background()
	conf, e := NewSpotifyAuthConf()
	if e != nil {
		return nil, fmt.Errorf("can't get spotify configuration")
	}

	token, err := conf.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func GetSpotifyUser(token oauth2.Token) (*structs.SpotifyUserRaw, error) {
	url := "https://api.spotify.com/v1/me"
	accessToken := token.AccessToken

	c := &http.Client{
		Timeout: 3 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get user profile from spotify: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to get user profile from spotify: %w", err)
	}

	defer res.Body.Close()

	if res.Status != "200 OK" {
		return nil, fmt.Errorf("unable to get user profile from spotify (due to api): %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to get user profile from spotify: %w", err)
	}

	var user structs.SpotifyUserRaw
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal user profile from spotify: %w", err)
	}
	return &user, nil
}

func GetSpotifyPlayingState(db *gorm.DB, accessToken string) (*structs.SpotifySongCache, error) {
	url := "https://api.spotify.com/v1/me/player/currently-playing"

	c := &http.Client{
		Timeout: 3 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get user playing state from spotify: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	res, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to get user playing state from spotify: %w", err)
	}

	defer res.Body.Close()

	if res.Status == "204 No Content" {
		return nil, fmt.Errorf("there is no playing state")
	}

	if res.Status != "200 OK" {
		return nil, fmt.Errorf("unable to get user playing state from spotify (due to api): %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to get user playing state from spotify: %w", err)
	}

	var state structs.SpotifyPlayingRaw
	err = json.Unmarshal(body, &state)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal user playing state from spotify: %w", err)
	}

	var artistNames []string
    for _, artist := range state.Item.Artists {
        artistNames = append(artistNames, artist.Name)
    }
    artists := strings.Join(artistNames, ", ")

	response := &structs.SpotifySongCache{
		ID: state.Item.ID,
		Image: state.Item.Album.Images[0].URL,
		Name: state.Item.Name,
		ArtistsPureString: artists,
		Album: state.Item.Album.Name,
	}

	var chk structs.SpotifySongCache
	err = db.Where("id = ?", state.Item.ID).First(&chk).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		result := db.Create(&response)
		if result.Error != nil {
			log.Print(result.Error)
		}
	}
	return response, nil
}
