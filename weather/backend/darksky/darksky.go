package darksky

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/eugecm/weatherpi/weather"
)

// Config represents the configuration parameters for the DarkSky API client
type Config struct {
	Key     string
	RootURL string
}

// DarkSky is a simple client for the DarkSky API
type DarkSky struct {
	config Config
	c      http.Client
}

type forecastResponse struct {
	Hourly struct {
		Data []struct {
			Time                     int64
			Summary                  string
			PrecipitationIntensity   float64 `json:"precipIntensity"`
			PrecipitationProbability float64 `json:"precipProbability"`
		}
	}
}

// New initializes the DarkSky client
func New(c Config) *DarkSky {
	if c.RootURL == "" {
		c.RootURL = "https://api.darksky.net"
	}

	return &DarkSky{
		config: c,
	}
}

func (d *DarkSky) doForecast(lat, lon string, t time.Time) (*forecastResponse, error) {
	rsp, err := d.c.Get(fmt.Sprintf("%s/forecast/%s/%s,%s,%d", d.config.RootURL, d.config.Key, lat, lon, t.Unix()))
	if err != nil {
		log.Printf("Error making forecast request: %v\n", err)
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		log.Println(rsp.Request.URL)
		return nil, fmt.Errorf("invalid status code %v. Expected 200", rsp.StatusCode)
	}

	r, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("Error reading forecast response body: %v\n", err)
		return nil, err
	}

	fRsp := &forecastResponse{}
	if err := json.Unmarshal(r, fRsp); err != nil {
		log.Printf("Error unmarshalling forecast reponse: %v\n", err)
		return nil, err
	}

	return fRsp, nil
}

// Forecast returns the forecast for a given day using the DarkSky API
func (d *DarkSky) Forecast(lat, lon string, t time.Time) (weather.Forecast, error) {
	// Truncate time to day to make sure we get 24 hours of forecast
	t = t.Truncate(24 * time.Hour)

	// Make request
	rsp, err := d.doForecast(lat, lon, t)
	if err != nil {
		log.Printf("Error while making call to DarkSky forecast API: %v\n", err)
		return weather.Forecast{}, err
	}

	// TODO: Perform some validation

	f := weather.Forecast{}
	for _, hour := range rsp.Hourly.Data {
		f.Hourly = append(f.Hourly, weather.TimedForecast{
			At: time.Unix(hour.Time, 0),
			PrecipitationIntensity: hour.PrecipitationIntensity,
			PrecipitationChance:    hour.PrecipitationProbability,
		})
	}

	return f, err
}
