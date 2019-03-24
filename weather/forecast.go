package weather

import "time"

// TimedForecast is a weather forecast at a given time
type TimedForecast struct {
	// At specifies the time for which this forecast is made
	At time.Time
	// PrecipitationIntensity is the amount of rain in millimeters per hour
	PrecipitationIntensity float64
	// PrecipitationChance is the probability of rain (between 0 and 1)
	PrecipitationChance float64
}

// Forecast is the predicted weather conditions for a day
type Forecast struct {
	// Hourly represents a series of forecast made by the hour.
	// A valid Forecast has 24 Hourly TimedForecast elements
	Hourly []TimedForecast
}

// Forecaster can Forecast the weather at a certain time
type Forecaster interface {
	Forecast(string, string, time.Time) (Forecast, error)
}
