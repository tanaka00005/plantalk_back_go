package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/tanaka00005/plantalk_back_go/middleware"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

type Question struct{
	Question string `json:"question" binding:"required"`
}

type User struct {
	Email string `json:"email" binding:"required" gorm:"uniqueIndex;size:255"`
	Name string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	ID	uint	`json:"id" gorm:"primaryKey;autoIncrement"`
	ChatLogs []ChatLog `json:"chat_logs" gorm:"foreignKey:UserID"`
	Plants []Plant `json:"plant" gorm:"foreignKey:UserID"`
}
type ChatLog struct{
	Message string `json:"message"`
	//Email string `json:"email"`
	IsAi bool `json:"is_ai"`
	UserID uint `json:"user_id" gorm:"not null"`
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`
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


func Chat(r *gin.Engine, db *gorm.DB){
	r.GET("/chat/test",func (c *gin.Context)  {
		c.JSON(200,gin.H{
			"message":"Hello world",
		})
	})

	r.POST("/chat/response",middleware.JWTAuthMiddleware(), func (c *gin.Context) {

		userEmail,exists := c.Get("user_email")

		if !exists {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"ユーザー情報の取得に失敗しました"})
			return 
		}

		var question Question

		if err := c.ShouldBindJSON(&question); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		fmt.Printf("res:%v\n",question.Question)

		var chatLog User
		
		fmt.Printf("userEmail:%v\n",userEmail)
		result := db.Where("email = ?",userEmail).Find(&chatLog)

		if result.Error != nil {
    		c.JSON(http.StatusInternalServerError, gin.H{"error": "データベースのクエリ実行に失敗しました。"})
    		return
		}

		var getUserID User
		result_plant := db.Preload("Plants").Where("email = ?",userEmail).First(&getUserID)

		if result_plant.Error != nil {
			// ユーザーが見つからない場合もエラーとして処理されます。
			log.Printf("Database query error: %v", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー情報の取得に失敗しました。"})
			return
		}

		fmt.Printf("email info:%v\n",chatLog)

		userID := chatLog.ID

		userMessageLog := ChatLog{
			Message: question.Question,
			IsAi: false, //user
			UserID: userID,
			//Email:userEmailStr,
		}
		fmt.Printf("userMessageLog:%v\n",userMessageLog)

		resultMessageLog := db.Create(&userMessageLog)

		if resultMessageLog.Error != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "データベースのクエリ実行に失敗しました。"})
    		return
		}

		ctx := context.Background()

		client, err := genai.NewClient(ctx,option.WithAPIKey(os.Getenv("API_KEY")))

		if err!= nil {
			log.Fatal(err)
		}
		defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	cs := model.StartChat()

	//履歴
	cs.History = []*genai.Content{
		&genai.Content{
			Parts: []genai.Part{
				genai.Text("how old are you"),
			},
			Role: "user",
		},
		&genai.Content{
			Parts: []genai.Part{
				genai.Text("Great to meet you. What would you like to know?"),
			},
			Role: "model",
		},
	}

	prompt := fmt.Sprintf(`
		あなたは、植物の育成状況を評価する専門家%sです。
		
		質問が来たら質問に簡潔に答えて。


		質問は%s。

		この質問に対して、適切に答えてください。

		では次のjson形式で返答してください。
		{
			"message":"質問に対する答えをここに書いてください"
		}
	`,getUserID.Plants[0].Species,question.Question)

	//質問
	response, err := cs.SendMessage(ctx,genai.Text(prompt+question.Question))

	if err != nil{
		fmt.Printf("Gemini APIエラー%v\n",err)
		c.JSON(http.StatusInternalServerError,gin.H{"error":"AIの応答取得に失敗しました。"})
		return 
	}

	//レスポンステキストを取得
	responseText := ""
	for _,part := range response.Candidates[0].Content.Parts{
		if textPart, ok := part.(genai.Text); ok{
			responseText += string(textPart)
		}
	}

	var aiResponse ChatLog

	cleanTextFront := strings.ReplaceAll(responseText,"```json","")
	cleanTextEnd := strings.ReplaceAll(cleanTextFront,"```","")
	cleanText := strings.TrimSpace(cleanTextEnd)
	
	if err := json.Unmarshal([]byte(cleanText), &aiResponse); err != nil {
		fmt.Printf("JSONパースエラー: %v\n", err)
		// パースに失敗した場合は生のテキストを使用
		aiResponse.Message = cleanText
	}

	aiResponseLog := ChatLog{
			Message: aiResponse.Message,
			IsAi: true, //ai
			UserID: userID,
			//Email:userEmailStr,
		}
		fmt.Printf("aiResponse:%v\n",aiResponseLog)

		resultAiResponse := db.Create(&aiResponseLog)

		if resultAiResponse.Error != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "データベースのクエリ実行に失敗しました。"})
    		return
		}

	c.JSON(http.StatusOK,aiResponseLog)
	
	})

	r.GET("/chat/history",middleware.JWTAuthMiddleware(),func (c *gin.Context)  {
		userEmail,exists := c.Get("user_email")

		if !exists {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"ユーザー情報の取得に失敗しました"})
			return 
		}

		fmt.Printf("userEmainnmml:%v\n",userEmail)

		var getUserID User

		result := db.Preload("Plants").Where("email = ?",userEmail).First(&getUserID)

		if result.Error != nil {
    		c.JSON(http.StatusInternalServerError, gin.H{"error": "データベースのクエリ実行に失敗しました。"})
    		return
		}

		var ChatLogMessage []ChatLog

		resultMessage := db.Where("user_id = ?",getUserID.ID).Find(&ChatLogMessage)

		if resultMessage.Error != nil {
    		c.JSON(http.StatusInternalServerError, gin.H{"error": "データベースのクエリ実行に失敗しました。"})
    		return
		}

		c.JSON(http.StatusOK, gin.H{
    		"plant_species_name": getUserID.Plants[0].Species,
   			"chat_logs":          ChatLogMessage,
		})

	})
	
}