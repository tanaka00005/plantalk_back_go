package calendar

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tanaka00005/plantalk_back_go/middleware"
	"gorm.io/gorm"
)
type User struct {
	Email string `json:"email" binding:"required" gorm:"uniqueIndex;size:255"`
	Name string `json:"name" binding:"required"`
	ID	uint	`json:"id" gorm:"primaryKey;autoIncrement"`
	Diaries []Diary `json:"diaries" gorm:"foreignKey:UserID"`
}
type Diary struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"not null"`       
	Content     string    `json:"content"`                       
	HealthState int       `json:"health_state"`               
	GrowthState int       `json:"growth_state"`                  
	RecordedAt  time.Time `json:"recorded_at" gorm:"type:timestamp; default:CURRENT_TIMESTAMP"` 
}

func Calendar(r *gin.Engine, db *gorm.DB){
	r.GET("/calendar/event",middleware.JWTAuthMiddleware(), func (c *gin.Context) {

		userEmail,exists := c.Get("user_email")

		if !exists {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"ユーザー情報の取得に失敗しました"})
			return 
		}

		fmt.Printf("userEmainnmml:%v\n",userEmail)

		yearStr := c.Query("year")
        monthStr := c.Query("month")

        // パラメータが指定されていない場合は、現在の年月をデフォルト値とする
        now := time.Now()
        year, err := strconv.Atoi(yearStr)
        if err != nil {
            year = now.Year()
        }
        month, err := strconv.Atoi(monthStr)
        if err != nil {
            month = int(now.Month())
        }


		// 3. 取得したい月の開始日時と終了日時を計算
		// 例: 7月の場合、7月1日 00:00:00 から 8月1日 00:00:00 の直前まで
		startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
		endDate := startDate.AddDate(0, 1, 0)

		fmt.Printf("startDate:%v\n",startDate)
		fmt.Printf("endDate:%v\n",endDate)


		// 4. GORMのPreloadに条件を追加して、指定した範囲の日記のみを読み込む
		var user User
		result := db.Preload("Diaries", "recorded_at >= ? AND recorded_at < ?", startDate, endDate).
			Where("email = ?", userEmail).First(&user)

		if result.Error != nil {
    		if result.Error == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "指定されたユーザーが見つかりません"})
				return
			}
			// その他のデータベースエラー
			c.JSON(http.StatusInternalServerError, gin.H{"error": "データベースエラーが発生しました"})
			return
		}
		fmt.Printf("user:%v\n",user)

		
		c.JSON(http.StatusOK,user.Diaries)

	})

	r.GET("/calendar/diary",middleware.JWTAuthMiddleware(), func (c *gin.Context) {
		userEmail,exists := c.Get("user_email")

		if !exists {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"ユーザー情報の取得に失敗しました"})
			return 
		}

		fmt.Printf("userEmainnmml:%v\n",userEmail)

		var user User

		result := db.Where("email = ?",userEmail).First(&user)

		if result.Error != nil {
    		if result.Error == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "指定されたユーザーが見つかりません"})
				return
			}
			// その他のデータベースエラー
			c.JSON(http.StatusInternalServerError, gin.H{"error": "データベースエラーが発生しました"})
			return
		}

		nowYear := 2025;
  		month := 5;
  		date := 4;
  		
		saveData := time.Date(nowYear, time.Month(month), date+1, 0, 0, 0, 0, time.Local)

		fmt.Println(saveData)


	})
}