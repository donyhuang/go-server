package service

import "github.com/gin-gonic/gin"

func ReturnParamError(c *gin.Context) {
	ReturnError(c, -1, "param error")
}
func ReturnError(c *gin.Context, code int, msg string) {
	c.JSON(200, gin.H{"code": code, "msg": msg})
}
func ReturnSuccess(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "msg": "success"})
}
func ReturnData(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{"code": 0, "msg": "success", "data": data})

}

type handleFunc func(*gin.Context)

func WrapHandleFunc(h handleFunc) gin.HandlerFunc {
	return func(context *gin.Context) {
		h(context)
	}
}
