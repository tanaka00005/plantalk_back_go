package main

import (
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"github.com/tanaka00005/plantalk_back_go/login/token"
)

type User struct {
	Email string `json:"email" binding:"required"`
	Name string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func main(){
	r := gin.Default()

	r.GET("/",func (c *gin.Context)  {
		c.JSON(200,gin.H{
			"message":"Hello world",
		})
	})

	r.POST("/auth/register",func (c *gin.Context){
		var user User

		if err := c.ShouldBindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		// パスワードのハッシュ化
		//soltはランダムな文字
		hash,err := EncryptPassword("password")

		if err != nil{
			panic(err)
		}
		fmt.Println(hash)

		// compareResult := CompareHashAndPassword(hash,"password")

		// if compareResult{
		// 	fmt.Println("一致")
		// }else{
		// 	fmt.Println("不一致")
		// }

		accessToken,err := token.Token()

		if err != nil{
			panic(err)
		}

		c.JSON(http.StatusOK,accessToken)
	})

	r.POST("auth/login",func(c *gin.Context){
		var user User

		if err := c.ShouldBindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		
	})

	r.Run(":8080")

}

func EncryptPassword(password string) (string,error){
	hash,err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)

	return string(hash),err
}

func CompareHashAndPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))

	return err == nil
}



