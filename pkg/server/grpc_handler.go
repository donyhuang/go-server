package server

import "google.golang.org/grpc"

type GrpcHandler struct {
	serviceDesc grpc.ServiceDesc
	handle      interface{}
}

func NewGrpcHandler(serviceDesc grpc.ServiceDesc, handle interface{}) *GrpcHandler {
	return &GrpcHandler{
		serviceDesc: serviceDesc,
		handle:      handle,
	}
}
