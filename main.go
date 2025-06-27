package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"go_schedule_server/icsProcessing"
	pb "go_schedule_server/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"strings"
	"time"
)

//TODO: BETTER ERROR HANDLING, CONFIGURATION FILES

var grpcConn *grpc.ClientConn
var GrpcClient pb.ScheduleScraperClient

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

type MessageTemplate struct {
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

func createGrpcConnection() {
	flag.Parse()
	var err error
	grpcConn, err = grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create scrapper connection: %v\n", err)
		os.Exit(1)
	}
	GrpcClient = pb.NewScheduleScraperClient(grpcConn)
}

func getScheduleTypeshandler(w http.ResponseWriter, r *http.Request) {
	jsonEcoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	message := MessageTemplate{Error: false}

	//print("checking if cache is valid")
	//if cache.CachedScheduleTypes.ScheduleTypes.ScheduleTypes != nil && cache.CachedScheduleTypes.UpdateTimeStamp == cache.GetUpdateTime() {
	//	print("cache is valid")
	//	jsonEcoder.Encode(cache.CachedScheduleTypes.ScheduleTypes)
	//	return
	//}

	print("cache is not valid")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := GrpcClient.GetScheduleTypes(ctx, &pb.Empty{})
	if err != nil {
		message.Error = true
		message.Message = "Could not get to the scraper"
		w.WriteHeader(http.StatusInternalServerError)
		jsonEcoder.Encode(message)
		return
	}

	jsonEcoder.Encode(req)
}

func getUpdateTimeHandler(w http.ResponseWriter, r *http.Request) {
	jsonEcoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	message := MessageTemplate{Error: false}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := GrpcClient.GetUpdateTime(ctx, &pb.Empty{})
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := GrpcClient.GetAvailableScheduleTimeGroups(ctx, &pb.Empty{})
	if err != nil {
		message.Error = true
		message.Message = "Could not get to the scraper"
		w.WriteHeader(http.StatusInternalServerError)
		jsonEcoder.Encode(message)
		return
	}

	jsonEcoder.Encode(req)
}

func getScheduleHandler(w http.ResponseWriter, r *http.Request) {
	jsonEcoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	message := MessageTemplate{Error: false}

	reqType := strings.ToLower(r.URL.Query().Get("type"))
	reqId := strings.ToLower(r.URL.Query().Get("id"))
	reqTimeGroup := strings.ToLower(r.URL.Query().Get("time-group"))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := GrpcClient.GetScheduleFileAsStr(ctx, &pb.ScheduleFileRequest{SchedType: reqType, SchedId: reqId, TimeGroup: reqTimeGroup})
	if err != nil {
		message.Error = true
		message.Message = "Could not get to the scraper"
		w.WriteHeader(http.StatusInternalServerError)
		jsonEcoder.Encode(message)
		return
	}

	cal, err := icsProcessing.Parse(req.GetContent())
	if err != nil {
		message.Error = true
		message.Message = "Could not parse the ics file"
		w.WriteHeader(http.StatusInternalServerError)
	}

	jsonEcoder.Encode(cal)
}

func getScheduleListHandler(w http.ResponseWriter, r *http.Request) {
	jsonEcoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	message := MessageTemplate{Error: false}

	reqType := strings.ToLower(r.URL.Query().Get("type"))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := GrpcClient.GetScheduleList(ctx, &pb.ScheduleTypeRequest{Type: reqType})
	if err != nil {
		message.Error = true
		message.Message = "Could not get to the scraper"
		w.WriteHeader(http.StatusInternalServerError)
		jsonEcoder.Encode(message)
		return
	}

	jsonEcoder.Encode(req)
}

func main() {
	createGrpcConnection()
	defer grpcConn.Close()
	http.HandleFunc("GET /scheduleTypes", getScheduleTypeshandler)
	http.HandleFunc("GET /updateTime", getUpdateTimeHandler)
	http.HandleFunc("GET /avaibleScheduleTimeGroups", getAvaibleScheduleTimeGroupsHandler)
	http.HandleFunc("GET /schedule", getScheduleHandler)
	http.HandleFunc("GET /scheduleList", getScheduleListHandler)
	http.ListenAndServe(":8080", nil)
}
