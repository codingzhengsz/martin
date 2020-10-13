package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	router.Use(Recover)

	router.GET("/ping/fail", func(c *gin.Context) {
		var slice = []int{1, 2, 3}
		slice[4] = 4
	})

	router.GET("/ping/success", successPingHandler)

	_ = router.Run(":8081")
}

func successPingHandler(c *gin.Context) {
	actionId := c.GetHeader("actionId")

	println(actionId)

	c.Header("actionId", "Return ActionId: " + actionId)
	c.String(http.StatusOK, "Success")
}
