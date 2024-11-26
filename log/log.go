package Log

import (
	"bytes"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Log struct {
	logger *zap.Logger
}

func (l *Log) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *Log) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *Log) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}
func (l *Log) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *Log) With(fields ...zap.Field) *Log {
	return &Log{logger: l.logger.With(fields...)}
}

var LogInstance *Log

// Init Logger config the zap logger with JSON encoding and additional options
func InitializeLogger() error {
	config := zap.NewProductionConfig()
	config.Encoding = "json"
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	logger, err := config.Build(zap.AddCaller())
	if err != nil {
		return err
	}
	LogInstance = &Log{logger: logger}
	return nil

}

func GetLogger(c *gin.Context) *Log {
	log, exists := c.Get("log")
	if exists {
		if logger, ok := log.(*Log); ok {
			return logger
		}
	}
	// Fallback to global logger if not found
	return LogInstance
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
		logger := LogInstance.With(zap.String("request_id", requestID))

		// Attach logger to request context
		c.Set("log", logger)

		// Capture response
		rw := NewResponseWriter(c.Writer)
		c.Writer = rw

		start := time.Now()
		c.Next()
		duration := time.Since(start)

		// Log request and response details
		logger.Info("API request",
			zap.String("url", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.Int("status", c.Writer.Status()),
			zap.String("response_content", rw.body.String()),
			zap.String("response", c.Writer.Header().Get("Content-Type")),
			zap.Duration("duration", duration),
		)
	}
}
