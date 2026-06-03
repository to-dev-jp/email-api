/*
* エントリーポイント
 */

package main

import (
	"email-api/errors"
	"email-api/routers"
	"fmt"
	"os"

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
        err := godotenv.Load()
        if err != nil {
            fmt.Println("Warning: .env file not found")
        }
	}
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
	routers.SetupRouter(e)

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "1323" // ローカル開発用のフォールバック
	}
	e.Logger.Fatal(e.Start(":" + port))
}
