package helper

import "github.com/gin-gonic/gin"

func SendResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	obj := map[string]interface{}{
		"statusCode": statusCode,
		"message":    message,
		"data":       data,
	}
	c.JSON(statusCode, obj)
}
