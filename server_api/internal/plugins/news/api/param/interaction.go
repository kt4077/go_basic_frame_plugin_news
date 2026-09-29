package param

// ArticleUIDReq 文章业务标识路径参数。
type ArticleUIDReq struct {
	ArticleUID string `form:"article_uid" json:"article_uid" binding:"required,max=32" validate:"文章标识" comment:"文章业务唯一标识"`
}

// CommentUIDReq 评论业务标识路径参数。
type CommentUIDReq struct {
	CommentUID string `form:"comment_uid" json:"comment_uid" binding:"required,max=32" validate:"评论标识" comment:"评论业务唯一标识"`
}

// CommentListReq 评论分页请求。
type CommentListReq struct {
	Page     int `form:"page" validate:"页码" comment:"页码"`
	PageSize int `form:"page_size" binding:"omitempty,max=50" validate:"每页数量" comment:"每页数量，最大50"`
}

// CommentCreateReq 发布评论请求。
type CommentCreateReq struct {
	Content    string   `json:"content" binding:"omitempty,max=1000" validate:"评论内容" comment:"评论内容"`
	Images     []string `json:"images" binding:"omitempty,max=3,dive,max=1024" validate:"评论图片" comment:"评论图片相对路径，最多3张"`
	RequestUID string   `json:"request_uid" binding:"required,max=64" validate:"请求标识" comment:"客户端幂等请求标识"`
}

// ArticleCommentCreateReq 发布一级评论请求。
type ArticleCommentCreateReq struct {
	ArticleUID string `json:"article_uid" binding:"required,max=32" validate:"文章标识" comment:"文章业务唯一标识"`
	CommentCreateReq
}

// CommentReplyCreateReq 发布评论回复请求。
type CommentReplyCreateReq struct {
	CommentUID string `json:"comment_uid" binding:"required,max=32" validate:"评论标识" comment:"被回复评论业务唯一标识"`
	CommentCreateReq
}

// ShareReq 记录分享发起请求。
type ShareReq struct {
	ArticleUID   string `json:"article_uid" binding:"required,max=32" validate:"文章标识" comment:"文章业务唯一标识"`
	RequestUID   string `json:"request_uid" binding:"required,max=64" validate:"请求标识" comment:"客户端幂等请求标识"`
	ShareChannel int    `json:"share_channel" binding:"required,oneof=1 2 3 4" validate:"分享渠道" comment:"分享渠道，1微信好友，2朋友圈，3系统分享，4复制链接"`
}
