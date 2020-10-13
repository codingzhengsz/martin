package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic: %v\n", r)
			c.JSON(http.StatusBadRequest, gin.H{
				"code": "1",
				"msg":  errToString(r),
				"data": nil,
			})
			c.Abort()
		}
	}()
	c.Next()
}

func errToString(r interface{}) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	default:
		return r.(string)
	}
}
