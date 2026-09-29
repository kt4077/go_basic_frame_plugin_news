package param

// CommentListReq 评论分页查询参数。
type CommentListReq struct {
	ArticleUID string `form:"article_uid" binding:"omitempty,max=32" validate:"文章标识" comment:"文章业务标识"`
	RootUID    string `form:"root_uid" binding:"omitempty,max=32" validate:"一级评论标识" comment:"一级评论业务标识"`
	Keyword    string `form:"keyword" binding:"omitempty,max=64" validate:"评论关键词" comment:"评论内容或用户关键词"`
	Status     int    `form:"status" binding:"omitempty,oneof=1 2 3" validate:"评论状态" comment:"评论状态"`
	Page       int    `form:"page" binding:"omitempty,min=1" validate:"页码" comment:"页码"`
	PageSize   int    `form:"page_size" binding:"omitempty,min=1,max=50" validate:"每页数量" comment:"每页数量"`
}

// CommentStatusReq 评论展示状态修改参数。
type CommentStatusReq struct {
	CommentUID string `json:"comment_uid" binding:"required,max=32" validate:"评论标识" comment:"评论业务标识"`
	Status     int    `json:"status" binding:"required,oneof=2 3" validate:"评论状态" comment:"评论状态，2显示，3关闭"`
}
