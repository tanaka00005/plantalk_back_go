package check

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Check(r *gin.Engine) {
	
	r.GET("/check/private",func(c *gin.Context){
		tokenString := c.GetHeader("x-auth-token")
		fmt.Printf("tokenString:%v\n",tokenString)

		if tokenString == "" {
			c.JSON(400, gin.H{"error": "権限がありません"})
			c.Abort()
			return 
		}

		//受け取ったjwtトークンが本当にサーバーが発行したもので改竄されていないかを検証
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		fmt.Printf("token:%v\n",token)
		return []byte("ACCESS_SECRET_KEY"), nil
		
	})

	if err != nil || token == nil {
		fmt.Println("トークン解析エラー:",err)
		return
	}

	claims,ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid{
		fmt.Println("トークン無効またはclaims変換失敗")
		return
	}

	userEmail,ok := claims["user_email"].(string)

	if !ok{
		fmt.Println("user_idの方が不正")
		return
	}
	fmt.Println("user_email:", userEmail)

	if exp,ok := claims["exp"].(float64); ok{
		fmt.Println("exp:", int64(exp))
	}else{
		fmt.Println("expの方が不正")
	}

		fmt.Printf("token:%v",tokenString)
	})
}
