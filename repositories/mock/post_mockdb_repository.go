/*
* モックDBのリポジトリ実装
 */

package mock

import (
	"email-api/models"
)

// MockdbPostRepository はPostRepositoryのMock実装
type MockdbPostRepository struct {
	Posts []models.Post
}

// Create は投稿を作成
func (r *MockdbPostRepository) Create(post *models.Post) (*models.Post, error) {
	var newPost models.Post
	newPost.Subject = post.Subject
	newPost.Message = post.Message
	return &newPost, nil
}
