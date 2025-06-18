package chat

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)
func Chat(r *gin.Engine){
	r.GET("/chat/test",func (c *gin.Context)  {
		c.JSON(200,gin.H{
			"message":"Hello world",
		})
	})

	r.GET("/chat/response",func (c *gin.Context) {
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

	//質問
	resp, err := cs.SendMessage(ctx,genai.Text("How many paws are in my house?"))

	if err != nil{
		log.Fatal(err)
	}
	c.JSON(http.StatusOK,resp)
	})
	
}