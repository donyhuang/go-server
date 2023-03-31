package interceptor

import (
	"context"
	"gitlab.nongchangshijie.com/go-base/server/pkg/log"
	"gitlab.nongchangshijie.com/go-base/server/pkg/trace"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

var LogEntryKey = struct{}{}

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	traceId := trace.IdFromContext(ctx)
	ctx = log.InjectLogEntryToContext(traceId, ctx)
	grpclog.Warningf("gRPC method: %s, %v traceId:%v", info.FullMethod, req, traceId)
	resp, err := handler(ctx, req)
	grpclog.Warningf("gRPC method: %s, %v traceId:%v", info.FullMethod, resp, traceId)
	return resp, err
}

func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	traceId := trace.IdFromContext(ctx)
	defer func() {
		if e := recover(); e != nil {
			grpclog.Errorf("grpcRecover traceId %v,stack %v", traceId, string(debug.Stack()))
			//debug.PrintStack()
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()
	return handler(ctx, req)
}
