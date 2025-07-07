package gear_grpc

import (
	"context"
	"fmt"
	"github.com/misnaged/gear-go/internal/models"
	proto "github.com/misnaged/gear-go/lib/server_grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"sync"
	"time"
)

type IGrpcClient interface {
	CallVoucherIssue(spender, balance string, upload bool, dur int) (*models.EncodedVoucherIssue, error)
	Close() error
}
type GrpcClient struct {
	timeout time.Duration
	*grpc.ClientConn
	mu sync.Mutex
	proto.GearGrpcServiceClient
	grpc_health_v1.HealthClient
}

func New(addr string, timeOut time.Duration) (IGrpcClient, error) {
	client := &GrpcClient{timeout: timeOut}

	if err := client.initConn(addr); err != nil {
		return nil, fmt.Errorf("creation of Grpc Client failed:  %w", err)
	}
	client.HealthClient = grpc_health_v1.NewHealthClient(client.ClientConn)

	client.GearGrpcServiceClient = proto.NewGearGrpcServiceClient(client.ClientConn)
	return client, nil
}

func (cli *GrpcClient) initConn(addr string) (err error) {
	var kacp = keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}

	connParams := grpc.WithConnectParams(grpc.ConnectParams{
		Backoff: backoff.Config{
			BaseDelay:  100 * time.Millisecond,
			Multiplier: 1.2,
			MaxDelay:   1 * time.Second,
		},
		MinConnectTimeout: 5 * time.Second,
	})
	cli.ClientConn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithKeepaliveParams(kacp), connParams)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}
	return
}

func (cli *GrpcClient) CallVoucherIssue(spender, balance string, upload bool, dur int) (*models.EncodedVoucherIssue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cli.timeout)
	defer cancel()
	params := &models.VoucherParams{
		Spender:       spender,
		Balance:       balance,
		CodeUploading: upload,
		Duration:      int32(dur),
	}
	resp, err := cli.GearGrpcServiceClient.CallVoucherIssue(ctx, models.VoucherParamsToProto(params))
	if err != nil {
		return nil, fmt.Errorf("GetMeals api request: %w", err)
	}
	return models.EncodedCallFromProto(resp), nil
}
