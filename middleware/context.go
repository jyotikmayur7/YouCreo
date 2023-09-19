package middleware

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
)

func AddContextInterceptorUnary(ctxMain context.Context) grpc.UnaryServerInterceptor {
	return func(ctxSubMain context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log := hclog.Default()
		log.Error("Test Middleware from unary")
		test := "testStringFromMiddleware"
		ctx := context.WithValue(ctxMain, "test", test)
		return handler(ctx, req)
	}
}

func AddContextInterceptorStream(ctxMain context.Context) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log := hclog.Default()
		log.Error("Test Middleware from stream")
		test := "testStringFromMiddleware"
		ctx := context.WithValue(ctxMain, "test", test)
		ss.Context().Value(ctx)
		return handler(srv, ss)
	}
}
