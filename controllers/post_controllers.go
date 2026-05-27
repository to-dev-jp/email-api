/*
* コントローラ層
 */

package controllers

import (
	"fmt"
	"net/http"
	"sync"

	"email-api/models"
	"os"

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

	// バリデーション追加
	if err := c.Validate(post); err != nil {
    	return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	

	apiKey := os.Getenv("API_KEY")
	myEmail := os.Getenv("EMAIL_ADDRESS")

    client := resend.NewClient(apiKey)

	notifyParams := &resend.SendEmailRequest{
	From:    myEmail,
	To:      []string{"your@email.com"}, // 自分のアドレス
	ReplyTo: post.Email,                 // 返信先 = お問い合わせ者
	Subject: "【お問い合わせ】" + post.Subject,
	Html: `
		<div style="font-family: sans-serif; max-width: 600px; margin: 0 auto;">
			<h2 style="color: #333; border-bottom: 2px solid #333; padding-bottom: 8px;">
				新しいお問い合わせが届きました
			</h2>

			<table style="width: 100%; border-collapse: collapse; margin-top: 16px;">
				<tr>
					<td style="padding: 10px; background: #f5f5f5; width: 30%; font-weight: bold;">お名前</td>
					<td style="padding: 10px; border-bottom: 1px solid #eee;">` + post.Name + `</td>
				</tr>
				<tr>
					<td style="padding: 10px; background: #f5f5f5; font-weight: bold;">会社名</td>
					<td style="padding: 10px; border-bottom: 1px solid #eee;">` + post.Company + `</td>
				</tr>
				<tr>
					<td style="padding: 10px; background: #f5f5f5; font-weight: bold;">メールアドレス</td>
					<td style="padding: 10px; border-bottom: 1px solid #eee;">` + post.Email + `</td>
				</tr>
				<tr>
					<td style="padding: 10px; background: #f5f5f5; font-weight: bold;">件名</td>
					<td style="padding: 10px; border-bottom: 1px solid #eee;">` + post.Subject + `</td>
				</tr>
				<tr>
					<td style="padding: 10px; background: #f5f5f5; font-weight: bold; vertical-align: top;">メッセージ</td>
					<td style="padding: 10px; border-bottom: 1px solid #eee; white-space: pre-wrap;">` + post.Message + `</td>
				</tr>
			</table>

			<p style="margin-top: 24px; color: #666; font-size: 14px;">
				このメールに返信すると、` + post.Name + ` 様へ直接返信されます。
			</p>
		</div>
		`,
	}

	// 送信者への自動返信
	autoReplyParams := &resend.SendEmailRequest{
		From:    myEmail,
		To:      []string{post.Email},
		Subject: "【自動返信】お問い合わせありがとうございます",
		Html: `
			<div style="font-family: sans-serif; max-width: 600px; margin: 0 auto; color: #333;">
				<h2 style="border-bottom: 2px solid #333; padding-bottom: 8px;">
					お問い合わせありがとうございます
				</h2>

				<p>` + post.Name + ` 様</p>

				<p>
					この度はポートフォリオをご覧いただき、<br>
					お問い合わせいただきありがとうございます。
				</p>
				<p>
					内容を確認の上、<strong>2〜3営業日以内</strong>にご返信いたします。<br>
					しばらくお待ちいただけますと幸いです。
				</p>

				<div style="background: #f9f9f9; border-left: 4px solid #333; padding: 16px; margin: 24px 0;">
					<p style="margin: 0 0 8px; font-weight: bold; color: #555;">【お問い合わせ内容】</p>
					<p style="margin: 4px 0;"><strong>件名：</strong>` + post.Subject + `</p>
					<p style="margin: 4px 0;"><strong>メッセージ：</strong></p>
					<p style="margin: 4px 0; white-space: pre-wrap;">` + post.Message + `</p>
				</div>

				<p style="color: #666; font-size: 13px;">
					※このメールは自動送信されています。<br>
					※このメールへの返信はお受けしておりません。
				</p>

				<hr style="border: none; border-top: 1px solid #eee; margin: 24px 0;">
				<p style="font-size: 13px; color: #999;">
					岡本 匠<br>` + myEmail + `
				</p>
			</div>
		`,
	}

	_, err := client.Emails.Send(notifyParams)
	if err != nil {
    fmt.Println("通知メール送信エラー:", err) // ← ここのメッセージs
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "送信失敗"})
	}

	// 自動返信メール送信
	_, err = client.Emails.Send(autoReplyParams)
	if err != nil {
		// 自動返信の失敗はログに残すだけでもOK
		fmt.Println("自動返信の送信に失敗:", err)
	}

return c.JSON(http.StatusOK, map[string]string{"message": "送信完了"})

}
