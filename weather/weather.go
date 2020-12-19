package weather

import (
	"github.com/apex/log"
)

const (
	longSleeves  = "long sleeves"
	shortSleeves = "short sleeves"
	umbrella     = "bring an umbrella"
	raincoat     = "bring a rain jacket"
)

// Weather is the data entity we use for our business logic
type Weather struct {
	Temp       float64
	Pressure   float64
	MinTemp    float64
	MaxTemp    float64
	ChanceRain float64
}

// Provider returns a Weather struct for a given set of inputs
type Provider interface {
	GetWeatherByCity(string) (Weather, error)
}

// ForecastService implements the business logic surrounding Weather
type ForecastService struct {
	provider Provider
}

func NewForecastService(p Provider) *ForecastService {
	return &ForecastService{provider: p}
}

func (fs *ForecastService) HowToDress(city string) ([]string, error) {
	w, err := fs.provider.GetWeatherByCity(city)
	if err != nil {
		log.WithField("city", city).Error("unable to retrieve weather data for city")
		return nil, err
	}
	options := make([]string, 0, 3)

	if w.Temp < 21 {
		options = append(options, longSleeves)
	} else {
		options = append(options, shortSleeves)
	}

	if w.ChanceRain >= 0.6 {
		options = append(options, umbrella)
	} else if w.ChanceRain >= 0.2 {
		options = append(options, raincoat)
	}
	return options, nil
}
