package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Weather struct {
	Currently CurWeather
}

type CurWeather struct {
	Summary     string
	Temperature float64
}

func getWeatherString() string {
	url := "https://api.darksky.net/forecast/" + os.Getenv("DARKSKY_KEY") + "/40.1164,-88.2434"
	log.Println(url)
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
	response := "The current weather is " + summary + ", with a temperature of " + strconv.FormatFloat(temperature, 'f', 2, 32)
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
	var responseData map[string]interface{}
	json.NewDecoder(r.Body).Decode(&responseData)
	if responseData["text"] == "weatherbot" {
		message := getWeatherString()
		postToChat(message)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
