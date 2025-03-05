package main

import (
	"fmt"
	"math"
	"net/http"

	owm "github.com/briandowns/openweathermap"
)

var cfg *Config

func SetConfig(c *Config) {
	cfg = c
}

func GetCurrentTemperature() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		weather, err := owm.NewCurrent(cfg.Unit, cfg.Lang, cfg.OWMAPIKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = weather.CurrentByCoordinates(
			&owm.Coordinates{
				Longitude: cfg.Longitude,
				Latitude:  cfg.Latitude,
			})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		temp := int(math.Round(weather.Main.Temp))
		w.Write([]byte(fmt.Sprintf("%d\n", temp)))
	}
}
