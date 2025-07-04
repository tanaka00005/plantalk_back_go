package chat

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"github.com/tanaka00005/plantalk_back_go/middleware"
)

type Question struct{
	Question string `json:"question" binding:"required"`
}

type Response struct{
	Message string `json:"message"`
	Email string `json:"email"`
}

func Chat(r *gin.Engine){
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

	species := "ラディッシュ"
	sunlight := 100
	water := 100
	temperature := 20
	humidity := 50
	character := "モスモス"

	prompt := fmt.Sprintf(`
		あなたは、植物の育成状況を評価する専門家%sです。%sとして植物の状態を1行でしゃべってください。
		語尾に「モス」をつけてください。

		植物の種類は%sです。
		日光は%d。
		水分は%d。
		温度は%d。
		湿度は%d。

		質問は%s。

		この質問に対して、モスモスとして適切に答えてください。植物に関する質問の場合は上記の情報も参考に考えてください。植物に関係ない質問でもモスモスとして答えてください。

		では次のjson形式で返答してください。
		{
			"message":"質問に対する答えをここに書いてください"
		}
	`,character,character,species,sunlight,water,	temperature,humidity,question.Question)

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

	fmt.Printf("responseText:%v",responseText)

	

	//jsonの開始と終了を見つけて摘出
	var aiResponse Response

	cleanTextFront := strings.ReplaceAll(responseText,"```json","")
	cleanTextEnd := strings.ReplaceAll(cleanTextFront,"```","")
	cleanText := strings.TrimSpace(cleanTextEnd)
	
	// JSONとして解析を試行

	//[]byte(cleanText)でcleanTextをバイトの配列に
	//&aiResponseでresponseにmessageのメモリアドレスを取得
	//unmarshlでjsonパーサーがバイト配列を解析
	//"message"キーを発見
	//cleanTextの値を取得
	//aiResponse.MessageフィールドにcleanTextの値を代入
	if err := json.Unmarshal([]byte(cleanText), &aiResponse); err != nil {
		fmt.Printf("JSONパースエラー: %v\n", err)
		// パースに失敗した場合は生のテキストを使用
		aiResponse.Message = cleanText
	}

	aiResponse.Email = userEmail.(string)

	fmt.Printf("prompt:%v\n",question.Question)
	fmt.Printf("response:%v\n",response)
	fmt.Printf("aiResponse:%v\n",aiResponse)



	c.JSON(http.StatusOK,aiResponse)
	
	})
	
}