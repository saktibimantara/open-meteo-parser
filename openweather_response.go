package open_meteo_parser

import (
	"encoding/json"
	"math"
	"strings"
	"time"
)

type Response3HoursStepForecast struct {
	Cod     string     `json:"cod"`
	Message int        `json:"message"`
	Cnt     int        `json:"cnt"`
	List    []Forecast `json:"list"`
}

type ResponseAQI struct {
	List []AQI `json:"list"`
}

type AQI struct {
	Main struct {
		Aqi int `json:"aqi"`
	}
	Components struct {
		Co    float64 `json:"co"`
		No    float64 `json:"no"`
		No2   float64 `json:"no2"`
		O3    float64 `json:"o3"`
		So2   float64 `json:"so2"`
		Pm2_5 float64 `json:"pm2_5"`
		Pm10  float64 `json:"pm10"`
		Nh3   float64 `json:"nh3"`
	}
	Dt int `json:"dt"`
}

type AQIBuilder struct {
	AQI *AQI
}

func NewAQIBuilder() *AQIBuilder {
	return &AQIBuilder{
		AQI: &AQI{},
	}
}

func (b *AQIBuilder) SetAqi(aqi int) *AQIBuilder {
	b.AQI.Main.Aqi = aqi
	return b
}

func (b *AQIBuilder) SetCo(co float64) *AQIBuilder {
	b.AQI.Components.Co = co
	return b
}

func (b *AQIBuilder) SetNo(no float64) *AQIBuilder {
	b.AQI.Components.No = no
	return b
}

func (b *AQIBuilder) SetNo2(no2 float64) *AQIBuilder {
	b.AQI.Components.No2 = no2
	return b
}

func (b *AQIBuilder) SetO3(o3 float64) *AQIBuilder {
	b.AQI.Components.O3 = o3
	return b
}

func (b *AQIBuilder) SetSo2(so2 float64) *AQIBuilder {
	b.AQI.Components.So2 = so2
	return b
}

func (b *AQIBuilder) SetPm2_5(pm2_5 float64) *AQIBuilder {
	b.AQI.Components.Pm2_5 = pm2_5
	return b
}

func (b *AQIBuilder) SetPm10(pm10 float64) *AQIBuilder {
	b.AQI.Components.Pm10 = pm10
	return b
}

func (b *AQIBuilder) SetNh3(nh3 float64) *AQIBuilder {
	b.AQI.Components.Nh3 = nh3
	return b
}

func (b *AQIBuilder) SetDt(dt int) *AQIBuilder {
	b.AQI.Dt = dt
	return b
}

func (b *AQIBuilder) SetUsAqi(usAQI float64) *AQIBuilder {

	b.AQI.Main.Aqi = int(usAQI)

	if usAQI == 0 {
		b.AQI.Main.Aqi = CalculateAQI(b.AQI.Components.Pm2_5, b.AQI.Components.Pm10, b.AQI.Components.O3, b.AQI.Components.No2, b.AQI.Components.So2, b.AQI.Components.Co)
	}

	return b
}

func (b *AQIBuilder) Build() *AQI {
	return b.AQI
}

func (a *AQI) GetDate() *time.Time {
	t := time.Unix(int64(a.Dt), 0)
	return &t
}

func (a *AQI) MarshalJSON() ([]byte, error) {
	type Alias AQI

	// Calculate AQI
	aqi := CalculateAQI(a.Components.Pm2_5, a.Components.Pm10, a.Components.O3, a.Components.No2, a.Components.So2, a.Components.Co)

	// Override the AQI field
	a.Main.Aqi = aqi

	return json.Marshal((*Alias)(a))
}

func (r *ResponseAQI) FindNearestAQIBasedOnStartTime(time *time.Time) *AQI {
	var nearestAQI AQI
	minDiff := math.MaxInt64

	for _, aqi := range r.List {
		aqiDate := aqi.GetDate()
		diff := int(math.Abs(float64(time.Unix() - aqiDate.Unix())))
		if diff < minDiff {
			minDiff = diff
			nearestAQI = aqi
		}
	}

	return &nearestAQI
}

type (
	Forecast struct {
		Dt         int       `json:"dt"`
		Main       Main      `json:"main"`
		Weather    []Weather `json:"weather"`
		Clouds     Clouds    `json:"clouds"`
		Visibility int       `json:"visibility"`
		Pop        float64   `json:"pop"`
		Wind       Wind      `json:"wind"`
		Sys        Sys       `json:"sys"`
		Rain       Rain      `json:"rain"`
		DtTxt      string    `json:"dt_txt"`
	}

	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		SeaLevel  int     `json:"sea_level"`
		GrndLevel int     `json:"grnd_level"`
		Humidity  int     `json:"humidity"`
		TempKf    float64 `json:"temp_kf"`
		Location  string  `json:"location"`
		Lat       string  `json:"lat"`
		Lng       string  `json:"lng"`
	}

	Weather struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}

	Clouds struct {
		All int `json:"all"`
	}

	Wind struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	}

	Sys struct {
		Pod string `json:"pod"`
	}

	Rain struct {
		ThreeH float64 `json:"3h"`
	}

	Snow struct {
		ThreeH float64 `json:"3h"`
	}

	City struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Coord struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		}
		Country    string `json:"country"`
		Timezone   int    `json:"timezone"`
		Sunrise    int    `json:"sunrise"`
		Sunset     int    `json:"sunset"`
		Population int    `json:"population"`
	}

	Error struct {
		Cod     string `json:"cod"`
		Message string `json:"message"`
	}
)

func (w *Weather) MarshalJSON() ([]byte, error) {
	type Alias Weather

	pod := extractChar(w.Icon)

	// Override the icon field
	switch w.ID {
	case 800:
		w.Icon = "113"
	case 801, 802, 803:
		w.Icon = "116"
	case 804:
		w.Icon = "119"
	case 701:
		w.Icon = "143"
	case 500:
		w.Icon = "176"
	case 600:
		w.Icon = "179"
	case 300, 321:
		w.Icon = "263"
	case 301:
		w.Icon = "266"
	case 313, 520:
		w.Icon = "293"
	case 302, 310, 311, 312:
		w.Icon = "296"
	case 314, 521:
		w.Icon = "299"
	case 501:
		w.Icon = "302"
	case 502:
		w.Icon = "308"
	case 611:
		w.Icon = "317"
	case 602:
		w.Icon = "320"
	case 601:
		w.Icon = "332"
	case 511:
		w.Icon = "350"
	case 522:
		w.Icon = "356"
	case 503, 504, 531:
		w.Icon = "359"
	case 612:
		w.Icon = "362"
	case 613:
		w.Icon = "365"
	case 620:
		w.Icon = "368"
	case 621, 622:
		w.Icon = "371"
	case 200, 210, 230, 231:
		w.Icon = "386"
	case 201, 202, 211, 212, 221, 232:
		w.Icon = "389"
	case 615:
		w.Icon = "615"
	case 616:
		w.Icon = "616"
	case 711:
		w.Icon = "701"
	case 721, 731:
		w.Icon = "731"
	case 741:
		w.Icon = "741"
	case 751:
		w.Icon = "751"
	case 761:
		w.Icon = "761"
	default:
		w.Icon = "800"
	}

	w.Icon = "https://d1c40hpuz0tre6.cloudfront.net/weathers/" + w.Icon + pod + ".png"

	return json.Marshal((*Alias)(w))
}

func extractChar(input string) string {
	// Check if the string contains "n"
	if strings.Contains(input, "n") {
		return "n"
	}

	return "d"
}

func (f *Forecast) GetDate() *time.Time {
	t := time.Unix(int64(f.Dt), 0)
	return &t
}

func (r *Response3HoursStepForecast) FindNearestWeatherBasedOnStartTime(time *time.Time) *Forecast {
	var nearestForecast Forecast

	// case
	// start Time at 06:00:00
	// forecast at 03:00:00, 06:00:00, 09:00:00
	// nearest forecast is 06:00:00

	// Find the nearest forecast based on the start time
	minDiff := math.MaxInt64
	for _, forecast := range r.List {
		forecastDate := forecast.GetDate()
		if forecastDate.Day() == time.Day() {
			diff := int(math.Abs(float64(time.Unix() - forecastDate.Unix())))
			if diff < minDiff {
				minDiff = diff
				nearestForecast = forecast
			}
		}
	}

	return &nearestForecast
}
