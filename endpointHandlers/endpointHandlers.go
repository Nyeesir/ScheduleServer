package endpointHandlers

import (
	"encoding/json"
	"fmt"
	"go_schedule_server/cache"
	"go_schedule_server/grpcConnection"
	pb "go_schedule_server/protos"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MessageTemplate struct {
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

func GetScheduleTypesHandler(w http.ResponseWriter, r *http.Request) {
	jsonEcoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	message := MessageTemplate{Error: false}

	types, err := cache.GetScheduleTypes(r.Context())
	if err != nil {
		fmt.Println(err)
		message.Error = true
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.DeadlineExceeded:
			message.Message = "Could not get the data from outside server"
			w.WriteHeader(http.StatusGatewayTimeout)
		case codes.NotFound:
			message.Message = "Could not get the schedule types"
			w.WriteHeader(http.StatusInternalServerError)
		default:
			message.Message = "Could not get to the scraper"
			w.WriteHeader(http.StatusGatewayTimeout)
		}
		jsonEcoder.Encode(message)
		return
	}

	jsonEcoder.Encode(types)
}

func GetUpdateTimeHandler(w http.ResponseWriter, r *http.Request) {
	jsonEcoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	message := MessageTemplate{Error: false}

	req, err := grpcConnection.GrpcClient.GetUpdateTime(r.Context(), &pb.Empty{})
	if err != nil {
		fmt.Println(err)
		message.Error = true
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.DeadlineExceeded:
			message.Message = "Could not get the data from outside server"
			w.WriteHeader(http.StatusGatewayTimeout)
		case codes.NotFound:
			message.Message = "Could not get the schedule types"
			w.WriteHeader(http.StatusInternalServerError)
		default:
			message.Message = "Could not get to the scraper"
			w.WriteHeader(http.StatusGatewayTimeout)
		}
		jsonEcoder.Encode(message)
		return
	}

	updateTime := time.Unix(int64(req.GetTime()), 0).String()
	jsonEcoder.Encode(map[string]string{"time": updateTime})
}

func GetAvailableScheduleTimeGroupsHandler(w http.ResponseWriter, r *http.Request) {
	jsonEcoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	message := MessageTemplate{Error: false}

	timeGroups, err := cache.GetAvailableTimeGroups(r.Context())
	if err != nil {
		fmt.Println(err)
		message.Error = true
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.DeadlineExceeded:
			message.Message = "Could not get the data from outside server"
			w.WriteHeader(http.StatusGatewayTimeout)
		case codes.NotFound:
			message.Message = "Could not get the schedule types"
			w.WriteHeader(http.StatusInternalServerError)
		default:
			message.Message = "Could not get to the scraper"
			w.WriteHeader(http.StatusGatewayTimeout)
		}
		jsonEcoder.Encode(message)
		return
	}

	jsonEcoder.Encode(timeGroups)
}

func GetScheduleHandler(w http.ResponseWriter, r *http.Request) {
	jsonEcoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	message := MessageTemplate{Error: false}

	reqType := strings.ToLower(r.URL.Query().Get("type"))
	reqId := strings.ToLower(r.URL.Query().Get("id"))
	reqTimeGroup := strings.ToLower(r.URL.Query().Get("time-group"))
	reqTimeGroupType := strings.ToLower(r.URL.Query().Get("time-group-type"))

	cal, err := cache.GetSchedule(r.Context(), reqType, reqId, reqTimeGroup, reqTimeGroupType)
	if err != nil {
		fmt.Println(err)
		message.Error = true
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.DeadlineExceeded:
			message.Message = "Could not get the data from outside server"
			w.WriteHeader(http.StatusGatewayTimeout)
		case codes.InvalidArgument:
			message.Message = "Invalid arguments"
			w.WriteHeader(http.StatusNotFound)
		default:
			message.Message = "Could not get to the scraper"
			w.WriteHeader(http.StatusGatewayTimeout)
		}
		jsonEcoder.Encode(message)
		return
	}

	jsonEcoder.Encode(cal)

}

func GetScheduleListHandler(w http.ResponseWriter, r *http.Request) {
	jsonEcoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	message := MessageTemplate{Error: false}

	reqType := strings.ToLower(r.URL.Query().Get("type"))

	scheduleList, err := cache.GetScheduleList(r.Context(), reqType)
	if err != nil {
		fmt.Println(err)
		message.Error = true
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.DeadlineExceeded:
			message.Message = "Could not get the data from outside server"
			w.WriteHeader(http.StatusGatewayTimeout)
		case codes.InvalidArgument:
			message.Message = "Invalid arguments"
			w.WriteHeader(http.StatusNotFound)
		default:
			message.Message = "Could not get to the scraper"
			w.WriteHeader(http.StatusGatewayTimeout)
		}
		jsonEcoder.Encode(message)
		return
	}

	jsonEcoder.Encode(scheduleList)

}
