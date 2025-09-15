package cache

//TODO: ADD LIMIT FOR CACHE AND CACHING TO FILE

import (
	"context"
	"fmt"
	"go_schedule_server/grpcConnection"
	"go_schedule_server/icsProcessing"
	pb "go_schedule_server/protos"
	"sync"
	"time"
)

type scheduleListStorage struct {
	items map[string]*pb.ScheduleListResponse
}

type scheduleFileStorage struct {
	items []*schedule
}

type schedule struct {
	Id            string
	Type          string
	TimeGroup     string
	TimeGroupType string
	content       *icsProcessing.Calendar
}

var scheduleUpdateTime time.Time
var lastTimeUpdateCheck time.Time
var cacheUpdateTimestamp time.Time

var cachedScheduleTypes *pb.ScheduleTypes
var cachedScheduleLists *scheduleListStorage
var cachedAvailableTimeGroups *pb.AvailableTimeGroups
var cachedScheduleFiles *scheduleFileStorage

var (
	scheduleTypesMutex       sync.RWMutex
	availableTimeGroupsMutex sync.RWMutex
	scheduleListsMutex       sync.RWMutex
	scheduleFilesMutex       sync.RWMutex
	initOnce                 sync.Once
)

func init() {
	initOnce.Do(func() {
		cachedScheduleTypes = &pb.ScheduleTypes{}
		cachedScheduleLists = &scheduleListStorage{
			items: make(map[string]*pb.ScheduleListResponse),
		}
		cachedScheduleFiles = &scheduleFileStorage{}
		cachedAvailableTimeGroups = &pb.AvailableTimeGroups{}
		cacheUpdateTimestamp = time.Time{}
	})
}

func CheckIfCacheIsUpToDate() bool {
	if cacheUpdateTimestamp.Before(getUpdateTime()) {
		clearCache()
		return false
	} else {
		return true
	}
}

func getUpdateTime() time.Time {
	if time.Now().Sub(lastTimeUpdateCheck) > time.Hour {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		req, err := grpcConnection.GrpcClient.GetUpdateTime(ctx, &pb.Empty{})
		if err == nil {
			scheduleUpdateTime = time.Unix(int64(req.GetTime()), 0)
			lastTimeUpdateCheck = time.Now()
		}
		return scheduleUpdateTime
	}
	return scheduleUpdateTime

}

func clearCache() {
	cachedScheduleTypes.ScheduleTypes = nil
	cachedScheduleLists.items = map[string]*pb.ScheduleListResponse{}
	cachedAvailableTimeGroups.Periods = nil
	cachedAvailableTimeGroups.Weeks = nil
	cachedScheduleFiles.items = nil
	return
}

func GetScheduleTypes(ctx context.Context) (*pb.ScheduleTypes, error) {
	scheduleTypesMutex.RLock()
	if !CheckIfCacheIsUpToDate() {
		fmt.Println("Schedule cache is outdated")
		scheduleTypesMutex.RUnlock()
		return refreshScheduleTypes(ctx)
	}

	if cachedScheduleTypes != nil && len(cachedScheduleTypes.ScheduleTypes) > 0 {
		fmt.Println("Returning cached schedule types")
		result := cachedScheduleTypes
		scheduleTypesMutex.RUnlock()
		return result, nil
	}
	scheduleTypesMutex.RUnlock()

	fmt.Println("Schedule cache is empty")
	return refreshScheduleTypes(ctx)
}

func refreshScheduleTypes(ctx context.Context) (*pb.ScheduleTypes, error) {
	fmt.Println("Refreshing schedule types cache")
	scheduleTypesMutex.Lock()
	defer scheduleTypesMutex.Unlock()

	if CheckIfCacheIsUpToDate() && cachedScheduleTypes != nil && len(cachedScheduleTypes.ScheduleTypes) > 0 {
		return cachedScheduleTypes, nil
	}

	scheduleTypes, err := grpcConnection.GrpcClient.GetScheduleTypes(ctx, &pb.Empty{})
	if err != nil {
		fmt.Println("Error while getting schedule types from grpc server")
		return nil, err
	}

	cachedScheduleTypes = scheduleTypes
	cacheUpdateTimestamp = time.Now()

	return scheduleTypes, nil
}

func GetAvailableTimeGroups(ctx context.Context) (*pb.AvailableTimeGroups, error) {
	availableTimeGroupsMutex.RLock()
	if !CheckIfCacheIsUpToDate() {
		fmt.Println("Time groups cache is outdated")
		availableTimeGroupsMutex.RUnlock()
		return refreshAvailableTimeGroups(ctx)
	}

	if cachedAvailableTimeGroups != nil && (len(cachedAvailableTimeGroups.Periods) > 0 || len(cachedAvailableTimeGroups.Weeks) > 0) {
		fmt.Println("Returning cached time groups")
		result := cachedAvailableTimeGroups
		availableTimeGroupsMutex.RUnlock()
		return result, nil
	}
	availableTimeGroupsMutex.RUnlock()

	fmt.Println("Time groups cache is empty")
	return refreshAvailableTimeGroups(ctx)
}

func refreshAvailableTimeGroups(ctx context.Context) (*pb.AvailableTimeGroups, error) {
	fmt.Println("Refreshing time groups cache")
	availableTimeGroupsMutex.Lock()
	defer availableTimeGroupsMutex.Unlock()

	if CheckIfCacheIsUpToDate() && cachedAvailableTimeGroups != nil &&
		(len(cachedAvailableTimeGroups.Periods) > 0 || len(cachedAvailableTimeGroups.Weeks) > 0) {
		return cachedAvailableTimeGroups, nil
	}

	timeGroups, err := grpcConnection.GrpcClient.GetAvailableScheduleTimeGroups(ctx, &pb.Empty{})
	if err != nil {
		fmt.Println("Error while getting time groups from grpc server")
		return nil, err
	}

	cachedAvailableTimeGroups = timeGroups
	cacheUpdateTimestamp = time.Now()

	return timeGroups, nil
}

func GetScheduleList(ctx context.Context, scheduleType string) (*pb.ScheduleListResponse, error) {
	scheduleListsMutex.RLock()
	if !CheckIfCacheIsUpToDate() {
		fmt.Println("Schedule list cache is outdated")
		scheduleListsMutex.RUnlock()
		return refreshScheduleList(ctx, scheduleType)
	}

	if cachedScheduleLists != nil {
		if list, exists := cachedScheduleLists.items[scheduleType]; exists {
			fmt.Printf("Returning cached schedule list for type: %s\n", scheduleType)
			scheduleListsMutex.RUnlock()
			return list, nil
		}
	}
	scheduleListsMutex.RUnlock()

	fmt.Printf("Schedule list cache is empty for type: %s\n", scheduleType)
	return refreshScheduleList(ctx, scheduleType)
}

func refreshScheduleList(ctx context.Context, scheduleType string) (*pb.ScheduleListResponse, error) {
	fmt.Printf("Refreshing schedule list cache for type: %s\n", scheduleType)
	scheduleListsMutex.Lock()
	defer scheduleListsMutex.Unlock()

	if CheckIfCacheIsUpToDate() {
		if list, exists := cachedScheduleLists.items[scheduleType]; exists {
			return list, nil
		}
	}

	scheduleList, err := grpcConnection.GrpcClient.GetScheduleList(ctx, &pb.ScheduleTypeRequest{Type: scheduleType})
	if err != nil {
		fmt.Printf("Error while getting schedule list from grpc server for type: %s\n", scheduleType)
		return nil, err
	}

	cachedScheduleLists.items[scheduleType] = scheduleList
	cacheUpdateTimestamp = time.Now()

	return scheduleList, nil
}

func GetSchedule(ctx context.Context, schedType string, schedId string, timeGroup string, timeGroupType string) (*icsProcessing.Calendar, error) {
	scheduleFilesMutex.RLock()
	if !CheckIfCacheIsUpToDate() {
		fmt.Println("Schedule cache is outdated")
		scheduleFilesMutex.RUnlock()
		return refreshSchedule(ctx, schedType, schedId, timeGroup, timeGroupType)
	}

	if cachedScheduleFiles != nil && cachedScheduleFiles.items != nil {
		for _, s := range cachedScheduleFiles.items {
			if s.Type == schedType && s.Id == schedId && s.TimeGroup == timeGroup && s.TimeGroupType == timeGroupType {
				fmt.Printf("Returning cached schedule for type: %s, id: %s, timeGroup: %s\n",
					schedType, schedId, timeGroup)
				result := s.content
				scheduleFilesMutex.RUnlock()
				return result, nil
			}
		}
	}
	scheduleFilesMutex.RUnlock()

	fmt.Printf("Schedule cache is empty for type: %s, id: %s, timeGroup: %s\n",
		schedType, schedId, timeGroup)
	return refreshSchedule(ctx, schedType, schedId, timeGroup, timeGroupType)
}

func refreshSchedule(ctx context.Context, schedType string, schedId string, timeGroup string, timeGroupType string) (*icsProcessing.Calendar, error) {
	fmt.Printf("Refreshing schedule cache for type: %s, id: %s, timeGroup: %s, timeGroupType: %s\n",
		schedType, schedId, timeGroup, timeGroupType)
	scheduleFilesMutex.Lock()
	defer scheduleFilesMutex.Unlock()

	if CheckIfCacheIsUpToDate() {
		for _, s := range cachedScheduleFiles.items {
			if s.Type == schedType && s.Id == schedId && s.TimeGroup == timeGroup && s.TimeGroupType == timeGroupType {
				return s.content, nil
			}
		}
	}

	req, err := grpcConnection.GrpcClient.GetScheduleFileAsStr(ctx, &pb.ScheduleFileRequest{
		SchedType:     schedType,
		SchedId:       schedId,
		TimeGroup:     timeGroup,
		TimeGroupType: timeGroupType,
	})
	if err != nil {
		fmt.Printf("Error while getting schedule from grpc server: %v\n", err)
		return nil, err
	}

	cal, err := icsProcessing.Parse(req.GetContent())
	if err != nil {
		fmt.Printf("Error while parsing schedule: %v\n", err)
		return nil, err
	}

	newSchedule := &schedule{
		Id:            schedId,
		Type:          schedType,
		TimeGroup:     timeGroup,
		TimeGroupType: timeGroupType,
		content:       &cal,
	}

	if cachedScheduleFiles.items == nil {
		cachedScheduleFiles.items = make([]*schedule, 0)
	}
	cachedScheduleFiles.items = append(cachedScheduleFiles.items, newSchedule)
	cacheUpdateTimestamp = time.Now()

	return &cal, nil
}
