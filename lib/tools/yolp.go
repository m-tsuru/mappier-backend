package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/m-tsuru/mappier-backend/lib/structs"
	"gopkg.in/ini.v1"
)

type YOLP struct {
	CLIENT_ID string
}

var yolp YOLP

func Init() error {
	conf, err := ini.Load("config.ini")
	if err != nil {
		return err
	}
	yolp.CLIENT_ID = conf.Section("yolp").Key("CLIENT_ID").String()
	return nil
}

func GetBuildingFromLatLon(latitude float64, longitude float64) (*[]structs.BuildingAbbr, error) {
	conf, err := ini.Load("config.ini")
	if err != nil {
		return nil, err
	}
	yolp.CLIENT_ID = conf.Section("yolp").Key("CLIENT_ID").String()

	latitudeString := strconv.FormatFloat(latitude, 'f', -1, 64)
	longitudeString := strconv.FormatFloat(longitude, 'f', -1, 64)
	url := fmt.Sprintf("https://map.yahooapis.jp/placeinfo/V1/get?lat=%s&lon=%s&appid=%s&output=json", latitudeString, longitudeString, yolp.CLIENT_ID)

	c := &http.Client{
		Timeout: 3 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.Status != "200 OK" {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var raw structs.PlaceInfoRaw
	err = json.Unmarshal(body, &raw)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal location list from yolp: %w", err)
	}

	var listing []structs.BuildingAbbr
	for i, v := range raw.ResultSet.Result {
		listing = append(listing, structs.BuildingAbbr{
			ID:       v.UID,
			Name:     v.Label,
			AreaName: v.Where,
		})
		if i == 2 {
			break
		}
	}

	return &listing, nil
}

func GetBuildingOfLatLon(buildingList []structs.BuildingAbbr) (*[]structs.Building, error) {
	conf, err := ini.Load("config.ini")
	if err != nil {
		return nil, err
	}

	yolp.CLIENT_ID = conf.Section("yolp").Key("CLIENT_ID").String()

	var listing []structs.Building

	for _, v := range buildingList {
		url := fmt.Sprintf("https://map.yahooapis.jp/search/local/V1/localSearch?appid=%s&uid=%s&output=json", yolp.CLIENT_ID, v.ID)
		c := &http.Client{
			Timeout: 3 * time.Second,
		}

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		res, err := c.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.Status != "200 OK" {
			return nil, err
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		var raw structs.LocalSearchRaw
		err = json.Unmarshal(body, &raw)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal location lat/lon [%s] from yolp: %w", v.ID, err)
		}

		coords := strings.SplitN(raw.Feature[0].Geometry.Coordinates, ",", 2)
		if len(coords) != 2 {
			return nil, fmt.Errorf("invalid coordinates format for uid [%s]", v.ID)
		}
		longitude, err := strconv.ParseFloat(coords[0], 64)
		if err != nil {
			return nil, fmt.Errorf("unable to parse longitude for uid [%s]: %w", v.ID, err)
		}
		latitude, err := strconv.ParseFloat(coords[1], 64)
		if err != nil {
			return nil, fmt.Errorf("unable to parse latitude for uid [%s]: %w", v.ID, err)
		}

		listing = append(listing, structs.Building{
			ID:        v.ID,
			Name:      v.Name,
			AreaName:  v.AreaName,
			Latitude:  latitude,
			Longitude: longitude,
		})

	}
	return &listing, nil
}
