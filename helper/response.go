package helper

import "github.com/gin-gonic/gin"

func HandleResponse(c *gin.Context, status int, message string, detail any) {
	if detail != nil {
		c.JSON(status, gin.H{"message": message, "details": detail})
	} else {
		c.JSON(status, gin.H{"message": message})
	}

}
