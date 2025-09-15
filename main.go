package main

import (
	"fmt"
	"go_schedule_server/configLoader"
	"go_schedule_server/endpointHandlers"
	"go_schedule_server/grpcConnection"
	"net/http"
	"os"
)

func main() {
	config := configLoader.AppConfig{ScraperUrl: "localhost:50051", ServerPort: "8080"}
	err := configLoader.LoadOrCreateYamlConfig("config.yaml", &config, true)
	if err != nil {
		panic(err)
	}

	grpcConnection.CreateGrpcConnection(config)
	defer grpcConnection.GrpcConn.Close()

	http.HandleFunc("GET /scheduleTypes", endpointHandlers.GetScheduleTypesHandler)
	http.HandleFunc("GET /updateTime", endpointHandlers.GetUpdateTimeHandler)
	http.HandleFunc("GET /avaibleScheduleTimeGroups", endpointHandlers.GetAvailableScheduleTimeGroupsHandler)
	http.HandleFunc("GET /schedule", endpointHandlers.GetScheduleHandler)
	http.HandleFunc("GET /scheduleList", endpointHandlers.GetScheduleListHandler)
	err = http.ListenAndServe(":"+config.ServerPort, nil)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to set up server: %v\n", err)
		os.Exit(1)
	}
}
