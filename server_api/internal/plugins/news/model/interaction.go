package model

import (
	"time"

	"gorm.io/gorm"
)

// ArticleLike 新闻文章点赞记录。表：plg_news_article_like。
type ArticleLike struct {
	ID         uint      `gorm:"primaryKey;comment:主键ID" json:"id"`
	UID        string    `gorm:"size:32;not null;uniqueIndex;comment:点赞业务唯一标识" json:"uid"`
	ArticleUID string    `gorm:"size:32;not null;uniqueIndex:uk_plg_news_article_like_member;index;comment:文章业务唯一标识" json:"article_uid"`
	MemberSN   string    `gorm:"size:32;not null;uniqueIndex:uk_plg_news_article_like_member;index;comment:会员业务编号" json:"-"`
	Status     int       `gorm:"not null;default:1;comment:状态，1点赞，2取消" json:"status"`
	CreatedAt  time.Time `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt  time.Time `gorm:"comment:更新时间" json:"updated_at"`
}

func (ArticleLike) TableName() string { return "plg_news_article_like" }

// ArticleCollection 新闻文章收藏记录。表：plg_news_article_collection。
type ArticleCollection struct {
	ID         uint      `gorm:"primaryKey;comment:主键ID" json:"id"`
	UID        string    `gorm:"size:32;not null;uniqueIndex;comment:收藏业务唯一标识" json:"uid"`
	ArticleUID string    `gorm:"size:32;not null;uniqueIndex:uk_plg_news_article_collection_member;index;comment:文章业务唯一标识" json:"article_uid"`
	MemberSN   string    `gorm:"size:32;not null;uniqueIndex:uk_plg_news_article_collection_member;index;comment:会员业务编号" json:"-"`
	Status     int       `gorm:"not null;default:1;comment:状态，1收藏，2取消" json:"status"`
	CreatedAt  time.Time `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt  time.Time `gorm:"comment:更新时间" json:"updated_at"`
}

func (ArticleCollection) TableName() string { return "plg_news_article_collection" }

// Comment 新闻文章评论。表：plg_news_comment。
type Comment struct {
	ID            uint           `gorm:"primaryKey;comment:主键ID" json:"id"`
	UID           string         `gorm:"size:32;not null;uniqueIndex;comment:评论业务唯一标识" json:"uid"`
	ArticleUID    string         `gorm:"size:32;not null;index;comment:文章业务唯一标识" json:"article_uid"`
	MemberSN      string         `gorm:"size:32;not null;index;comment:评论会员业务编号" json:"-"`
	RequestUID    string         `gorm:"size:64;not null;uniqueIndex:uk_plg_news_comment_request;comment:客户端幂等请求标识" json:"request_uid"`
	ParentUID     string         `gorm:"size:32;not null;default:'';index;comment:直接父评论业务标识" json:"parent_uid"`
	RootUID       string         `gorm:"size:32;not null;index;comment:一级评论业务标识" json:"root_uid"`
	ReplyMemberSN string         `gorm:"size:32;not null;default:'';comment:被回复会员业务编号" json:"-"`
	Content       string         `gorm:"size:1000;not null;default:'';comment:评论内容" json:"content"`
	Status        int            `gorm:"not null;default:1;index;comment:状态，1待审核，2显示，3关闭" json:"status"`
	LikeCount     uint64         `gorm:"not null;default:0;comment:点赞数量" json:"like_count"`
	ReplyCount    uint64         `gorm:"not null;default:0;comment:回复数量" json:"reply_count"`
	Platform      int            `gorm:"not null;default:1;comment:发布平台" json:"platform"`
	IP            string         `gorm:"size:64;not null;default:'';comment:发布IP" json:"-"`
	CreatedAt     time.Time      `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"comment:更新时间" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index;comment:删除时间" json:"-"`
}

func (Comment) TableName() string { return "plg_news_comment" }

// CommentImage 新闻评论图片。表：plg_news_comment_image。
type CommentImage struct {
	ID         uint      `gorm:"primaryKey;comment:主键ID" json:"id"`
	UID        string    `gorm:"size:32;not null;uniqueIndex;comment:评论图片业务唯一标识" json:"uid"`
	CommentUID string    `gorm:"size:32;not null;index;comment:评论业务唯一标识" json:"comment_uid"`
	Path       string    `gorm:"size:1024;not null;comment:图片相对路径" json:"path"`
	Sort       int       `gorm:"not null;default:0;comment:排序值" json:"sort"`
	CreatedAt  time.Time `gorm:"comment:创建时间" json:"created_at"`
}

func (CommentImage) TableName() string { return "plg_news_comment_image" }

// CommentLike 新闻评论点赞记录。表：plg_news_comment_like。
type CommentLike struct {
	ID         uint      `gorm:"primaryKey;comment:主键ID" json:"id"`
	UID        string    `gorm:"size:32;not null;uniqueIndex;comment:评论点赞业务唯一标识" json:"uid"`
	CommentUID string    `gorm:"size:32;not null;uniqueIndex:uk_plg_news_comment_like_member;index;comment:评论业务唯一标识" json:"comment_uid"`
	MemberSN   string    `gorm:"size:32;not null;uniqueIndex:uk_plg_news_comment_like_member;index;comment:会员业务编号" json:"-"`
	Status     int       `gorm:"not null;default:1;comment:状态，1点赞，2取消" json:"status"`
	CreatedAt  time.Time `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt  time.Time `gorm:"comment:更新时间" json:"updated_at"`
}

func (CommentLike) TableName() string { return "plg_news_comment_like" }

// ShareRecord 新闻文章分享发起记录。表：plg_news_share_record。
type ShareRecord struct {
	ID           uint      `gorm:"primaryKey;comment:主键ID" json:"id"`
	UID          string    `gorm:"size:32;not null;uniqueIndex;comment:分享记录业务唯一标识" json:"uid"`
	ArticleUID   string    `gorm:"size:32;not null;index;comment:文章业务唯一标识" json:"article_uid"`
	MemberSN     string    `gorm:"size:32;not null;index;comment:会员业务编号" json:"-"`
	RequestUID   string    `gorm:"size:64;not null;uniqueIndex:uk_plg_news_share_request;comment:客户端请求唯一标识" json:"request_uid"`
	ShareChannel int       `gorm:"not null;default:1;comment:分享渠道" json:"share_channel"`
	Platform     int       `gorm:"not null;default:1;comment:发起平台" json:"platform"`
	CreatedAt    time.Time `gorm:"comment:创建时间" json:"created_at"`
}

func (ShareRecord) TableName() string { return "plg_news_share_record" }
