package main

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type health struct {
	Status grpc_health_v1.HealthCheckResponse_ServingStatus
	Reason string
}

func (h *health) Watch(*grpc_health_v1.HealthCheckRequest, grpc_health_v1.Health_WatchServer) error {
	return nil
}
func (h *health) OffLine(reason string) {
	h.Status = grpc_health_v1.HealthCheckResponse_NOT_SERVING
	h.Reason = reason
	fmt.Println(reason)
}
func (h *health) OnLine(reason string) {
	h.Status = grpc_health_v1.HealthCheckResponse_SERVING
	h.Reason = reason
	fmt.Println(reason)
}
func (h *health) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	spew.Dump(req)

	return &grpc_health_v1.HealthCheckResponse{Status: h.Status}, nil
}
