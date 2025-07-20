package login

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	
)

type User struct {
	Email string `json:"email" binding:"required" gorm:"uniqueIndex;size:255"`
	Name string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	ID	uint	`json:"id" gorm:"primaryKey;autoIncrement"`
	Plant []Plant `json:"plant" gorm:"foreignKey:UserID"`
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

type ChatLog struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null"` 
	Message    string    `json:"message"`                 
	Sender     string    `json:"sender"`                  
	RecordedAt time.Time `json:"recorded_at" gorm:"type:timestamp; default:CURRENT_TIMESTAMP"`
}


func Login(r *gin.Engine, db *gorm.DB){

	solt := "1234567890"
	
	r.GET("/test",func (c *gin.Context)  {
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

		// var dbUer User

		// //userにdbから取ってきた一番最初の値を上書き
		// result := db.First(&dbUer,"email = ?",user.Email)

		var dbUser User
		fmt.Printf("user:%v\n",user)
		fmt.Printf("dbuser:%v\n",dbUser)

		//dbUserは空
		//user.Emailは入力したemail
		//データベースからuser.Emailと同じemailを探して挿入
		data := db.First(&dbUser,"email = ?",user.Email)

		fmt.Printf("dbdata:%v",dbUser)

		//もしデータベースとuser.Emailに同じemailが見つかったら->エラーが起こらなかったら
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


		accessToken,err := Token(user.Email)

		result := db.Create(&user)

		if result.Error != nil {
			panic("failed to insert record")
		}

		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"サーバーエラー"})
			return 
		}

		c.JSON(http.StatusOK,gin.H{"token":accessToken})
	})

	r.POST("/auth/login",func(c *gin.Context){
		var user User

		//userに取得してきた値を上書き
		if err := c.ShouldBindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		fmt.Printf("first:%v",user)

		var dbUser User

		//userにdbから取ってきた一番最初の値を上書き
		result := db.First(&dbUser,"email = ?",user.Email)

		fmt.Printf("dbUser:%v",dbUser)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound,gin.H{"error":"User not found"})
			return 
		}

		
		if(user.Name != dbUser.Name){
			c.JSON(http.StatusUnauthorized,gin.H{"message":"ユーザー名が一致しません"})
		}
		
		compareResult := CompareHashAndPassword(dbUser.Password,user.Password+solt)

		if compareResult{
			fmt.Println("一致")
		}else{
			fmt.Println("不一致")
		}


		accessToken,err := Token(user.Email)

		if result.Error != nil {
			panic("failed to insert record")
		}

		if err != nil{
			panic(err)
		}

		c.JSON(http.StatusOK,gin.H{"token":accessToken})

	})


}


func EncryptPassword(password string) (string,error){
	hash,err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)

	return string(hash),err
}

func CompareHashAndPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))

	return err == nil
}




