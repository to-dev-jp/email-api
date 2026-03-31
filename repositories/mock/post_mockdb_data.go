/*
* モックDBデータ
 */

package mock

import (
	"email-api/models"
)

// サンプルデータ
// 本来はデータベースから取得するが、ここでは簡易的にサンプルデータを使用
var Posts = []models.Post{
	{Name: "1", Subject: "投稿1", Message: "サンプル投稿1"},
	{Name: "2", Subject: "投稿2", Message: "サンプル投稿2"},
	{Name: "3", Subject: "投稿3", Message: "サンプル投稿3"},
	{Name: "4", Subject: "投稿4", Message: "サンプル投稿4"},
	{Name: "5", Subject: "投稿5", Message: "サンプル投稿5"},
	{Name: "6", Subject: "投稿6", Message: "サンプル投稿6"},
	{Name: "7", Subject: "投稿7", Message: "サンプル投稿7"},
	{Name: "8", Subject: "投稿8", Message: "サンプル投稿8"},
	{Name: "9", Subject: "投稿9", Message: "サンプル投稿9"},
	{Name: "10", Subject: "投稿10", Message: "サンプル投稿10"},
}