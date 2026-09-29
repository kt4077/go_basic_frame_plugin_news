package resp

import "time"

// CommentItem 管理端评论响应项，不返回会员SN与原始IP。
type CommentItem struct {
	UID           string    `json:"uid" gorm:"column:uid" comment:"评论业务标识"`
	ArticleUID    string    `json:"article_uid" gorm:"column:article_uid" comment:"文章业务标识"`
	ParentUID     string    `json:"parent_uid" gorm:"column:parent_uid" comment:"直接父评论业务标识"`
	RootUID       string    `json:"root_uid" gorm:"column:root_uid" comment:"一级评论业务标识"`
	Content       string    `json:"content" gorm:"column:content" comment:"评论内容"`
	Status        int       `json:"status" gorm:"column:status" comment:"评论状态"`
	LikeCount     uint64    `json:"like_count" gorm:"column:like_count" comment:"点赞数量"`
	ReplyCount    uint64    `json:"reply_count" gorm:"column:reply_count" comment:"回复数量"`
	Platform      int       `json:"platform" gorm:"column:platform" comment:"发布平台"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at" comment:"创建时间"`
	Nickname      string    `json:"nickname" gorm:"-" comment:"评论用户昵称"`
	Account       string    `json:"account" gorm:"-" comment:"脱敏账号"`
	AvatarURL     string    `json:"avatar_url" gorm:"-" comment:"用户头像完整地址"`
	ReplyNickname string    `json:"reply_nickname" gorm:"-" comment:"被回复用户昵称"`
	IPLocation    string    `json:"ip_location" gorm:"-" comment:"脱敏IP位置"`
	Images        []string  `json:"images" gorm:"-" comment:"评论图片完整地址"`
	MemberSN      string    `json:"-" gorm:"column:member_sn" comment:"内部会员业务编号"`
	ReplyMemberSN string    `json:"-" gorm:"column:reply_member_sn" comment:"内部被回复会员业务编号"`
	IP            string    `json:"-" gorm:"column:ip" comment:"内部原始IP"`
}

// CommentListRes 管理端评论分页响应。
type CommentListRes struct {
	List  []CommentItem `json:"list" comment:"评论列表"`
	Total int64         `json:"total" comment:"评论总数"`
}
