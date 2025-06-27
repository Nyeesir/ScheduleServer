package cache

import (
	"context"
	pb "go_schedule_server/protos"
	"time"
)

type cachedScheduleTypesTemplate struct {
	ScheduleTypes   pb.ScheduleTypes
	UpdateTimeStamp int64
}

var scheduleUpdateTime int64 = 0
var lastTimeUpdateCheck int64 = 0
var CachedScheduleTypes cachedScheduleTypesTemplate = cachedScheduleTypesTemplate{}

func GetUpdateTime() int64 {
	if time.Now().Unix()-lastTimeUpdateCheck > int64(time.Hour.Seconds()) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		req, err := main.GrpcClient.GetUpdateTime(ctx, &pb.Empty{})
		if err == nil {
			scheduleUpdateTime = int64(req.Time)
			lastTimeUpdateCheck = time.Now().Unix()
		}
		return scheduleUpdateTime
	}
	return scheduleUpdateTime
}
