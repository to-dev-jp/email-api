/*
* ルーティングの設定を行う
 */

package routers

import (
	"email-api/controllers"

	"github.com/labstack/echo/v4"
	"github.com/resend/resend-go/v3"
)

// ルーティングの設定
func SetupRouter(e *echo.Echo, client *resend.Client, myEmail string) {
	// コントローラを初期化
	postController := &controllers.PostController{
        Client:  client,
        MyEmail: myEmail,
    }

	// /v1/api に関連するエンドポイントをグループ化
	api := e.Group("/v1/api")
	api.POST("/email", postController.SendEmail)
}