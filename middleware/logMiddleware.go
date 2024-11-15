package middleware

import (
	"bytes"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// Init Logger config the zap logger with JSON encoding and additional options
func InitializeLogger() error {
	config := zap.NewProductionConfig()
	config.Encoding = "json"
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	var err error
	Logger, err = config.Build(zap.AddCaller())
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(Logger) // zap set global logger
	return nil
}

type contextKey string

const loggerKey = contextKey("logger")

// WithLogger returns a context with the logger attached for passing through the application.
func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)

}

// zap logger is set globbally so below funtion is not required but still can be useful for calling zap logger
func GetLogger(ctx context.Context) *zap.Logger {
	logger, _ := ctx.Value(loggerKey).(*zap.Logger)
	if logger != nil {
		return logger
	}
	return Logger
}

type ResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func NewResponseWriter(w gin.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w, body: &bytes.Buffer{}}
}

// Write overrides the Write method to capture the response body
func (w *ResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		logger := Logger.With(zap.String("request_id", requestID))
		zap.ReplaceGlobals(logger)

		// Retrieve the logger with the request context
		ctx := WithLogger(c.Request.Context(), logger)

		// Attach the updated context to the request
		c.Request = c.Request.WithContext(ctx)

		rw := NewResponseWriter(c.Writer)
		c.Writer = rw

		c.Next()

		defer zap.ReplaceGlobals(Logger)

		zap.L().Info("API request",
			zap.String("url", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.Int("status", c.Writer.Status()),
			zap.String("response_content", rw.body.String()),
			zap.String("response", c.Writer.Header().Get("Content-Type")),
		)
	}
}
