package open_meteo_parser

import (
	"testing"
	"time"
)

func TestParser_GetOpenWeatherForecast(t *testing.T) {
	type fields struct {
		APIKey        string
		CloudfrontURL string
	}
	type args struct {
		latitude  float64
		longitude float64
		startTime time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Forecast
		wantErr bool
	}{
		{
			name: "Test GetOpenWeatherForecast",
			fields: fields{
				"xxx",
				"https://ddd.cloudfront.net",
			},
			args: args{
				latitude:  -8.68163896537287,
				longitude: 115.19724863873421,
				startTime: time.Now(),
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.fields.APIKey, tt.fields.CloudfrontURL)
			wc, err := p.GetOpenWeatherForecast(tt.args.latitude, tt.args.longitude, tt.args.startTime)

			aq, err := p.GetOpenWeatherAQI(tt.args.latitude, tt.args.longitude, tt.args.startTime)

			if err != nil {
				t.Error(err)
				t.Fail()
				return
			}

			if wc == nil {
				t.Error("Weather code is nil")
				t.Fail()
				return
			}

			if aq == nil {
				t.Error("AQI is nil")
				t.Fail()
				return
			}
		})
	}
}

func TestPackage(t *testing.T) {

	type fields struct {
		omp IOpenMeteoParser
	}

	type args struct {
		latitude  float64
		longitude float64
		startTime time.Time
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Test GetWeatherForecastByLatLon",
			fields{
				NewParser("xxx", "https://ddd.cloudfront.net"),
			},
			args{
				latitude:  -8.68163896537287,
				longitude: 115.19724863873421,
				startTime: time.Now(),
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.fields.omp

			wc, err := s.GetOpenWeatherForecast(tt.args.latitude, tt.args.longitude, tt.args.startTime)
			if err != nil || wc == nil {
				t.Error(err)
				t.Fail()
				return
			}

			aq, err := s.GetOpenWeatherAQI(tt.args.latitude, tt.args.longitude, tt.args.startTime)

			if err != nil || aq == nil {
				t.Error(err)
				t.Fail()
				return
			}

		})
	}

}
