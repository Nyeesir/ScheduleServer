package grpcConnection

import (
	"flag"
	"fmt"
	pb "go_schedule_server/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

var GrpcConn *grpc.ClientConn
var GrpcClient pb.ScheduleScraperClient

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func CreateGrpcConnection() {
	flag.Parse()
	var err error
	GrpcConn, err = grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create scrapper connection: %v\n", err)
		os.Exit(1)
	}
	GrpcClient = pb.NewScheduleScraperClient(GrpcConn)
}
