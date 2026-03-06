package requests

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}
