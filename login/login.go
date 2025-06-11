package login

import (
	"fmt"
	"net/http"
	"errors"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Email string `json:"email" binding:"required" gorm:"uniqueIndex;size:255"`
	Name string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	ID	uint	`json:"id" gorm:"primaryKey;autoIncrement"`
}

func Login(r *gin.Engine, db *gorm.DB){

	solt := "1234567890"
	

	r.POST("/auth/register",func (c *gin.Context){
		var user User

		if err := c.ShouldBindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		// var dbUer User

		// //userにdbから取ってきた一番最初の値を上書き
		// result := db.First(&dbUer,"email = ?",user.Email)

		var dbUer User

		data := db.First(&dbUer,"email = ?",user.Email)

		fmt.Printf("dbdata:%v",user)

		if(data.Error == nil){
			c.JSON(http.StatusConflict,gin.H{"error":"このメールアドレスはすでに存在しています"})
		}else if !errors.Is(data.Error, gorm.ErrRecordNotFound) {
			// その他のDBエラー
			c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラー"})
			return
		}

		// パスワードのハッシュ化
		//soltはランダムな文字
		hash,err := EncryptPassword(user.Password+solt)

		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"サーバーエラー"})
			return 
		
		}
		fmt.Println(hash)

		user.Password = hash
		accessToken,err := Token()

		result := db.Create(&user)

		if result.Error != nil {
			panic("failed to insert record")
		}
		

		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"サーバーエラー"})
			return 
		}

		c.JSON(http.StatusOK,accessToken)
	})

	r.POST("auth/login",func(c *gin.Context){
		var user User

		//userに取得してきた値を上書き
		if err := c.ShouldBindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		fmt.Printf("first:%v",user)

		var dbUer User

		//userにdbから取ってきた一番最初の値を上書き
		result := db.First(&dbUer,"email = ?",user.Email)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound,gin.H{"error":"User not found"})
			return 
		}

		
		if(user.Name != dbUer.Name){
			c.JSON(http.StatusUnauthorized,gin.H{"message":"ユーザー名が一致しません"})
		}
		
		compareResult := CompareHashAndPassword(dbUer.Password,user.Password+solt)

		if compareResult{
			fmt.Println("一致")
		}else{
			fmt.Println("不一致")
		}

		accessToken,err := Token()

		if result.Error != nil {
			panic("failed to insert record")
		}

		if err != nil{
			panic(err)
		}



		c.JSON(http.StatusOK,accessToken)

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




