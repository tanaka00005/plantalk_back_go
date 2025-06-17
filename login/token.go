package login

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)


func Token(user string) (string,error){

	fmt.Println("id:",user)
	//jwtに付与する構造体

	//ペイロード(中身)の作成
	claims := jwt.MapClaims{
		"user_email":user,
		"exp":time.Now().Add(time.Hour * 72).Unix(),
	}

	//ヘッダーとペイロード生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	fmt.Println("token1:",token)
	//トークンに著者を付与
	accessToken,_ := token.SignedString([]byte("ACCESS_SECRET_KEY"))
	fmt.Println("accessToken:",accessToken)
	
	return accessToken,nil
}

 
