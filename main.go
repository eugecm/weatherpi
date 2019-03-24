package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/eugecm/weatherpi/weather"
	"github.com/eugecm/weatherpi/weather/backend/darksky"
	rpio "github.com/stianeikeland/go-rpio"
)

var (
	lat  = flag.String("lat", "", "Latitude")
	lon  = flag.String("lon", "", "Longitude")
	gpio = flag.Bool("gpio", false, "Output to GPIO")
)

const (
	gpioPinMorningIndicator      = 4
	gpioPinEveningIndicator      = 18
	gpioPinLowChanceIndicator    = 22
	gpioPinMediumChanceIndicator = 23
	gpioPinHighChanceIndicator   = 24
)

func main() {
	flag.Parse()

	switch {
	case *lat == "", *lon == "":
		flag.Usage()
		os.Exit(1)
	}

	f, err := setupDarkSky()
	if err != nil {
		log.Fatalln(err)
	}

	forecast, err := f.Forecast(*lat, *lon, time.Now())
	if err != nil {
		log.Fatalln(err)
	}

	var (
		worstMorningChance float64
		worstEveningChance float64
	)
	for _, h := range forecast.Hourly {
		switch {
		case h.At.Hour() >= 7 && h.At.Hour() < 11:
			if h.PrecipitationChance > worstMorningChance {
				worstMorningChance = h.PrecipitationChance
			}
		case h.At.Hour() >= 17 && h.At.Hour() <= 23:
			if h.PrecipitationChance > worstMorningChance {
				worstEveningChance = h.PrecipitationChance
			}
		}
	}

	log.Printf("morning chance of rain: %v\n", worstMorningChance)
	log.Printf("evening chance of rain: %v\n", worstEveningChance)

	if !*gpio {
		return
	}

	initGPIO()

	if worstMorningChance > 0.05 {
		switchPin(gpioPinMorningIndicator, true)
	} else {
		switchPin(gpioPinMorningIndicator, false)
	}

	if worstEveningChance > 0.05 {
		switchPin(gpioPinEveningIndicator, true)
	} else {
		switchPin(gpioPinEveningIndicator, false)
	}

	worstChance := math.Max(worstMorningChance, worstEveningChance)
	switchChanceIndicator(worstChance)
}

func setupDarkSky() (weather.Forecaster, error) {
	secret := os.Getenv("DARKSKY_SECRET_KEY")
	if secret == "" {
		return nil, fmt.Errorf("No secret key specified. Set up DARKSKY_SECRET_KEY environment variable")
	}

	return darksky.New(darksky.Config{
		Key: secret,
	}), nil
}

func switchChanceIndicator(level float64) {
	switch {
	case level <= 0.05:
		switchPin(gpioPinLowChanceIndicator, false)
		switchPin(gpioPinHighChanceIndicator, false)
		switchPin(gpioPinMediumChanceIndicator, false)
	case level <= 0.33:
		switchPin(gpioPinLowChanceIndicator, false)
		switchPin(gpioPinMediumChanceIndicator, true)
		switchPin(gpioPinHighChanceIndicator, false)
	case level <= 0.66:
		switchPin(gpioPinLowChanceIndicator, false)
		switchPin(gpioPinMediumChanceIndicator, false)
		switchPin(gpioPinHighChanceIndicator, true)
	case level <= 1.1:
		switchPin(gpioPinLowChanceIndicator, true)
		switchPin(gpioPinMediumChanceIndicator, false)
		switchPin(gpioPinHighChanceIndicator, false)
	}
}

func switchPin(pin int, on bool) {
	p := rpio.Pin(pin)
	switch on {
	case true:
		p.High()
	case false:
		p.Low()
	}
}

func initGPIO() {
	if err := rpio.Open(); err != nil {
		panic(err)
	}

	for _, p := range []int{gpioPinMorningIndicator, gpioPinEveningIndicator, gpioPinLowChanceIndicator, gpioPinHighChanceIndicator, gpioPinMediumChanceIndicator} {
		pin := rpio.Pin(p)
		pin.Output()
	}
}
