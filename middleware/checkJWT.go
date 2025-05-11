package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)


func main(){
	//jwtに付与する構造体
	claims := jwt.MapClaims{
		"user_id":"user_id1234",
		"exp":time.Now().Add(time.Hour * 72).Unix(),
	}

	//ヘッダーとペイロード生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	//トークンに著者を付与
	accessToken,_ := token.SignedString([]byte("ACCESS_SECRET_KEY"))
	fmt.Println("accessToken:",accessToken)
}