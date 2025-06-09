package main

import (
	"encoding/json"
	"fmt"
	"go_schedule_server/icsProcessing"
	"net/http"
	"os"
)

type MessageTemplate struct {
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	jsonEcoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	message := MessageTemplate{Error: false}

	dat, err := os.ReadFile("INŻIIIIO1.ics")
	if err != nil {
		message.Error = true
		message.Message = "Could not read file"
		w.WriteHeader(http.StatusInternalServerError)
		jsonEcoder.Encode(message)
	}
	//fmt.Println(string(dat))

	scheduleJSON, err := icsProcessing.IcsToJson(string(dat))
	if err != nil {
		message.Error = true
		message.Message = "Could not parse file"
		w.WriteHeader(http.StatusInternalServerError)
		jsonEcoder.Encode(message)
	}
	fmt.Println(string(scheduleJSON))
}

func main() {
	http.HandleFunc("GET /schedules/test", handler)
	http.ListenAndServe(":8080", nil)
}
