package temperature

import (
	"fmt"
	"math"
	"net/http"

	"github.com/ezebunandu/oteller/pkg/config"

	owm "github.com/briandowns/openweathermap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const name = "temperature.service"

var (
	cfg           *config.Config
	meter         = otel.Meter(name)
	tempHistogram metric.Float64Histogram
)

func init() {
	var err error
	tempHistogram, err = meter.Float64Histogram(
		"temperature.celsius",
		metric.WithDescription("Distribution of temperature measurements"),
		metric.WithUnit("°C"),
	)
	if err != nil {
		panic(err)
	}
}

func SetConfig(c *config.Config) {
	cfg = c
}

func GetCurrentTemperature() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tracer := otel.Tracer("get_current_temp")
		ctx, span := tracer.Start(ctx, "fetch_temperature")
		defer span.End()

		weather, err := owm.NewCurrent(cfg.Unit, cfg.Lang, cfg.OWMAPIKey)
		if err != nil {
			span.SetAttributes(attribute.Int("http.status_code", http.StatusInternalServerError))
			span.SetAttributes(attribute.String("error", err.Error()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = weather.CurrentByCoordinates(
			&owm.Coordinates{
				Longitude: cfg.Longitude,
				Latitude:  cfg.Latitude,
			})
		if err != nil {
			span.SetAttributes(attribute.Int("http.status_code", http.StatusInternalServerError))
			span.SetAttributes(attribute.String("error", err.Error()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Record successful API call
		span.SetAttributes(
			attribute.Int("http.status_code", http.StatusOK),
			attribute.String("owm.response.id", fmt.Sprintf("%d", weather.ID)),
			attribute.String("owm.response.name", weather.Name),
		)

		temp := weather.Main.Temp
		tempInt := int(math.Round(temp))

		// Record the temperature in the histogram
		attrs := []attribute.KeyValue{
			attribute.String("unit", cfg.Unit),
			attribute.Float64("longitude", cfg.Longitude),
			attribute.Float64("latitude", cfg.Latitude),
		}
		tempHistogram.Record(ctx, temp, metric.WithAttributes(attrs...))

		// Add temperature to span attributes
		span.SetAttributes(attribute.Float64("temperature", temp))

		w.Write([]byte(fmt.Sprintf("%d", tempInt)))
	}
}
