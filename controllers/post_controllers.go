/*
* コントローラ層
 */

package controllers

import (
	"net/http"
	"sync"

	"email-api/models"

	"github.com/labstack/echo/v4"
	"github.com/resend/resend-go/v3"
)

var lock = sync.Mutex{}

// PostController は投稿に関連するコントローラ
type PostController struct {
	Repo models.PostRepository
}

// 投稿の作成のコントローラ
func (p *PostController) SendEmail(c echo.Context) error {
	// mutexを使用して排他制御を行う
	lock.Lock()
	defer lock.Unlock() // 関数を抜ける際にmutexを解放する

	// リクエストボディの取得
	post := new(models.Post)
	if err := c.Bind(post); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body.")
	}
	
	apiKey := "re_DuxRwdR2_2datak4nRcxY91Znhmj2KyxL"

    client := resend.NewClient(apiKey)

    params := &resend.SendEmailRequest{
        From:    "onboarding@resend.dev",
        To:      []string{"halot01025@gmail.com"},
		ReplyTo: post.Email,
        Subject: post.Subject,
        Html:    "<p>Name:" + post.Name + "</p><p>Company:" + post.Company + "</p><p>Message:" + post.Message + "</p>",
    }

    sent, err := client.Emails.Send(params)
    if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{
        "message": "Failed to send email."})
	}

	return c.JSON(http.StatusOK, sent)
}
