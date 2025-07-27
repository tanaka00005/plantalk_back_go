package calendar

import (
	"fmt"
	"net/http"
	"time"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/tanaka00005/plantalk_back_go/middleware"
	"gorm.io/gorm"
)

type Plant struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`
	Species string `json:"species" binding:"required"`
	SpeciesName string `json:"speciesname" binding:"required"`
	UserID uint `json:"user_id" gorm:"not null"`
	Diaries []Diary `json:"diaries" gorm:"foreignKey:PlantID"`
}
type User struct {
	Email string `json:"email" binding:"required" gorm:"uniqueIndex;size:255"`
	Name string `json:"name" binding:"required"`
	ID	uint	`json:"id" gorm:"primaryKey;autoIncrement"`
	Diaries []Diary `json:"diaries" gorm:"foreignKey:UserID"`
	Plants []Plant `json:"plant" gorm:"foreignKey:UserID"`
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
type DiaryInput struct {
    Content     string `json:"content"`
    Month       int    `json:"month"`
    Day         int    `json:"day"`
	HealthState int       `json:"health_state"`               
	GrowthState int       `json:"growth_state"`                  
}

func Calendar(r *gin.Engine, db *gorm.DB){
	r.GET("/calendar/event",middleware.JWTAuthMiddleware(), func (c *gin.Context) {

		userEmail,exists := c.Get("user_email")

		if !exists {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"ユーザー情報の取得に失敗しました"})
			return 
		}

		fmt.Printf("userEmainnmml:%v\n",userEmail)

		// 4. GORMのPreloadに条件を追加して、指定した範囲の日記のみを読み込む
		var user User
		result := db.Preload("Diaries").Where("email = ?", userEmail).First(&user)

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

	r.POST("/calendar/post-diary",middleware.JWTAuthMiddleware(), func (c *gin.Context) {
		userEmail,exists := c.Get("user_email")

		if !exists {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"ユーザー情報の取得に失敗しました"})
			return 
		}

		fmt.Printf("userEmainnmml:%v\n",userEmail)

		var user User
    	if err := db.Preload("Plants").Where("email = ?", userEmail).First(&user).Error; err != nil {
        	if errors.Is(err, gorm.ErrRecordNotFound) {
            	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            	return
        		}
        	c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error while fetching user"})
       	 	return
   		 }


		var stateLog DiaryInput

		if err := c.ShouldBindJSON(&stateLog); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		fmt.Printf("res:%v\n",stateLog)
		
		// 現在時刻を取得
		now := time.Now()
		nowYear := now.Year()
		saveDate := time.Date(nowYear, time.Month(stateLog.Month), stateLog.Day, 0, 0, 0, 0, time.Local)

		fmt.Printf("saveDate:%v\n",saveDate)
		
		var getState Diary
		startOfDay := saveDate
        endOfDay := startOfDay.AddDate(0, 0, 1) // 翌日の0時
        eventHistory := db.Where("user_id = ? AND recorded_at >= ? AND recorded_at < ?", user.ID, startOfDay, endOfDay).First(&getState)

		if eventHistory.Error == nil{
			updateData := Diary{
                HealthState: getState.HealthState,
                GrowthState: getState.GrowthState,
                Content:     getState.Content,
            }
            if err := db.Model(&eventHistory).Updates(updateData).Error; err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update diary"})
                return
            }
		}else if errors.Is(eventHistory.Error, gorm.ErrRecordNotFound) {
			if len(user.Plants) == 0 {
                c.JSON(http.StatusBadRequest, gin.H{"error": "日記を紐付ける植物がありません"})
                return
            }
			newData := Diary{
				UserID:     user.ID,
				HealthState: stateLog.HealthState,
				GrowthState: stateLog.GrowthState,
				Content:     stateLog.Content,
				RecordedAt:  saveDate,
				PlantID: user.Plants[0].ID,
			}
			if err := db.Create(&newData).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create diary"})
                return
			}
		}

		var updatedUser User
        db.Preload("Plants").Preload("Diaries").Where("email = ?", userEmail).First(&updatedUser)

		c.JSON(http.StatusOK,updatedUser) 
	})
}