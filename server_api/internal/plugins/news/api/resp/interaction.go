package resp

import "time"

// CommentItem 用户端评论响应，不包含会员业务编号、账号和IP。
type CommentItem struct {
	UID           string    `json:"uid" gorm:"column:uid" comment:"评论业务唯一标识"`
	Nickname      string    `json:"nickname" gorm:"-" comment:"评论用户昵称"`
	AvatarURL     string    `json:"avatar_url" gorm:"-" comment:"评论用户头像完整地址"`
	ReplyNickname string    `json:"reply_nickname,omitempty" gorm:"-" comment:"被回复用户昵称"`
	Content       string    `json:"content" gorm:"column:content" comment:"评论内容"`
	Images        []string  `json:"images" gorm:"-" comment:"评论图片完整地址"`
	LikeCount     uint64    `json:"like_count" gorm:"column:like_count" comment:"点赞数量"`
	ReplyCount    uint64    `json:"reply_count" gorm:"column:reply_count" comment:"回复数量"`
	IsLiked       bool      `json:"is_liked" gorm:"-" comment:"当前用户是否点赞，公开接口固定为否"`
	IsMine        bool      `json:"is_mine" gorm:"-" comment:"是否当前用户评论，公开接口固定为否"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at" comment:"创建时间"`
	MemberSN      string    `json:"-" gorm:"column:member_sn" comment:"会员业务编号，仅内部组装"`
	ReplyMemberSN string    `json:"-" gorm:"column:reply_member_sn" comment:"被回复会员业务编号，仅内部组装"`
}

func (CommentItem) TableName() string { return "plg_news_comment" }

// CommentListRes 评论分页响应。
type CommentListRes struct {
	List  []CommentItem `json:"list" comment:"评论列表"`
	Total int64         `json:"total" comment:"评论总数"`
}

// CommentCreateRes 发布评论响应。
type CommentCreateRes struct {
	UID    string `json:"uid" comment:"评论业务唯一标识"`
	Status int    `json:"status" comment:"审核状态，1待审核"`
}

// InteractionStateRes 当前用户文章互动状态。
type InteractionStateRes struct {
	Liked     bool `json:"liked" comment:"是否已点赞文章"`
	Collected bool `json:"collected" comment:"是否已收藏文章"`
}

// InteractionToggleRes 文章互动切换结果。
type InteractionToggleRes struct {
	Active bool   `json:"active" comment:"切换后的激活状态"`
	Count  uint64 `json:"count" comment:"切换后的互动数量"`
}
