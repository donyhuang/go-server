package ginx

import (
	"github.com/gin-gonic/gin"
)

const HealthCheck = "/health"

func AddHealthCheck(r *gin.Engine) {
	r.GET(HealthCheck, func(context *gin.Context) {})
}
