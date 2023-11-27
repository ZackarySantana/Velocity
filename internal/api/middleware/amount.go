package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func QueryAmount(value int) gin.HandlerFunc {
	return func(c *gin.Context) {
		amountAsString := c.DefaultQuery("amount", strconv.Itoa(value))
		amount, err := strconv.Atoi(amountAsString)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "amount must be a number",
			})
			c.Abort()
			return
		}
		c.Set("amount", amount)
		c.Next()
	}
}

func GetQueryAmount(c *gin.Context) int {
	amount, exists := c.Get("amount")
	if !exists {
		return -1
	}
	return amount.(int)
}
