package types

type StoreData struct {
	Name        string  `json:"name"`
	Region      string  `json:"region"`
	Country     string  `json:"country"`
	Latitude    float32 `json:"latitude"`
	Longitude   float32 `json:"longitude"`
	LocalTime   string  `json:"localtime"`
	TempC       float32 `json:"temp_c"`
	TempF       float32 `json:"temp_f"`
	LastUpdated string  `json:"last_updated"`
	Text        string  `json:"text"`
	Icon        string  `json:"icon"`
	Code        int     `json:"code"`
	Uv          int     `json:"uv"`
}

type Api struct {
	Url    string
	City   string
	ApiKey string
}
