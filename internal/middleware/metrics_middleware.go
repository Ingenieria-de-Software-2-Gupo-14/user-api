package middleware

import (
	"strconv"
	"time"

	"github.com/Ingenieria-de-Software-2-Gupo-14/go-core/pkg/telemetry"
	"github.com/gin-gonic/gin"
)

func NewTelemetryStack(client telemetry.Client) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		telemetry.AddTracer(client),
		MetricsMiddleware,
	}
}

func MetricsMiddleware(c *gin.Context) {
	startTime := time.Now()

	tags := []string{
		"method:" + c.Request.Method,
		"path:" + c.Request.URL.Path,
	}

	ctx := c.Request.Context()

	telemetry.Incr(ctx, "traffic.request", tags...)

	c.Next()

	tags = append(tags, "status:"+strconv.Itoa(c.Writer.Status()))
	telemetry.Timing(ctx, "traffic.request.time", time.Since(startTime), tags...)
	telemetry.Incr(ctx, "traffic.response", tags...)
}
