package main

import (
	"fmt"
	gear_grpc "github.com/misnaged/gear-go/internal/gear-grpc"
	"github.com/misnaged/gear-go/pkg/logger"
	"os"
	"time"
)

var AliceSecret = "0xe5be9a5092b81bca64be81d212e7f2f9eba183bb7a90954f7b76361f6edb5c0a"

// make sure grpc server is running properly!
func main() {
	grpcClient, err := gear_grpc.New("127.0.0.1:9090", 5*time.Second)
	if err != nil {
		logger.Log().Errorf("%v", err)
		os.Exit(1)
	}

	signed, err := grpcClient.CallVoucherIssue(AliceSecret, "10000000000000000000", true, 1000000)
	if err != nil {
		logger.Log().Errorf("%v", err)
		os.Exit(1)
	}

	fmt.Println("gear_grpc_call is", signed.EncodedCall)
}
