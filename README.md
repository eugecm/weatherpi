# WeatherPI

Simple rain indicator made for the Raspberry PI.

```bash
# Check the weather in London
./weatherpi -lat 51.5074 -lon 0.1278 -gpio
```

Setting the `-gpio` flag will output to the following GPIO pins:
* Pin 4: Rain in the morning
* Pin 18: Rain in the evening
* Pin 22: Low chance of rain
* Pin 23: Medium chance of rain
* Pin 24: High chance of rain

It can use any weather API as long as it implements a basic `Forecaster` interface. A
`Forecaster` for the DarkSky API is included (but you need your own key).