/*
* PostRepository インターフェース
 */

package models

// PostRepository インターフェース
type PostRepository interface {
	Create(post *Post) (*Post, error)         // 投稿を作成
}