/*
* エントリーポイント
 */

package main

import (
	"email-api/errors"
	"email-api/routers"
	"log"
	"os"

	"github.com/resend/resend-go/v3"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomValidator struct {
    validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
    return cv.validator.Struct(i)
}

func main() {
	 if os.Getenv("APP_ENV") != "production" {
        if err := godotenv.Load(); err != nil {
            log.Println("Warning: .env file not found")
        }
    }

    apiKey := os.Getenv("API_KEY")
    if apiKey == "" {
        log.Fatal("API_KEY が設定されていません")
    }
    myEmail := os.Getenv("EMAIL_ADDRESS")
    if myEmail == "" {
        log.Fatal("EMAIL_ADDRESS が設定されていません")
    }
	client := resend.NewClient(apiKey)

	// Echoのインスタンス作成
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()} 

	e.Use(middleware.RequestID()) // リクエストごとの一意のIDを生成
	e.Use(middleware.RequestLogger())    // ロギング
	e.Use(middleware.Recover())   // パニック時のリカバリ
	e.Use(middleware.Gzip())      // Gzip圧縮

	// カスタムエラーハンドラを登録
	e.HTTPErrorHandler = errors.CustomHTTPErrorHandler

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    	AllowOrigins: []string{"https://to-dev.jp"},
	}))

	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5)))

	// ルーティング
	routers.SetupRouter(e, client, myEmail)

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "1323" // ローカル開発用のフォールバック
	}
	e.Logger.Fatal(e.Start(":" + port))
}
