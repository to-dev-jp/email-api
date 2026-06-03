/*
* ルーティングの設定を行う
 */

package routers

import (
	"email-api/controllers"

	"github.com/labstack/echo/v4"
)

// ルーティングの設定
func SetupRouter(e *echo.Echo) {
	// コントローラを初期化
	// Repositoryをコントローラに注入
	postController := &controllers.PostController{}

	// /v1/api に関連するエンドポイントをグループ化
	api := e.Group("/v1/api")
	api.POST("/email", postController.SendEmail)
}