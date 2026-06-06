/*
* エントリーポイント
 */

package main

import (
	"email-api/errors"
	"email-api/routers"
	"log"
	"net/http"
	"os"
	"time"

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

	e.Use(middleware.RequestID())     // リクエストごとの一意のIDを生成
	e.Use(middleware.RequestLogger()) // ロギング
	e.Use(middleware.Recover())       // パニック時のリカバリ
	e.Use(middleware.Gzip())          // Gzip圧縮

	// カスタムエラーハンドラを登録
	e.HTTPErrorHandler = errors.CustomHTTPErrorHandler

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://to-dev.jp"},
	}))

	e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      0.5,             // 平均 2秒に1リクエスト補充
				Burst:     5,               // 同一IPから短時間に最大5件まで
				ExpiresIn: 3 * time.Minute, // 3分アクセスのないIPは掃除
			},
		),
		IdentifierExtractor: func(c echo.Context) (string, error) {
			return c.RealIP(), nil
		},
		DenyHandler: func(c echo.Context, identifier string, err error) error {
			return echo.NewHTTPError(http.StatusTooManyRequests, "リクエストが多すぎます。しばらく待って再度お試しください。")
		},
	}))

	// ルーティング
	routers.SetupRouter(e, client, myEmail)

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "1323" // ローカル開発用のフォールバック
	}
	e.Logger.Fatal(e.Start(":" + port))
}
