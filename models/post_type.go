/*
* データ構造体を定義する
 */

package models

// 投稿データ構造体を定義する
// 「タグ」機能を用いることで、構造体のフィールドとJSONデータの間で変換を行う
type (
	Post struct {
		Name    string `json:"name" validate:"required,max=50"`
		Company string `json:"company" validate:"max=100"`
		Email   string `json:"email"   validate:"required,email"`
		Subject string `json:"subject" validate:"required,max=100"`
		Message string `json:"message" validate:"required,max=1000"`
	}
)
