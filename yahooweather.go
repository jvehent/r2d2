package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const weatherHelp = "usage: weather <location>. example: weather sarasota, fl"

type Response struct {
	Query Query `json:"query"`
}

type Query struct {
	Count   float64 `json:"count"`
	Lang    string  `json:"lang"`
	Results Results `json:"results"`
}

type Results struct {
	Channel Channel `json:"channel"`
}

type Channel struct {
	Title         string     `json:"title"`
	Link          string     `json:"link"`
	Description   string     `json:"description"`
	Language      string     `json:"language"`
	LastBuildDate string     `json:"lastBuildDate"`
	TTL           string     `json:"ttl"`
	Location      Location   `json:"location"`
	Wind          Wind       `json:"wind"`
	Atmosphere    Atmosphere `json:"atmosphere"`
	Astronomy     Astronomy  `json:"astronomy"`
	Item          Item       `json:"item"`
}
type Location struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Region  string `json:"region"`
}

type Wind struct {
	Chill     string `json:"chill"`
	Direction string `json:"direction"`
	Speed     string `json:"speed"`
}

type Atmosphere struct {
	Humidity   string `json:"humidity"`
	Pressure   string `json:"pressure"`
	Rising     string `json:"rising"`
	Visibility string `json:"visibility"`
}

type Astronomy struct {
	Sunrise string `json:"sunrise"`
	Sunset  string `json:"sunset"`
}

type Item struct {
	Title       string     `json:"title"`
	Lat         string     `json:"lat"`
	Long        string     `json:"long"`
	Link        string     `json:"link"`
	PubDate     string     `json:"pubDate"`
	Condition   Condition  `json:"condition"`
	Description string     `json:"description"`
	Forecast    []Forecast `json:"forecast"`
}

type Condition struct {
	Code string `json:"code"`
	Date string `json:"date"`
	Temp string `json:"temp"`
	Text string `json:"text"`
}

type Forecast struct {
	Code string `json:"code"`
	Date string `json:"date"`
	Day  string `json:"day"`
	High string `json:"high"`
	Low  string `json:"low"`
	Text string `json:"text"`
}

func getYahooForecast(location string) string {
	query := `select * from weather.forecast where woeid in (select woeid from geo.places(1) where text="` + location + `")`
	query = url.QueryEscape(query)
	target := `https://query.yahooapis.com/v1/public/yql?q=` + query + `&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys`
	log.Print(target)
	resp, err := http.Get(target)
	defer resp.Body.Close()
	if err != nil {
		return "Failed to retrieve forecast from Yahoo Weather"
	}
	if resp.StatusCode != 200 {
		return "Yahoo Weather api call returned: " + resp.Status
	}
	var r Response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "Failed to read response from Yahoo Weather api"
	}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return "Invalid response from Yahoo Weather api"
	}
	city := r.Query.Results.Channel.Location.City
	country := r.Query.Results.Channel.Location.Country
	cond := r.Query.Results.Channel.Item.Condition
	var tomorrow, in2days Forecast
	if len(r.Query.Results.Channel.Item.Forecast) > 3 {
		tomorrow = r.Query.Results.Channel.Item.Forecast[1]
		in2days = r.Query.Results.Channel.Item.Forecast[2]
	}
	ast := r.Query.Results.Channel.Astronomy
	f, err := strconv.Atoi(cond.Temp)
	if err != nil {
		return "Temperature convertion failed for: " + cond.Temp
	}
	return fmt.Sprintf("overcast in %s, %s: %dF/%dC. daylight from %s to %s. forecast: %s high %s low %s (%s); %s high %s low %s (%s)",
		city, country, f, int32(float64(f-32)*0.5555), ast.Sunrise, ast.Sunset,
		tomorrow.Day, tomorrow.High, tomorrow.Low, tomorrow.Text,
		in2days.Day, in2days.High, in2days.Low, in2days.Text)

}
