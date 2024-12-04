package helper

import "github.com/gin-gonic/gin"

func HandleError(c *gin.Context, status int, message string, err error) {
	if err != nil {
		c.JSON(status, gin.H{"details": err.Error(), "error": message})
	} else {
		c.JSON(status, gin.H{"error": message})
	}
}
