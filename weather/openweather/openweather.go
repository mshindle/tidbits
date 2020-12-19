package openweather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/apex/log"
	"gitlab.com/mshindle/tidbits/weather"
)

const (
	endpoint   = "https://api.openweathermap.org/data/2.5"
	pathByCity = "/weather?q=%s&appid=%s&units=metric"
)

type client struct {
	apikey string
}

type apiResponse struct {
	message string
	main    struct {
		Temp     float64 `json:"temp"`
		Pressure float64 `json:"pressure"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
	}
}

func (a apiResponse) toWeather() weather.Weather {
	return weather.Weather{
		Temp:     a.main.Temp,
		Pressure: a.main.Pressure,
		MaxTemp:  a.main.TempMax,
		MinTemp:  a.main.TempMin,
	}
}

func New(apikey string) *client {
	return &client{apikey: apikey}
}

func (c *client) GetWeatherByCity(city string) (weather.Weather, error) {
	u := endpoint + fmt.Sprintf(pathByCity, city, c.apikey)
	res, err := http.Get(u)
	if err != nil {
		log.WithError(err).Error("failed http get from openweather")
		return weather.Weather{}, weather.ErrProviderFailure
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"code":   res.StatusCode,
			"status": res.Status,
		}).Error("invalid response code from openweather")
		return weather.Weather{}, weather.ErrProviderFailure
	}

	// read the response body and encode it into the response struct
	bodyRaw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.WithError(err).Error("failed http response from openweather")
		return weather.Weather{}, weather.ErrProviderFailure
	}

	var apiRes apiResponse
	if err = json.Unmarshal(bodyRaw, &apiRes); err != nil {
		log.WithError(err).Error("failed to parse http response")
		return weather.Weather{}, weather.ErrProviderFailure
	}

	// return the external response converted into an entity
	return apiRes.toWeather(), nil
}
