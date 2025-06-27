package grpcserver

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		logger := zap.L()

		/* log incoming request */
		logger.Info("Incoming gRPC request",
			zap.String("method", info.FullMethod),
			zap.Any("request", req),
		)

		/* handle request */
		resp, err = handler(ctx, req)

		/* log response or error */
		if err != nil {
			st, _ := status.FromError(err)
			logger.Error("gRPC request failed",
				zap.String("method", info.FullMethod),
				zap.Error(err),
				zap.String("code", st.Code().String()),
			)
		} else {
			logger.Info("gRPC response",
				zap.String("method", info.FullMethod),
				zap.Any("response", resp),
			)
		}

		return resp, err
	}
}
