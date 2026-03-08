package health

import (
	"context"

	"google.golang.org/grpc/health/grpc_health_v1"
)

type HealthHandler struct {
	grpc_health_v1.HealthServer
}

func NewHealthCheckHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}
