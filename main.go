package main

import (
	"fmt"
	"os"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/tanaka00005/plantalk_back_go/login/login.go"
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
	
	dbUer := os.Getenv("DB_USERNAME")
	dbPort := os.Getenv("DB_PORT")
	dbScheema := os.Getenv("DB_SCHEEMA")
	dbPassword := os.Getenv("DB_PASSWORD")

	
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",dbUer,dbPassword,dbPort,dbScheema)
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

	login.Login()
	

	r.Run(":8080")

}
