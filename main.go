package main

import (
	"context"
	"encoding/json"
	"go_schedule_server/cache"
	"go_schedule_server/grpcConnection"
	pb "go_schedule_server/protos"
	"net/http"
	"strings"
	"time"
)

//TODO: BETTER ERROR HANDLING, CONFIGURATION FILES, LOGGING, SECURING CONNECTION

type MessageTemplate struct {
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

func getScheduleTypeshandler(w http.ResponseWriter, r *http.Request) {
	jsonEcoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	message := MessageTemplate{Error: false}

	types, err := cache.GetScheduleTypes(r.Context())
	if err != nil {
		message.Error = true
		message.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		jsonEcoder.Encode(message)
		return
	}

	jsonEcoder.Encode(types)
}

func getUpdateTimeHandler(w http.ResponseWriter, r *http.Request) {
	jsonEcoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	message := MessageTemplate{Error: false}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := grpcConnection.GrpcClient.GetUpdateTime(ctx, &pb.Empty{})
	if err != nil {
		message.Error = true
		message.Message = "Could not get to the scraper"
		w.WriteHeader(http.StatusInternalServerError)
		jsonEcoder.Encode(message)
		return
	}
	updateTime := time.Unix(int64(req.GetTime()), 0).String()
	jsonEcoder.Encode(map[string]string{"time": updateTime})
}

func getAvaibleScheduleTimeGroupsHandler(w http.ResponseWriter, r *http.Request) {
	jsonEcoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	message := MessageTemplate{Error: false}

	timeGroups, err := cache.GetAvailableTimeGroups(r.Context())
	if err != nil {
		message.Error = true
		message.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		jsonEcoder.Encode(message)
		return
	}

	jsonEcoder.Encode(timeGroups)
}

func getScheduleHandler(w http.ResponseWriter, r *http.Request) {
	jsonEcoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	message := MessageTemplate{Error: false}

	reqType := strings.ToLower(r.URL.Query().Get("type"))
	reqId := strings.ToLower(r.URL.Query().Get("id"))
	reqTimeGroup := strings.ToLower(r.URL.Query().Get("time-group"))

	cal, err := cache.GetSchedule(r.Context(), reqType, reqId, reqTimeGroup)
	if err != nil {
		message.Error = true
		message.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		jsonEcoder.Encode(message)
		return
	}

	jsonEcoder.Encode(cal)

}

func getScheduleListHandler(w http.ResponseWriter, r *http.Request) {
	jsonEcoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	message := MessageTemplate{Error: false}

	reqType := strings.ToLower(r.URL.Query().Get("type"))

	scheduleList, err := cache.GetScheduleList(r.Context(), reqType)
	if err != nil {
		message.Error = true
		message.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		jsonEcoder.Encode(message)
		return
	}

	jsonEcoder.Encode(scheduleList)

}

func main() {
	grpcConnection.CreateGrpcConnection()
	defer grpcConnection.GrpcConn.Close()
	http.HandleFunc("GET /scheduleTypes", getScheduleTypeshandler)
	http.HandleFunc("GET /updateTime", getUpdateTimeHandler)
	http.HandleFunc("GET /avaibleScheduleTimeGroups", getAvaibleScheduleTimeGroupsHandler)
	http.HandleFunc("GET /schedule", getScheduleHandler)
	http.HandleFunc("GET /scheduleList", getScheduleListHandler)
	http.ListenAndServe(":8080", nil)
}
