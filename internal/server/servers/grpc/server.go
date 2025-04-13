package grpcserver

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	pvz "github.com/Sanchir01/golang-avito-proto/pkg/gen/golang/pvz"
	sl "github.com/Sanchir01/golang-avito/pkg/lib/log"
	grpclogging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClientPVZ struct {
	api pvz.PVZServiceClient
	log *slog.Logger
}

func NewGRPCClient(
	ctx context.Context,
	log *slog.Logger,
	port string,
	host string,
	retries int,
) (*GRPCClientPVZ, error) {
	const op = "grpc.New.Client"
	serverAddr := fmt.Sprintf("%s:%s", host, port)
	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.Aborted, codes.NotFound, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retries)),
		grpcretry.WithPerRetryTimeout(1 * time.Second),
	}
	logOpts := []grpclogging.Option{
		grpclogging.WithLogOnEvents(grpclogging.PayloadReceived, grpclogging.PayloadSent),
	}

	cc, err := grpc.DialContext(
		ctx,
		serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpcretry.UnaryClientInterceptor(retryOpts...),
			grpclogging.UnaryClientInterceptor(sl.InterceptorLogger(log), logOpts...),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	client := pvz.NewPVZServiceClient(cc)

	return &GRPCClientPVZ{
		api: client,
		log: log,
	}, nil
}

func (c *GRPCClientPVZ) AllPVZHandler(ctx context.Context) (*pvz.GetPVZListResponse, error) {
	pvz, err := c.api.GetPVZList(ctx, &pvz.GetPVZListRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get PVZ list: %w", err)
	}
	return pvz, nil
}
