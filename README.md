# Go Darksky Client

[DarkSky](https://darksky.net/)

Because it seemed like the internet needed another one... :D

## Example

```golang
package main

import (
	"fmt"

	"github.com/fishnix/darksky"
)

func main() {
	ex := []string{"hourly", "minutely"}
	request := darksky.ApiRequest{
		Lat:     "51.894444",
		Long:    "1.482500",
		Key:     "supersecretapikey",
		Exclude: ex,
		Lang:    "",
		Units:   "",
	}
	f, _ := darksky.GetForecast(request)
	fmt.Println("Currently it is:", f.Currently.Temperature)
	for _, v := range f.Daily.Data {
		fmt.Printf("%d:\n    Summary: %s\n    will be between %f and %f\n", v.Time, v.Summary, v.TemperatureMin, v.TemperatureMax)
	}
}
```

[![GoDoc](https://godoc.org/github.com/fishnix/darksky?status.svg)](https://godoc.org/github.com/fishnix/darksky)
