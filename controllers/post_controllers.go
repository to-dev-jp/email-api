/*
* コントローラ層
 */

package controllers

import (
	"bytes"
	"html/template"
	"net/http"

	"email-api/models"

	"github.com/labstack/echo/v4"
	"github.com/resend/resend-go/v3"
)

// PostController は投稿に関連するコントローラ
type PostController struct {
	Client  *resend.Client
	MyEmail string
}

const notifyTmplText = `
		<div style="font-family: sans-serif; max-width: 600px; margin: 0 auto;">
			<h2 style="color: #333; border-bottom: 2px solid #333; padding-bottom: 8px;">
				新しいお問い合わせが届きました
			</h2>

			<table style="width: 100%; border-collapse: collapse; margin-top: 16px;">
				<tr>
					<td style="padding: 10px; background: #f5f5f5; width: 30%; font-weight: bold;">お名前</td>
					<td style="padding: 10px; border-bottom: 1px solid #eee;">{{.Name}}</td>
				</tr>
				<tr>
					<td style="padding: 10px; background: #f5f5f5; font-weight: bold;">会社名</td>
					<td style="padding: 10px; border-bottom: 1px solid #eee;">{{.Company}}</td>
				</tr>
				<tr>
					<td style="padding: 10px; background: #f5f5f5; font-weight: bold;">メールアドレス</td>
					<td style="padding: 10px; border-bottom: 1px solid #eee;">{{.Email}}</td>
				</tr>
				<tr>
					<td style="padding: 10px; background: #f5f5f5; font-weight: bold;">件名</td>
					<td style="padding: 10px; border-bottom: 1px solid #eee;">{{.Subject}}</td>
				</tr>
				<tr>
					<td style="padding: 10px; background: #f5f5f5; font-weight: bold; vertical-align: top;">メッセージ</td>
					<td style="padding: 10px; border-bottom: 1px solid #eee; white-space: pre-wrap;">{{.Message}}</td>
				</tr>
			</table>

			<p style="margin-top: 24px; color: #666; font-size: 14px;">
				このメールに返信すると、{{.Name}}様へ直接返信されます。
			</p>
		</div>
		`

const autoReplyTmplText = `
			<div style="font-family: sans-serif; max-width: 600px; margin: 0 auto; color: #333;">
				<h2 style="border-bottom: 2px solid #333; padding-bottom: 8px;">
					お問い合わせありがとうございます
				</h2>

				<p>{{.Name}}様</p>

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
					<p style="margin: 4px 0;"><strong>件名：</strong>{{.Subject}}</p>
					<p style="margin: 4px 0;"><strong>メッセージ：</strong></p>
					<p style="margin: 4px 0; white-space: pre-wrap;">{{.Message}}</p>
				</div>

				<p style="color: #666; font-size: 13px;">
					※このメールは自動送信されています。<br>
					※このメールへの返信はお受けしておりません。
				</p>

				<hr style="border: none; border-top: 1px solid #eee; margin: 24px 0;">
				<p style="font-size: 13px; color: #999;">
					岡本 匠<br>{{.MyEmail}}
				</p>
			</div>
		`

var notifyTmpl = template.Must(template.New("notify").Parse(notifyTmplText))

var autoReplyTmpl = template.Must(template.New("autoReply").Parse(autoReplyTmplText))

// 投稿の作成のコントローラ
func (p *PostController) SendEmail(c echo.Context) error {

	// リクエストボディの取得
	post := new(models.Post)
	if err := c.Bind(post); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "リクエストの形式が不正です。")
	}

	// バリデーション追加
	if err := c.Validate(post); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "入力内容に誤りがあります。")
	}

	var notifyBuf bytes.Buffer
	if err := notifyTmpl.Execute(&notifyBuf, post); err != nil {
		c.Logger().Error("通知テンプレート処理失敗:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "サーバーでエラーが発生しました。")
	}

	notifyParams := &resend.SendEmailRequest{
		From:    p.MyEmail,
		To:      []string{p.MyEmail}, // 自分のアドレス
		ReplyTo: post.Email,          // 返信先 = お問い合わせ者
		Subject: "【お問い合わせ】" + post.Subject,
		Html:    notifyBuf.String(),
	}

	data := struct {
		*models.Post
		MyEmail string
	}{Post: post, MyEmail: p.MyEmail}

	var replyBuf bytes.Buffer
	if err := autoReplyTmpl.Execute(&replyBuf, data); err != nil {
		c.Logger().Error("自動返信テンプレート処理失敗:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "サーバーでエラーが発生しました。")
	}
	// 送信者への自動返信
	autoReplyParams := &resend.SendEmailRequest{
		From:    p.MyEmail,
		To:      []string{post.Email},
		Subject: "【自動返信】お問い合わせありがとうございます",
		Html:    replyBuf.String(),
	}

	_, err := p.Client.Emails.Send(notifyParams)
	if err != nil {
		c.Logger().Error("通知メール送信エラー:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "メールの送信に失敗しました。")
	}

	// 自動返信メール送信
	_, err = p.Client.Emails.Send(autoReplyParams)
	if err != nil {
		// 自動返信の失敗はログに残すだけ
		c.Logger().Error("自動返信メール送信エラー:", err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "送信完了"})

}
