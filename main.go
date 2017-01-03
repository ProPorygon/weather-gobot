package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"googlemaps.github.io/maps"
)

type Weather struct {
	Currently CurWeather
}

type CurWeather struct {
	Summary     string
	Temperature float64
}

type CallbackMessage struct {
	Text string
}

//G maps client
var client *maps.Client

func getWeatherString(location string) string {
	request := &maps.GeocodingRequest{
		Address: location,
	}

	resp, err := client.Geocode(context.Background(), request)
	if err != nil {
		log.Println(err.Error())
		return "Location not found"
	}

	lat := strconv.FormatFloat(resp[0].Geometry.Location.Lat, 'f', 4, 64)
	lon := strconv.FormatFloat(resp[0].Geometry.Location.Lng, 'f', 4, 64)
	city := resp[0].AddressComponents[0].ShortName

	url := "https://api.darksky.net/forecast/" + os.Getenv("DARKSKY_KEY") + "/" + lat + "," + lon
	res, err := http.Get(url)
	if err != nil {
		log.Println(err.Error())
	}
	if res != nil {
		defer res.Body.Close()
	}
	weather := &Weather{
		Currently: CurWeather{
			Summary:     "summary",
			Temperature: 0.0,
		},
	}
	if err := json.NewDecoder(res.Body).Decode(&weather); err != nil {
		log.Println(err.Error())
	}
	summary := weather.Currently.Summary
	temperature := weather.Currently.Temperature
	response := "The current weather in " + city + " is " + summary + ", with a temperature of " + strconv.FormatFloat(temperature, 'f', 2, 32)
	return response
}

func postToChat(message string) {
	jsonMap := map[string]string{"bot_id": os.Getenv("BOT_ID"), "text": message}
	jsonVal, _ := json.Marshal(jsonMap)
	_, err := http.Post("https://api.groupme.com/v3/bots/post",
		"application/json", bytes.NewReader(jsonVal))
	if err != nil {
		log.Println(err.Error())
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	var responseData CallbackMessage
	json.NewDecoder(r.Body).Decode(&responseData)
	chatMessage := strings.SplitN(responseData.Text, " ", 2)
	if chatMessage[0] == "weatherbot" {
		message := getWeatherString(chatMessage[1])
		postToChat(message)
	}
}

func main() {
	log.Println("Starting")
	var err error
	client, err = maps.NewClient(maps.WithAPIKey(os.Getenv("GMAPS_KEY")))
	if err != nil {
		log.Println(err.Error())
	}
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
