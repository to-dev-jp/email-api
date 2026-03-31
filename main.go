/*
* エントリーポイント
 */

package main

import (
	"email-api/errors"
	"email-api/routers"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echoのインスタンス作成
	e := echo.New()

	e.Use(middleware.RequestID()) // リクエストごとの一意のIDを生成
	e.Use(middleware.RequestLogger())    // ロギング
	e.Use(middleware.Recover())   // パニック時のリカバリ
	e.Use(middleware.Gzip())      // Gzip圧縮

	// カスタムエラーハンドラを登録
	e.HTTPErrorHandler = errors.CustomHTTPErrorHandler

	// ルーティング
	routers.SetupRouter(e)

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "1323" // ローカル開発用のフォールバック
	}
	e.Logger.Fatal(e.Start(":" + port))
}
