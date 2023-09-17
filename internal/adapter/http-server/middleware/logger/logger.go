package logger

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"time"
)

type Handler func(next http.Handler) http.Handler

func New(log *slog.Logger) Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(
			slog.String("component", "middleware/logger"),
		)
		log.Info("Logging middleware enabled!")

		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("URL", r.URL.Path),
				slog.String("user_address", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			startTime := time.Now()
			defer func() {
				entry.Info(
					"request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes_written", ww.BytesWritten()),
					slog.String("duration", time.Since(startTime).String()),
				)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
