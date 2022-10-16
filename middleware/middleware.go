package middleware

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"fmt"
	"google.golang.org/grpc/metadata"
)

type Middleware func(ctx context.Context) (context.Context, error)

func MiddlewareFunc(middlewareFunc Middleware) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		fmt.Println("获取用户信息")
		newCtx, err := middlewareFunc(ctx)
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

func GetUserInfo(ctx context.Context) (newCtx context.Context, err error) {
	md, _ := metadata.FromIncomingContext(ctx)
	token, _ := md["token"]
	newCtx = context.WithValue(ctx, "token", token)
	return
}
