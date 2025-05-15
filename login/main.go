package main

import (
	"fmt"
	"net/http"
	"os"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tanaka00005/plantalk_back_go/login/token"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	Email string `json:"email" binding:"required" gorm:"uniqueIndex;size:255"`
	Name string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	ID	uint	`json:"id" gorm:"primaryKey;autoIncrement"`
}

func main(){


	err := godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
	}
	
	dbUsername := os.Getenv("DB_USERNAME")
	dbPort := os.Getenv("DB_PORT")
	dbScheema := os.Getenv("DB_SCHEEMA")
	dbPassword := os.Getenv("DB_PASSWORD")

	fmt.Printf("dbUsername:%s\n",dbUsername)
	fmt.Printf("dbPort:%s\n",dbPort)
	fmt.Printf("dbScheema:%s\n",dbScheema)
	fmt.Printf("dbPassword:%s\n",dbPassword)

	
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",dbUsername,dbPassword,dbPort,dbScheema)
	//dsn := "root:password@tcp(127.0.0.1:53306)/plantalk_go?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
		panic("failed to connect database")
	}
	//テーブルのマイグレーション
	db.AutoMigrate(&User{})
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
		hash,err := EncryptPassword(user.Password)

		if err != nil{
			panic(err)
		}
		fmt.Println(hash)

		user.Password = hash

		// compareResult := CompareHashAndPassword(hash,"password")

		// if compareResult{
		// 	fmt.Println("一致")
		// }else{
		// 	fmt.Println("不一致")
		// }

		accessToken,err := token.Token()

		result := db.Create(&user)

		if result.Error != nil {
			panic("failed to insert record")
		}

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






