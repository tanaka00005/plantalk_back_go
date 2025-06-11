package login

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)


func Token() (string,error){
	//jwtに付与する構造体

	//ペイロード部分に含まれる情報を定義するためのマップ形式の型
	claims := jwt.MapClaims{
		"user_id":"user_id1234",
		"exp":time.Now().Add(time.Hour * 72).Unix(),
	}

	//ヘッダーとペイロード生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	fmt.Println("token:",token)
	//トークンに著者を付与
	accessToken,_ := token.SignedString([]byte("ACCESS_SECRET_KEY"))
	fmt.Println("accessToken:",accessToken)
	
	return accessToken,nil
}

 
