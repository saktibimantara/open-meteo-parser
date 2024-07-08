package open_meteo_parser

import (
	"fmt"
	pom "github.com/saktibimantara/go-open-meteo"
	"time"
)

type Parser struct {
	APIKey        string
	CloudfrontURL string
	om            pom.IGoOpenMeteo
}

func NewParser(apiKey, cloudfrontURL string) *Parser {

	om := pom.New(pom.NewConfig())

	return &Parser{
		APIKey:        apiKey,
		CloudfrontURL: cloudfrontURL,
		om:            om,
	}
}

type IParser interface {
	GetOpenWeatherForecast(latitude, longitude float64, startTime time.Time) (*Forecast, error)
	GetOpenWeatherAQI(latitude, longitude float64, startTime time.Time)
}

func (p Parser) GetOpenWeatherForecast(latitude, longitude float64, startTime time.Time) (*Forecast, error) {

	weatherForecast, err := p.getWeatherWithOpenWeatherFormat(latitude, longitude, startTime)
	if err != nil {
		return nil, err
	}

	return weatherForecast, nil
}

func generateAQIParam(lat, lon float64) *pom.AQIParams {
	params, err := pom.NewAQIParamsBuilder().
		SetLatitude(-8.68163896537287).
		SetLongitude(115.19724863873421).
		SetForecastDays(5).
		AddHourlyParam(pom.PM10, pom.PM2_5, pom.PM2_5, pom.CarbonMonoxide, pom.NitrogenDioxide, pom.SulphurDioxide, pom.Ozone, pom.UVIndex, pom.USAQI).
		Build()

	if err != nil {
		return nil
	}

	return &params
}

func GenerateParams(lat, lon float64) *pom.ForecastParams {
	params, err := pom.NewForecastParamsBuilder().
		SetLatitude(lat).
		SetLongitude(lon).
		SetForecastDays(12).
		AddHourlyParam(
			pom.Temperature2m,
			pom.WindSpeed10m,
			pom.Precipitation,
			pom.Rain,
			pom.WeatherCode,
			pom.WindSpeed10m,
			pom.WindDirection10m,
			pom.WindGusts10m,
			pom.SurfacePressure,
			pom.PressureMSL,
			pom.IsDay,
		).
		AddMinutely15Param(
			pom.Minutely15Temperature2m,
			pom.Minutely15Precipitation,
			pom.Minutely15Rain,
			pom.Minutely15WeatherCode,
			pom.Minutely15RelativeHumidity2m,
			pom.Minutely15WindDirection10m,
			pom.Minutely15WindSpeed10m,
			pom.Minutely15WindGusts10m,
			pom.Minutely15ApparentTemperature,
		).
		AddDailyParam(
			pom.DailyTemperature2mMax,
			pom.DailyTemperature2mMin,
			pom.DailyWeatherCode,
		).
		Build()

	if err != nil {
		return nil
	}

	return &params
}

func (p Parser) GetOpenWeatherAQI(latitude, longitude float64, startTime time.Time) (*AQI, error) {

	aqi, err := p.getAQIWithOpenWeatherFormat(latitude, longitude, startTime)
	if err != nil {
		return nil, err
	}

	return aqi, nil
}

func (p Parser) getAQIWithOpenWeatherFormat(lat, lon float64, startTime time.Time) (*AQI, error) {

	aqi, err := p.om.GetAQI(generateAQIParam(lat, lon))
	if err != nil {
		return nil, err
	}

	wd := pom.NewWeatherData().SetAQIForecastResponse(aqi)

	wp := pom.NewWeatherProcessor(wd)

	nf, err := wp.FindNearestAQIForecastByTime(startTime)
	if err != nil {
		return nil, err
	}

	if nf == nil {
		return nil, fmt.Errorf("forecast is nil")
	}

	return ParseToAQI(*nf.AqiHourlyForecast), err
}

func (p Parser) getWeatherWithOpenWeatherFormat(lat, lon float64, startTime time.Time) (*Forecast, error) {

	openResp, err := p.om.Forecast(GenerateParams(lat, lon))
	if err != nil {
		return nil, err
	}

	wd := pom.NewWeatherData().SetForecastResponse(openResp)

	wp := pom.NewWeatherProcessor(wd)

	nf, err := wp.FindNearestForecastByTime(startTime)
	if err != nil {
		return nil, err
	}

	if nf == nil {
		return nil, fmt.Errorf("forecast is nil")
	}

	return ParseToForecast(*nf), err
}

func ParseToAQI(aqi pom.NearestAQIHourlyForecast) *AQI {

	aqiData := NewAQIBuilder().
		SetUsAqi(safeFloat64(aqi.USAQI)).
		SetCo(safeFloat64(aqi.CarbonMonoxide)).
		SetNo2(safeFloat64(aqi.NitrogenDioxide)).
		SetO3(safeFloat64(aqi.Ozone)).
		SetPm10(safeFloat64(aqi.PM10)).
		SetPm2_5(safeFloat64(aqi.PM2_5)).
		SetSo2(safeFloat64(aqi.SulphurDioxide)).
		SetDt(int(safeDate(aqi.Time).Unix())).
		Build()

	return aqiData
}

func ParseToForecast(forecast pom.NearestForecast) *Forecast {
	var temp *float64
	var dt *time.Time
	var feelsLike *float64
	var tempMax *float64
	var tempMin *float64
	var pressure *int
	var seaLevel *int
	var grndLevel *int
	var tempKf *float64
	var weather *Weather
	var humidity *int
	var windSpeed *float64
	var windDeg *int
	var windGust *float64
	var rain *float64
	var isDay *int

	if forecast.Minutely15Forecast != nil {
		dt = &forecast.Minutely15Forecast.Time.Time
		temp = forecast.Minutely15Forecast.Temperature2m
		feelsLike = forecast.Minutely15Forecast.ApparentTemperature
		if weatherCode := forecast.Minutely15Forecast.WeatherCode; weatherCode != nil {
			weather = ParseWeatherCode(*weatherCode)
		}

		if hm := forecast.Minutely15Forecast.RelativeHumidity2m; hm != nil {
			humidity = new(int)
			*humidity = int(*hm)
		}

		if ws := forecast.Minutely15Forecast.WindSpeed10m; ws != nil {
			windSpeed = new(float64)
			*windSpeed = *ws
		}

		if wd := forecast.Minutely15Forecast.WindDirection10m; wd != nil {
			windDeg = new(int)
			*windDeg = int(*wd)
		}

		if wg := forecast.Minutely15Forecast.WindGusts10m; wg != nil {
			windGust = new(float64)
			*windGust = *wg
		}

	}

	if forecast.HourlyForecast != nil {
		if dt == nil {
			dt = &forecast.HourlyForecast.Time.Time
		}

		if temp == nil {
			temp = forecast.HourlyForecast.Temperature2m
		}

		if weather == nil {
			if weatherCode := forecast.HourlyForecast.WeatherCode; weatherCode != nil {
				weather = ParseWeatherCode(*weatherCode)
			}
		}

		if hm := forecast.HourlyForecast.RelativeHumidity2m; hm != nil && humidity == nil {
			humidity = new(int)
			*humidity = int(*hm)
		}

		if ws := forecast.HourlyForecast.WindSpeed10m; ws != nil && windSpeed == nil {
			windSpeed = new(float64)
			*windSpeed = *ws
		}

		if wd := forecast.HourlyForecast.WindDirection10m; wd != nil && windDeg == nil {
			windDeg = new(int)
			*windDeg = int(*wd)
		}

		if wg := forecast.HourlyForecast.WindGusts10m; wg != nil && windGust == nil {
			windGust = new(float64)
			*windGust = *wg
		}

		if ps := forecast.HourlyForecast.PressureMSL; ps != nil {
			pressure = new(int)
			*pressure = int(*ps)
		}

		if rn := forecast.HourlyForecast.Rain; rn != nil {
			rain = new(float64)
			*rain = *rn
		}

		if id := forecast.HourlyForecast.IsDay; id != nil {
			isDay = new(int)
			*isDay = *id
		}

	}

	if forecast.DailyForecast != nil {
		if dt == nil {
			dt = &forecast.DailyForecast.Time.Time
		}

		if weather == nil {
			if weatherCode := forecast.DailyForecast.WeatherCode; weatherCode != nil {
				weather = ParseWeatherCode(*weatherCode)
			}
		}

		if tmpMax := forecast.DailyForecast.Temperature2mMax; tmpMax != nil {
			tempMax = new(float64)
			*tempMax = *tmpMax
		}

		if tmpMin := forecast.DailyForecast.Temperature2mMin; tmpMin != nil {
			tempMin = new(float64)
			*tempMin = *tmpMin
		}

	}

	return &Forecast{
		Dt: int(safeDate(dt).Unix()),
		Main: Main{
			Temp:      safeFloat64(temp),
			FeelsLike: safeFloat64(feelsLike),
			TempMax:   safeFloat64(tempMax),
			TempMin:   safeFloat64(tempMin),
			Pressure:  safeInt(pressure),
			SeaLevel:  safeInt(seaLevel),
			GrndLevel: safeInt(grndLevel),
			TempKf:    safeFloat64(tempKf),
			Humidity:  safeInt(humidity),
		},
		Weather: []Weather{
			safeWeather(weather, isDay),
		},
		Clouds: Clouds{
			All: 0,
		},
		Visibility: 0,
		Pop:        0,
		Wind: Wind{
			Speed: safeFloat64(windSpeed),
			Deg:   safeInt(windDeg),
		},
		Sys: Sys{},
		Rain: Rain{
			ThreeH: safeFloat64(rain),
		},
		DtTxt: safeDate(dt).String(),
	}

}

func safeWeather(w *Weather, isDay *int) Weather {
	if w == nil {
		return Weather{
			ID:          800,
			Main:        "Clear",
			Description: "clear sky",
			Icon:        "",
		}
	}

	if id := w.ID; id != 0 {
		w.Icon = "01n"
		return *w
	}

	return *w
}

func ParseWeatherCode(weatherCode pom.WeatherCodeResponse) *Weather {

	var openWeatherCode int
	var main string
	var description string

	switch weatherCode {
	case pom.WeatherCodeClearSky:
		openWeatherCode = 800
		main = "Clear"
		description = "clear sky"
	case pom.WeatherCodeMainlyClear:
		openWeatherCode = 801
		main = "Clouds"
		description = "few clouds"
	case pom.WeatherCodePartlyCloudy:
		openWeatherCode = 802
		main = "Clouds"
		description = "scattered clouds"
	case pom.WeatherCodeOvercast:
		openWeatherCode = 804
		main = "Clouds"
		description = "overcast clouds"
	case pom.WeatherCodeFog:
		openWeatherCode = 741
		main = "Fog"
		description = "fog"
	case pom.WeatherCodeDepositingRimeFog:
		openWeatherCode = 741
		main = "Fog"
		description = "fog"
	case pom.WeatherCodeLightDrizzle:
		openWeatherCode = 300
		main = "Drizzle"
		description = "light intensity drizzle"
	case pom.WeatherCodeModerateDrizzle:
		openWeatherCode = 301
		main = "Drizzle"
		description = "drizzle"
	case pom.WeatherCodeDenseDrizzle:
		openWeatherCode = 302
		main = "Drizzle"
		description = "heavy intensity drizzle"
	case pom.WeatherCodeLightFreezingDrizzle:
		openWeatherCode = 310
		main = "Drizzle"
		description = "light intensity drizzle rain"
	case pom.WeatherCodeDenseFreezingDrizzle:
		openWeatherCode = 313
		main = "Drizzle"
		description = "shower rain and drizzle"
	case pom.WeatherCodeSlightRain:
		openWeatherCode = 500
		main = "Rain"
		description = "light rain"
	case pom.WeatherCodeModerateRain:
		openWeatherCode = 501
		main = "Rain"
		description = "moderate rain"
	case pom.WeatherCodeHeavyRain:
		openWeatherCode = 502
		main = "Rain"
		description = "heavy intensity rain"
	case pom.WeatherCodeLightFreezingRain:
		openWeatherCode = 511
		main = "Rain"
		description = "freezing rain"
	case pom.WeatherCodeHeavyFreezingRain:
		openWeatherCode = 511
		main = "Rain"
	case pom.WeatherCodeSlightSnowFall:
		openWeatherCode = 600
		main = "Snow"
		description = "light snow"
	case pom.WeatherCodeModerateSnowFall:
		openWeatherCode = 601
		main = "Snow"
		description = "snow"
	case pom.WeatherCodeHeavySnowFall:
		openWeatherCode = 602
		main = "Snow"
		description = "heavy snow"
	case pom.WeatherCodeSnowGrains:
		openWeatherCode = 611
		main = "Snow"
		description = "sleet"
	case pom.WeatherCodeSlightRainShowers:
		openWeatherCode = 520
		main = "Rain"
		description = "light intensity shower rain"
	case pom.WeatherCodeModerateRainShowers:
		openWeatherCode = 521
		main = "Rain"
		description = "shower rain"
	case pom.WeatherCodeViolentRainShowers:
		openWeatherCode = 522
		main = "Rain"
		description = "heavy intensity shower rain"
	case pom.WeatherCodeSlightSnowShowers:
		openWeatherCode = 620
		main = "Snow"
		description = "light shower snow"
	case pom.WeatherCodeHeavySnowShowers:
		openWeatherCode = 622
		main = "Snow"
		description = "heavy shower snow"
	case pom.WeatherCodeThunderstorm:
		openWeatherCode = 200
		main = "Thunderstorm"
		description = "thunderstorm with light rain"
	case pom.WeatherCodeSlightHailThunder:
		openWeatherCode = 200
		main = "Thunderstorm"
		description = "thunderstorm with light rain"
	case pom.WeatherCodeHeavyHailThunder:
		openWeatherCode = 200
		main = "Thunderstorm"
		description = "thunderstorm with light rain"
	default:
		openWeatherCode = 800
		main = "Clear"
		description = "clear sky"
	}

	return &Weather{
		ID:          openWeatherCode,
		Main:        main,
		Description: description,
	}

}

func safeDate(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}

	return *t
}

func safeFloat64(f *float64) float64 {
	if f == nil {
		return 0
	}

	return *f
}

func safeInt(i *int) int {
	if i == nil {
		return 0
	}

	return *i
}
