package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tanaka00005/plantalk_back_go/chat"
	"github.com/tanaka00005/plantalk_back_go/login"
	"github.com/tanaka00005/plantalk_back_go/calendar"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)
type Question struct{
	Question string `json:"question" binding:"required"`
}


type ChatLog struct{
	Message string `json:"message"`
	//Email string `json:"email"`
	IsAi bool `json:"is_ai"`
	UserID uint `json:"user_id" gorm:"not null"`
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`
}

type User struct {
	Email string `json:"email" binding:"required" gorm:"uniqueIndex;size:255"`
	Name string `json:"name" binding:"required"`
	ID	uint	`json:"id" gorm:"primaryKey;autoIncrement"`
	Diaries []Diary `json:"diaries" gorm:"foreignKey:UserID"`
	Plants []Plant `json:"plant" gorm:"foreignKey:UserID"`
	ChatLogs []ChatLog `json:"chat_logs" gorm:"foreignKey:UserID"`
}

type Plant struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`
	Species string `json:"species" binding:"required"`
	SpeciesName string `json:"speciesname" binding:"required"`
	UserID uint `json:"user_id" gorm:"not null"`
	Diaries []Diary `json:"diaries" gorm:"foreignKey:PlantID"`
}

type Diary struct {
	ID          uint      `gorm:"primaryKey"`
	PlantID     uint      `json:"plant_id" gorm:"not null"`    
	UserID      uint      `json:"user_id" gorm:"not null"`       
	Content     string    `json:"content"`                       
	HealthState int       `json:"health_state"`               
	GrowthState int       `json:"growth_state"`                  
	RecordedAt  time.Time `json:"recorded_at" gorm:"type:timestamp; default:CURRENT_TIMESTAMP"`
}


func main(){

	err := godotenv.Load()
	// if err != nil {
	//   log.Fatal("Error loading .env file")
	// }
	
	dbUer := os.Getenv("DB_USERNAME")
	dbPort := os.Getenv("DB_PORT")
	dbScheema := os.Getenv("DB_SCHEEMA")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")

	
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",dbUer,dbPassword,dbHost,dbPort,dbScheema)
	//dsn := "root:password@tcp(127.0.0.1:53306)/plantalk_go?charset=utf8mb4&parseTime=True&loc=Local"

	var db *gorm.DB

    // ----- 👇ここからが重要！このループ処理を追加・修正してください -----
    // 10回まで再試行する
    for i := 0; i < 10; i++ {
        db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
        if err == nil {
            // 接続に成功したらループを抜ける
            log.Println("🎉 データベース接続成功！")
            break
        }
        log.Printf("DB接続に失敗しました。2秒後に再試行します... (%d/10)", i+1)
        time.Sleep(2 * time.Second)
    }

    // 10回試行してもダメだったら、プログラムを終了
    if err != nil {
        log.Fatalf("💀 データベースに接続できませんでした: %v", err)
    }

	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
		panic("failed to connect database")
	}

	//テーブルのマイグレーション
	err = db.AutoMigrate(&User{}, &Plant{}, &Diary{}, &ChatLog{})
	if err != nil {
    	panic("failed to migrate database")
	}
	
	r := gin.Default()

	r.GET("/",func (c *gin.Context)  {
		c.JSON(200,gin.H{
			"message":"Hello world",
		})
	})

	config := cors.DefaultConfig()
	config.AllowOriginFunc = func(origin string) bool {
        return origin == "http://localhost:5173"
    }

	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "x-auth-token"}
    config.AllowCredentials = true
	
	r.Use(cors.New(config))

	login.Login(r, db)
	chat.Chat(r,db)
	calendar.Calendar(r,db)
	r.Run(":8080")

}

// func corsMiddleware(allowOrigins []string) gin.HandlerFunc {
// 	config := cors.DefaultConfig()
// 	config.AllowOrigins = allowOrigins
// 	return cors.New(config)
// }
