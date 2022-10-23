package interceptor

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/status"
	"runtime/debug"
	"time"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	RequestID = "x-request-id"
)

func AccessLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	beginTime := time.Now()
	beginTimeUnix := beginTime.Local().Unix()
	zap.S().Infof("access request log: method: %s, begin_time: %d, request: %v, metadata: %v\n",
		info.FullMethod, beginTimeUnix, req, md)

	resp, err := handler(ctx, req)

	endTimeUnix := time.Now().Local().Unix()
	zap.S().Infof("access response log: method: %s, begin_time: %d, end_time: %d, cost:%s,response: %v, metadata: %v",
		info.FullMethod, beginTimeUnix, endTimeUnix, time.Since(beginTime), resp, md)
	return resp, err
}

// 普通错误记录的日志拦截器
func ErrorLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	resp, err := handler(ctx, req)
	if err != nil {
		s,_:= status.FromError(err)
		zap.S().Infof("error log: method: %s, code: %v, message: %v, details: %v,metadata: %v\n",
			info.FullMethod, s.Code(), s.Err().Error(), s.Details(), md)
	}
	return resp, err
}

// 异常捕抓拦截器
func Recovery(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	defer func() {
		if e := recover(); e != nil {
			zap.S().Info("recovery log: method: %s, message: %v, stack: %s", info.FullMethod, e, string(debug.Stack()[:]))
		}
	}()

	return handler(ctx, req)
}

func RequestIDServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.Pairs()
		}
		// Set request ID for context.
		requestIDs := md[string(RequestID)]
		if len(requestIDs) >= 1 {
			ctx = context.WithValue(ctx, RequestID, requestIDs[0])
			return handler(ctx, req)
		}

		// Generate request ID and set context if not exists.
		requestID := uuid.NewString()
		ctx = context.WithValue(ctx, RequestID, requestID)
		return handler(ctx, req)
	}
}


