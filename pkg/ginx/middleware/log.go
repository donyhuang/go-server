package middleware

import (
	"fmt"
	log2 "gitlab.nongchangshijie.com/go-base/server/pkg/log"
	"gitlab.nongchangshijie.com/go-base/server/pkg/trace"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	TraceKey = "TraceId"
)

func LoggerMiddleware(c *gin.Context) {
	startTime := time.Now()
	traceId := trace.IdFromContext(c)
	c.Set(TraceKey, traceId)
	log2.InjectLogEntryToGinContext(traceId, c)
	c.Next()
	stopTime := time.Since(startTime)
	spendTime := fmt.Sprintf("%d ms", int(math.Ceil(float64(stopTime.Nanoseconds()/1000000))))
	statusCode := c.Writer.Status()
	dataSize := c.Writer.Size()
	if dataSize < 0 {
		dataSize = 0
	}
	method := c.Request.Method
	url := c.Request.RequestURI
	Log := log2.Logger.WithFields(logrus.Fields{
		"SpendTime": spendTime,
		"path":      url,
		"Method":    method,
		"status":    statusCode,
		"TraceId":   traceId,
	})
	if len(c.Errors) > 0 {
		Log.Error(c.Errors.ByType(gin.ErrorTypePrivate))
	}
	if statusCode >= 500 {
		Log.Error()
	} else if statusCode >= 400 {
		Log.Warn()
	} else {
		Log.Info()
	}

}
