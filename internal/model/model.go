package model

type Weather struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float32 `json:"lat"`
		Lon            float32 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocalTimeEpoch uint32  `json:"localtime_epoch"`
		LocalTime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch uint32  `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float32 `json:"temp_c"`
		TempF            float32 `json:"temp_f"`
		IsDay            int     `json:"is_day"`
		Uv               int     `json:"-"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
	} `json:"current"`
}
