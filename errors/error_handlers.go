/*
* カスタムエラーハンドラーを実装
 */

package errors

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	var code int
	var message string

	// errが*echo.HTTPError型かどうかをチェック
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if m, ok := he.Message.(string); ok {
			message = m
		} else {
			message = http.StatusText(he.Code)
		}
	} else {
		code = http.StatusInternalServerError
		message = "Internal Server Error"
	}

	// クライアントに返す前にエラーログを記録する
	c.Logger().Errorf("Error: %v, Status code: %d, Request: %s %s", err, code, c.Request().Method, c.Request().URL)

	// レスポンスがクライアントに送信されていない場合のみ、エラーレスポンスを送信
	if !c.Response().Committed {
		c.JSON(code, map[string]interface{}{
			"error":       message,
			"description": http.StatusText(code),
		})
	}
}
