package check

import(
	"fmt"
	
	"github.com/gin-gonic/gin"
)

func Check(r *gin.Engine) {

	r.GET("/check/private",func(c *gin.Context){
		token := c.GetHeader("x-auth-token")

		if token == "" {
			c.JSON(401, gin.H{"error": "Token is missing"})
			c.Abort()
			return 
		}

		fmt.Printf("token:%v",token)
	})
}
