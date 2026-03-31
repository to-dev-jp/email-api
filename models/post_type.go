/*
* データ構造体を定義する
 */

package models

// 投稿データ構造体を定義する
// 「タグ」機能を用いることで、構造体のフィールドとJSONデータの間で変換を行う
type (
	Post struct {
		Name	string `json:"name"`
		Company string `json:"company"`
		Email	string `json:"email"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}
)