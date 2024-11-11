package middleware

// logger/logger.go

import (
	"bytes"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// InitializeLogger configures the zap logger with JSON encoding and additional options.
func InitializeLogger() error {
	var err error
	config := zap.NewProductionConfig()
	config.Encoding = "json"
	config.OutputPaths = []string{"stdout"}

	// Set custom time format
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339) // Format to ISO 8601 date-time
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder              // Add short caller path (file and line number)

	// Build the logger with caller info enabled
	Logger, err = config.Build(zap.AddCaller())
	if err != nil {
		return err
	}
	return nil
}

// WithLogger returns a context with the logger attached for passing through the application.
func WithLogger(ctx context.Context) context.Context {
	return context.WithValue(ctx, "logger", Logger)
}

// GetLogger retrieves the logger from the context, or defaults to the global logger.
func GetLogger(ctx context.Context) *zap.Logger {
	if ctxLogger, ok := ctx.Value("logger").(*zap.Logger); ok {
		return ctxLogger
	}
	return Logger
}

type ResponseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

// Write overrides the Write method to capture the response body
func (w *ResponseWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)                  // Capture the response body
	return w.ResponseWriter.Write(b) // Write the actual response to the client
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		body := &bytes.Buffer{}
		rw := &ResponseWriter{
			ResponseWriter: c.Writer,
			Body:           body,
		}
		c.Writer = rw
		// Call the next handler

		// Retrieve the logger with the request context
		// log := GetLogger(c.Request.Context())
		ctx := WithLogger(c.Request.Context())
		// Attach the updated context to the request
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		log := GetLogger(c.Request.Context())

		// Log details of the request and response
		log.Info("API request",
			zap.String("url", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.Int("status", c.Writer.Status()),          // Response status
			zap.String("response_content", body.String()), // The body of the response (captured in the custom ResponseWriter)
			zap.String("response", c.Writer.Header().Get("Content-Type")),
		)
	}
}
