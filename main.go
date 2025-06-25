package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	pb "go_schedule_server/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"time"
)

var grpcConn *grpc.ClientConn
var grpcClient pb.ScheduleScraperClient

var someMessage string = "Asd"

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
	grpcClient = pb.NewScheduleScraperClient(grpcConn)
}

func getScheduleTypeshandler(w http.ResponseWriter, r *http.Request) {
	jsonEcoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	message := MessageTemplate{Error: false}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := grpcClient.GetScheduleTypes(ctx, &pb.Empty{})
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

	req, err := grpcClient.GetUpdateTime(ctx, &pb.Empty{})
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

func main() {
	createGrpcConnection()
	defer grpcConn.Close()
	http.HandleFunc("GET /scheduleTypes", getScheduleTypeshandler)
	http.HandleFunc("GET /updateTime", getUpdateTimeHandler)
	http.ListenAndServe(":8080", nil)
}
