package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	var responseData map[string]interface{}
	json.NewDecoder(r.Body).Decode(&responseData)
	if responseData["text"] == "weatherbot" {
		jsonMap := map[string]string{"bot_id": os.Getenv("BOT_ID"), "text": "Hello"}
		jsonVal, _ := json.Marshal(jsonMap)
		_, err := http.Post("https://api.groupme.com/v3/bots/post",
			"application/json", bytes.NewReader(jsonVal))
		if err != nil {
			log.Printf("%v", err.Error())
		}
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
