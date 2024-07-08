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

			t.Log(wc)
			t.Log(aq)

			if err != nil {
				t.Error(err)
				t.Fail()
				return
			}
		})
	}
}
