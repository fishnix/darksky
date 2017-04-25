package darksky

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const url string = "https://api.darksky.net/forecast"

// ApiRequest contains all of the data for a request
type ApiRequest struct {
	Lat     string
	Long    string
	Key     string
	Exclude []string
	Lang    string
	Units   string
}

// Forecast is the forecast object returned from the API
type Forecast struct {
	Currently Currently   `json:"currently,omitempty"`
	Daily     Daily       `json:"daily,omitempty"`
	Hourly    Hourly      `json:"hourly,omitempty"`
	Latitude  float64     `json:"latitude,omitempty"`
	Longitude float64     `json:"longitude,omitempty"`
	Offset    int         `json:"offset,omitempty"`
	Timezone  string      `json:"timezone,omitempty"`
	Alerts    []AlertData `json:"alerts,omitempty"`
}

// Currently is the current weather data
type Currently struct {
	ApparentTemperature  float64 `json:"apparentTemperature,omitempty"`
	CloudCover           float64 `json:"cloudCover,omitempty"`
	DewPoint             float64 `json:"dewPoint,omitempty"`
	Humidity             float64 `json:"humidity,omitempty"`
	Icon                 string  `json:"icon,omitempty"`
	NearestStormBearing  float64 `json:"nearestStormBearing,omitempty"`
	NearestStormDistance float64 `json:"nearestStormDistance,omitempty"`
	Ozone                float64 `json:"ozone,omitempty"`
	PrecipIntensity      float64 `json:"precipIntensity,omitempty"`
	PrecipProbability    float64 `json:"precipProbability,omitempty"`
	Pressure             float64 `json:"pressure,omitempty"`
	Summary              string  `json:"summary,omitempty"`
	Temperature          float64 `json:"temperature,omitempty"`
	Time                 int     `json:"time"`
	Visibility           float64 `json:"visibility,omitempty"`
	WindBearing          float64 `json:"windBearing,omitempty"`
	WindSpeed            float64 `json:"windSpeed,omitempty"`
}

// Daily is the daily forecast information
type Daily struct {
	Icon    string      `json:"icon,omitempty"`
	Summary string      `json:"summary,omitempty"`
	Data    []DailyData `json:"data"`
}

// DailyData is the data block for the daily forecast
type DailyData struct {
	ApparentTemperatureMax     float64 `json:"apparentTemperatureMax"`
	ApparentTemperatureMaxTime int     `json:"apparentTemperatureMaxTime"`
	ApparentTemperatureMin     float64 `json:"apparentTemperatureMin"`
	ApparentTemperatureMinTime int     `json:"apparentTemperatureMinTime"`
	CloudCover                 float64 `json:"cloudCover"`
	DewPoint                   float64 `json:"dewPoint"`
	Humidity                   float64 `json:"humidity"`
	Icon                       string  `json:"icon"`
	MoonPhase                  float64 `json:"moonPhase"`
	Ozone                      float64 `json:"ozone"`
	PrecipIntensity            float64 `json:"precipIntensity"`
	PrecipIntensityMax         float64 `json:"precipIntensityMax"`
	PrecipIntensityMaxTime     int     `json:"precipIntensityMaxTime"`
	PrecipProbability          float64 `json:"precipProbability"`
	PrecipType                 string  `json:"precipType"`
	Pressure                   float64 `json:"pressure"`
	Summary                    string  `json:"summary"`
	SunriseTime                int     `json:"sunriseTime"`
	SunsetTime                 int     `json:"sunsetTime"`
	TemperatureMax             float64 `json:"temperatureMax"`
	TemperatureMaxTime         int     `json:"temperatureMaxTime"`
	TemperatureMin             float64 `json:"temperatureMin"`
	TemperatureMinTime         int     `json:"temperatureMinTime"`
	Time                       int     `json:"time"`
	Visibility                 float64 `json:"visibility"`
	WindBearing                float64 `json:"windBearing"`
	WindSpeed                  float64 `json:"windSpeed"`
}

// Hourly is the hourly forecast information
type Hourly struct {
	Icon    string       `json:"icon,omitempty"`
	Summary string       `json:"summary,omitempty"`
	Data    []HourlyData `json:"data"`
}

// HourlyData is the data block for the hourly forecast
type HourlyData struct {
	ApparentTemperature float64 `json:"apparentTemperature"`
	CloudCover          float64 `json:"cloudCover"`
	DewPoint            float64 `json:"dewPoint"`
	Humidity            float64 `json:"humidity"`
	Icon                string  `json:"icon"`
	Ozone               float64 `json:"ozone"`
	PrecipIntensity     float64 `json:"precipIntensity"`
	PrecipProbability   float64 `json:"precipProbability"`
	Pressure            float64 `json:"pressure"`
	Summary             string  `json:"summary"`
	Temperature         float64 `json:"temperature"`
	Time                int     `json:"time"`
	Visibility          float64 `json:"visibility"`
	WindBearing         float64 `json:"windBearing"`
	WindSpeed           float64 `json:"windSpeed"`
}

// Hourly is the minute by minute forecast information
type Minutely struct {
	Icon    string         `json:"icon,omitempty"`
	Summary string         `json:"string,omitempty"`
	Data    []MinutelyData `json:"data"`
}

// MinutelyData is the data block for the minutely forecast
type MinutelyData struct {
	PrecipIntensity   float64 `json:"precipIntensity"`
	PrecipProbability float64 `json:"precipProbability"`
	Time              int     `json:"time"`
}

// AlertData is the data block for weather alerts
type AlertData struct {
	Description string   `json:"description"`
	Expires     int      `json:"expires,omitempty"`
	Regions     []string `json:"regions"`
	Severity    string   `json:"severity"`
	Time        int      `json:"time"`
	Title       string   `json:"title"`
	URI         string   `json:"uri"`
}

// GetForecast calls the darksky API and responds with the forecast
func GetForecast(req ApiRequest) (*Forecast, error) {
	u := buildURL(req)
	f, err := getResponse(u)
	return f, err
}

// buildURL creates the API URL string for the request
func buildURL(req ApiRequest) string {
	var apiURL bytes.Buffer

	apiURL.WriteString(url + "/")
	apiURL.WriteString(req.Key + "/")
	apiURL.WriteString(req.Lat + ",")
	apiURL.WriteString(req.Long)

	if len(req.Units) > 0 {
		apiURL.WriteString("?units=" + req.Units)
	} else {
		apiURL.WriteString("?units=auto")
	}

	if len(req.Exclude) > 0 {
		apiURL.WriteString("&exclude=" + strings.Join(req.Exclude, ","))
	}

	if len(req.Lang) > 0 {
		apiURL.WriteString("&lang=" + req.Lang)
	}

	return apiURL.String()
}

// getResponse takes the URL and returns the parsed forecast
func getResponse(u string) (*Forecast, error) {
	log.Println("[INFO] fetching forecast from", u)

	var forecast *Forecast
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Get(u)
	if err != nil {
		log.Println("[ERROR] Unable to Get: ", err)
		return forecast, err
	}
	defer resp.Body.Close()

	forecast, err = decodeJSON(resp.Body)
	if err != nil {
		log.Println("[ERROR] Unable to decode JSON: ", err)
		return forecast, err
	}

	return forecast, nil
}

// decodeJSON unmarshalls the JSON from DarkSky into our types
func decodeJSON(reader io.Reader) (*Forecast, error) {
	var forecast Forecast
	if err := json.NewDecoder(reader).Decode(&forecast); err != nil {
		log.Println("[ERROR] decoding json:", err)
		return nil, err
	}

	return &forecast, nil
}
