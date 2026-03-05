package pkg

import "github.com/gin-gonic/gin"

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Success(c *gin.Context, statusCode int, data interface{}){
	c.JSON(statusCode, Response{
		Success: true,
		Data: data,
	})
}

func Error(c *gin.Context, statusCode int, message string){
	c.JSON(statusCode, Response{
		Success: false,
		Error: message,
	})
}