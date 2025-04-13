package sl

import (
	"context"
	"log/slog"

	grpclogging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
func InterceptorLogger(l *slog.Logger) grpclogging.Logger {
	return grpclogging.LoggerFunc(func(ctx context.Context, lvl grpclogging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
