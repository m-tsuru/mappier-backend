package structs

import (
	"time"

	"gorm.io/gorm"
)

type PlaceInfoRaw struct {
	//
	// Yahoo! 場所情報API から出力されます．
	// lat, lon を渡すと，確度が高い順に返答があります．Result のすべての地点は DB にキャッシュせず，登録した場合だけキャッシュする予定です．
	// Req: https://map.yahooapis.jp/placeinfo/V1/get?lat={緯度}&lon={経度}&｛YOLP_API_ID｝&output=json
	// Docs: https://developer.yahoo.co.jp/webapi/map/openlocalplatform/v1/placeinfo.html
	//
	ResultSet struct {
		Address []string `json:"Address"`
		Govcode string   `json:"Govcode"`
		Country struct {
			Code string `json:"Code"`
			Name string `json:"Name"`
		} `json:"Country"`
		Roadname interface{} `json:"Roadname"`
		Result   []struct {
			Name     string  `json:"Name"`
			UID      string  `json:"Uid"`
			Category string  `json:"Category"`
			Label    string  `json:"Label"`
			Where    string  `json:"Where"`
			Combined string  `json:"Combined"`
			Score    float64 `json:"Score"`
		} `json:"Result"`
		Area []struct {
			ID    string  `json:"Id"`
			Name  string  `json:"Name"`
			Score float64 `json:"Score"`
			Type  int     `json:"Type"`
		} `json:"Area"`
	} `json:"ResultSet"`
}

type LocalSearchRaw struct {
	//
	// Yahoo! ローカルサーチAPI から出力されます．
	// uid に 建物の uid を渡すと，返答があります．
	// Req: https://map.yahooapis.jp/search/local/V1/localSearch?appid={YOLP_API_ID}&uid={地点情報UID}&output=json
	// Docs: https://developer.yahoo.co.jp/webapi/map/openlocalplatform/v1/localsearch.html
	//
	ResultInfo struct {
		Count       int     `json:"Count"`
		Total       int     `json:"Total"`
		Start       int     `json:"Start"`
		Status      int     `json:"Status"`
		Description string  `json:"Description"`
		Copyright   string  `json:"Copyright"`
		Latency     float64 `json:"Latency"`
	} `json:"ResultInfo"`
	Feature []struct {
		ID       string `json:"Id"`
		Gid      string `json:"Gid"`
		Name     string `json:"Name"`
		Geometry struct {
			Type        string `json:"Type"`
			Coordinates string `json:"Coordinates"`
		} `json:"Geometry"`
		Category    []string      `json:"Category"`
		Description string        `json:"Description"`
		Style       []interface{} `json:"Style"`
		Property    struct {
			UID        string `json:"Uid"`
			CassetteID string `json:"CassetteId"`
			Yomi       string `json:"Yomi"`
			Country    struct {
				Code string `json:"Code"`
				Name string `json:"Name"`
			} `json:"Country"`
			Address              string `json:"Address"`
			GovernmentCode       string `json:"GovernmentCode"`
			AddressMatchingLevel string `json:"AddressMatchingLevel"`
			LandmarkCode         string `json:"LandmarkCode"`
			Genre                []struct {
				Code string `json:"Code"`
				Name string `json:"Name"`
			} `json:"Genre"`
			Area []struct {
				Code string `json:"Code"`
				Name string `json:"Name"`
			} `json:"Area"`
			Station []struct {
				ID       string `json:"Id"`
				SubID    string `json:"SubId"`
				Name     string `json:"Name"`
				Railway  string `json:"Railway"`
				Exit     string `json:"Exit"`
				ExitID   string `json:"ExitId"`
				Distance string `json:"Distance"`
				Time     string `json:"Time"`
				Geometry struct {
					Type        string `json:"Type"`
					Coordinates string `json:"Coordinates"`
				} `json:"Geometry"`
			} `json:"Station"`
			SmartPhoneCouponFlag string `json:"SmartPhoneCouponFlag"`
			KeepCount            string `json:"KeepCount"`
		} `json:"Property"`
	} `json:"Feature"`
}

// Not in Database
type BuildingAbbr struct {
	ID string `gorm:"primaryKey"`
	Name string
	AreaName string
}

type Building struct {
	ID string `gorm:"primaryKey"`
	Name string
	AreaName string
	Latitude float64
	Longitude float64
	CreatedAt time.Time `json:"omitempty"`
	UpdatedAt time.Time `json:"omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index",json:"omitempty"`
}
