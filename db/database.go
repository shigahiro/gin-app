package db

import (
	"os"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/shigahiro/gin-app/model"
	"golang.org/x/crypto/bcrypt"
)

func Init() {
	db := gormConnect()

	defer db.Close()
	db.AutoMigrate(&model.Tweet{}) //構造体に基づいてテーブルを作成
	db.AutoMigrate(&model.User{})
}

func gormConnect() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
	DBMS := os.Getenv("mytweet_DBMS")
	USER := os.Getenv("mytweet_USER")
	PASS := os.Getenv("mytweet_PASS")
	DBNAME := os.Getenv("mytweet_DBNAME")

	// time.Timeの処理のために?parseTime=trueを追加
	// ホスト名を指定する際はdbnameの前に()内で指定
	CONNECT := USER + ":" + PASS + "@/" + DBNAME + "?parseTime=true"

	db, err := gorm.Open(DBMS, CONNECT)

	// データベースと接続できないとプログラムの実行が危ういからパニック
	if err != nil {
		panic(err.Error())
	}
	return db
}

func Insert(content string) {
	db := gormConnect()

	defer db.Close()
	db.Create(&model.Tweet{Content: content})
}

func Update(id int, tweetText string) {
	db := gormConnect()
	var tweet model.Tweet
	db.First(&tweet, id)
	tweet.Content = tweetText
	db.Save(&tweet)
	db.Close()
}

func GetAll() []model.Tweet {
	db := gormConnect()

	defer db.Close()
	var tweets []model.Tweet
	// tweetsに登録順に並び替えたものを入れる
	db.Order("created_at desc").Find(&tweets)
	return tweets
}

func GetOne(id int) model.Tweet {
	db := gormConnect()
	var tweet model.Tweet
	db.First(&tweet, id)
	db.Close()
	return tweet
}

func Delete(id int) {
	db := gormConnect()
	var tweet model.Tweet
	db.First(&tweet, id)
	db.Delete(&tweet)
	db.Close()
}

// 以下ユーザの処理
func CreateUser(username string, password string) error {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	db := gormConnect()
	defer db.Close()
	// Insert処理
	if err := db.Create(&model.User{Username: username, Password: string(hashed)}).Error; err != nil {
		return err
	}
	return nil
}

func GetUser(username string) model.User {
	var user model.User

	db := gormConnect()
	db.Where("username = ?", username).First(&user)
	return user
}
