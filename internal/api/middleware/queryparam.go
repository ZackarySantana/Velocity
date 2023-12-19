package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func QueryAmount(value int) gin.HandlerFunc {
	return func(c *gin.Context) {
		amountAsString := c.DefaultQuery("amount", strconv.Itoa(value))
		amount, err := strconv.Atoi(amountAsString)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"message": "amount must be a number",
			})
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

func InstanceId() gin.HandlerFunc {
	return func(c *gin.Context) {
		instanceId := c.Query("instance_id")
		if instanceId == "" {
			instanceId = c.Param("instance_id")
			if instanceId == "" {
				c.AbortWithStatusJSON(400, gin.H{
					"message": "instance_id required",
				})
				return
			}
		}
		// convert to objectid
		i, err := primitive.ObjectIDFromHex(instanceId)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"message": "invalid instance_id",
			})
			return
		}
		c.Set("instance_id", i)
		c.Next()
	}
}

func GetInstanceId(c *gin.Context) primitive.ObjectID {
	instanceId, exists := c.Get("instance_id")
	if !exists {
		return primitive.NilObjectID
	}
	return instanceId.(primitive.ObjectID)
}
