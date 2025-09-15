package grpcConnection

import (
	"flag"
	"fmt"
	"go_schedule_server/configLoader"
	pb "go_schedule_server/protos"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var GrpcConn *grpc.ClientConn
var GrpcClient pb.ScheduleScraperClient

func CreateGrpcConnection(config configLoader.AppConfig) {
	flag.Parse()
	var err error
	GrpcConn, err = grpc.NewClient(config.ScraperUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to create scrapper connection: %v\n", err)
		os.Exit(1)
	}
	GrpcClient = pb.NewScheduleScraperClient(GrpcConn)
}
