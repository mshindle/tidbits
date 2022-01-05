package weather

import (
	"reflect"
	"testing"
)

func TestForecastService_HowToDress(t *testing.T) {
	type fields struct {
		provider Provider
	}
	type args struct {
		city string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "rainy cold",
			fields: fields{
				provider: mockProvider{Temp: 20.0, ChanceRain: 0.8},
			},
			args:    args{city: "rainy-cold"},
			want:    []string{longSleeves, umbrella},
			wantErr: false,
		},
		{
			name: "bad provider",
			fields: fields{
				provider: mockProvider{Temp: 20.0, ChanceRain: 0.8},
			},
			args:    args{city: "error"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sunny clear",
			fields: fields{
				provider: mockProvider{Temp: 30.0, ChanceRain: 0.1},
			},
			args:    args{city: "sunnyclear"},
			want:    []string{shortSleeves},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &ForecastService{
				provider: tt.fields.provider,
			}
			got, err := fs.HowToDress(tt.args.city)
			if (err != nil) != tt.wantErr {
				t.Errorf("HowToDress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HowToDress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockProvider Weather

func (m mockProvider) GetWeatherByCity(city string) (Weather, error) {
	if city == "error" {
		return Weather{}, ErrProviderFailure
	}

	return Weather(m), nil
}
