package logic

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"server_api/internal/common/app"
	"server_api/internal/common/auth"
	commonEnums "server_api/internal/common/enums"
	commonmiddleware "server_api/internal/common/middleware"
	commonModel "server_api/internal/common/model"
	commonupload "server_api/internal/common/upload"
	"server_api/internal/plugins/news/api/param"
	"server_api/internal/plugins/news/api/resp"
	newsEnums "server_api/internal/plugins/news/enums"
	"server_api/internal/plugins/news/model"
	"server_api/pkg/dberror"
	"server_api/pkg/pagination"
	"server_api/pkg/sn"
)

const (
	commentLimitPerMinute = 5
	likeLimitPerMinute    = 30
	shareLimitPerMinute   = 10
	collectLimitPerMinute = 30
)

func requireMemberSN(c *gin.Context) (string, error) {
	claims := auth.CtxClaims(c)
	if claims == nil || claims.Client != commonEnums.ClientApi || strings.TrimSpace(claims.SN) == "" {
		return "", errors.New("登录状态已过期，请重新登录")
	}
	return claims.SN, nil
}

func limitMemberAction(application *app.App, c *gin.Context, action, memberSN string, limit int) error {
	key := fmt.Sprintf("plugin:news:limit:%s:%s:%s", action, memberSN, time.Now().Format("200601021504"))
	ctx := c.Request.Context()
	count, err := application.Redis.Incr(ctx, key).Result()
	if err != nil {
		return errors.New("服务繁忙，请稍后重试")
	}
	if count == 1 {
		application.Redis.Expire(ctx, key, 2*time.Minute)
	}
	if count > int64(limit) {
		return errors.New("操作过于频繁，请稍后重试")
	}
	return nil
}

func publishedArticle(tx *gorm.DB, articleIdentifier string) (*model.Article, error) {
	var article model.Article
	err := tx.Where("(uid = ? OR slug = ?) AND status = ? AND enable_status = ?", articleIdentifier, articleIdentifier, newsEnums.ArticleStatusPublished, newsEnums.ArticleEnableStatusEnabled).First(&article).Error
	if err != nil {
		return nil, errors.New("文章不存在或已下线")
	}
	if strings.TrimSpace(article.UID) == "" {
		uid, generateErr := generateUID()
		if generateErr != nil {
			return nil, generateErr
		}
		if updateErr := tx.Model(&model.Article{}).Where("id = ? AND (uid IS NULL OR uid = '')", article.ID).Update("uid", uid).Error; updateErr != nil {
			return nil, errors.New("补齐文章业务标识失败")
		}
		article.UID = uid
	}
	return &article, nil
}

func generateUID() (string, error) {
	uid, err := sn.Generate()
	if err != nil {
		return "", errors.New("业务标识生成失败")
	}
	return uid, nil
}

// Comments 查询已展示的一级评论。
func (l *NewsLogic) Comments(c *gin.Context, articleUID string, req *param.CommentListReq) (*resp.CommentListRes, error) {
	if _, err := publishedArticle(l.App.DB.WithContext(c.Request.Context()), articleUID); err != nil {
		return nil, err
	}
	result, err := l.commentPage(c, "article_uid = ? AND parent_uid = ''", articleUID, req)
	if err != nil {
		return nil, err
	}
	if err := l.completeVisibleReplyCounts(c, result.List); err != nil {
		return nil, err
	}
	return result, nil
}

// completeVisibleReplyCounts 批量统计当前页一级评论实际可展示的回复数，避免历史缓存计数失真。
func (l *NewsLogic) completeVisibleReplyCounts(c *gin.Context, list []resp.CommentItem) error {
	if len(list) == 0 {
		return nil
	}
	uids := make([]string, 0, len(list))
	for _, item := range list {
		uids = append(uids, item.UID)
	}
	type replyCount struct {
		ParentUID string `gorm:"column:parent_uid"`
		Total     uint64 `gorm:"column:total"`
	}
	var counts []replyCount
	if err := l.App.DB.WithContext(c.Request.Context()).Model(&model.Comment{}).
		Select("parent_uid, COUNT(*) AS total").
		Where("parent_uid IN ? AND status = ?", uids, newsEnums.CommentStatusVisible).
		Group("parent_uid").Find(&counts).Error; err != nil {
		return errors.New("查询评论回复数量失败")
	}
	countMap := make(map[string]uint64, len(counts))
	for _, count := range counts {
		countMap[count.ParentUID] = count.Total
	}
	for index := range list {
		list[index].ReplyCount = countMap[list[index].UID]
	}
	return nil
}

// Replies 分页查询某条评论的直接已展示回复。
func (l *NewsLogic) Replies(c *gin.Context, commentUID string, req *param.CommentListReq) (*resp.CommentListRes, error) {
	var parent model.Comment
	if err := l.App.DB.WithContext(c.Request.Context()).Where("uid = ? AND status = ?", commentUID, newsEnums.CommentStatusVisible).First(&parent).Error; err != nil {
		return nil, errors.New("评论不存在或已关闭")
	}
	if _, err := publishedArticle(l.App.DB.WithContext(c.Request.Context()), parent.ArticleUID); err != nil {
		return nil, err
	}
	result, err := l.commentPage(c, "parent_uid = ?", parent.UID, req)
	if err != nil {
		return nil, err
	}
	if err := l.completeVisibleReplyCounts(c, result.List); err != nil {
		return nil, err
	}
	return result, nil
}

func (l *NewsLogic) commentPage(c *gin.Context, condition string, value string, req *param.CommentListReq) (*resp.CommentListRes, error) {
	var total int64
	var list []resp.CommentItem
	db := l.App.DB.WithContext(c.Request.Context()).Model(&model.Comment{}).Where(condition, value).Where("status = ?", newsEnums.CommentStatusVisible)
	if err := db.Count(&total).Error; err != nil {
		return nil, errors.New("查询评论失败")
	}
	page, size := pagination.Normalize(req.Page, req.PageSize)
	if size > 50 {
		size = 50
	}
	if err := db.Select("uid,member_sn,reply_member_sn,content,like_count,reply_count,created_at").Order("created_at ASC, id ASC").Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, errors.New("查询评论失败")
	}
	if err := l.completeCommentProfiles(c, list); err != nil {
		return nil, err
	}
	if err := l.completeCommentImages(c, list); err != nil {
		return nil, err
	}
	return &resp.CommentListRes{List: list, Total: total}, nil
}

func (l *NewsLogic) completeCommentProfiles(c *gin.Context, list []resp.CommentItem) error {
	sns := make([]string, 0, len(list)*2)
	for _, item := range list {
		sns = append(sns, item.MemberSN)
		if item.ReplyMemberSN != "" {
			sns = append(sns, item.ReplyMemberSN)
		}
	}
	if len(sns) == 0 {
		return nil
	}
	var members []commonModel.SysMember
	if err := l.App.DB.WithContext(c.Request.Context()).Select("sn,nickname,avatar").Where("sn IN ?", sns).Find(&members).Error; err != nil {
		return errors.New("查询评论用户失败")
	}
	memberMap := make(map[string]commonModel.SysMember, len(members))
	for _, member := range members {
		memberMap[member.SN] = member
	}
	for index := range list {
		member, ok := memberMap[list[index].MemberSN]
		if !ok {
			list[index].Nickname = "已注销用户"
			continue
		}
		list[index].Nickname = member.Nickname
		avatarURL, err := commonupload.FileURL(l.App, member.Avatar)
		if err != nil {
			return errors.New("头像地址解析失败")
		}
		list[index].AvatarURL = avatarURL
		if replyMember, exists := memberMap[list[index].ReplyMemberSN]; exists {
			list[index].ReplyNickname = replyMember.Nickname
		}
	}
	return nil
}

func (l *NewsLogic) completeCommentImages(c *gin.Context, list []resp.CommentItem) error {
	if len(list) == 0 {
		return nil
	}
	uids := make([]string, 0, len(list))
	for _, item := range list {
		uids = append(uids, item.UID)
	}
	var images []model.CommentImage
	if err := l.App.DB.WithContext(c.Request.Context()).Where("comment_uid IN ?", uids).Order("sort ASC,id ASC").Find(&images).Error; err != nil {
		return errors.New("查询评论图片失败")
	}
	imageMap := make(map[string][]string)
	for _, image := range images {
		url, err := commonupload.FileURL(l.App, image.Path)
		if err != nil {
			return errors.New("评论图片地址解析失败")
		}
		imageMap[image.CommentUID] = append(imageMap[image.CommentUID], url)
	}
	for index := range list {
		list[index].Images = imageMap[list[index].UID]
	}
	return nil
}

// CreateComment 发布一级评论或回复，身份只从鉴权上下文获取。
func (l *NewsLogic) CreateComment(c *gin.Context, articleUID, parentUID string, req *param.CommentCreateReq) (*resp.CommentCreateRes, error) {
	memberSN, err := requireMemberSN(c)
	if err != nil {
		return nil, err
	}
	if err := limitMemberAction(l.App, c, "comment", memberSN, commentLimitPerMinute); err != nil {
		return nil, err
	}
	content := strings.TrimSpace(req.Content)
	if content == "" && len(req.Images) == 0 {
		return nil, errors.New("评论内容和图片不能同时为空")
	}
	if utf8.RuneCountInString(content) > 1000 {
		return nil, errors.New("评论内容不能超过1000个字符")
	}
	uid, err := generateUID()
	if err != nil {
		return nil, err
	}
	comment := model.Comment{UID: uid, ArticleUID: articleUID, MemberSN: memberSN, RequestUID: req.RequestUID, ParentUID: parentUID, RootUID: uid, Content: content, Status: newsEnums.CommentStatusPending, Platform: commonmiddleware.CtxPlatformSource(c), IP: c.ClientIP()}
	err = l.App.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if _, err := publishedArticle(tx, articleUID); err != nil {
			return err
		}
		if parentUID != "" {
			var parent model.Comment
			if err := tx.Where("uid = ? AND article_uid = ? AND status = ?", parentUID, articleUID, newsEnums.CommentStatusVisible).First(&parent).Error; err != nil {
				return errors.New("被回复评论不存在或已关闭")
			}
			comment.RootUID = parent.RootUID
			comment.ReplyMemberSN = parent.MemberSN
		}
		if err := tx.Create(&comment).Error; err != nil {
			return err
		}
		resolver, err := commonupload.NewURLResolver(l.App)
		if err != nil && len(req.Images) > 0 {
			return errors.New("存储配置不可用")
		}
		for index, source := range req.Images {
			path, err := resolver.Relative(source)
			if err != nil {
				return err
			}
			imageUID, err := generateUID()
			if err != nil {
				return err
			}
			if err := tx.Create(&model.CommentImage{UID: imageUID, CommentUID: comment.UID, Path: path, Sort: index}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		if dberror.IsDuplicateKey(err) {
			return nil, errors.New("评论提交重复，请勿重复操作")
		}
		return nil, err
	}
	return &resp.CommentCreateRes{UID: comment.UID, Status: comment.Status}, nil
}

// ReplyComment 根据父评论业务标识发布回复，文章关系由服务端解析。
func (l *NewsLogic) ReplyComment(c *gin.Context, parentUID string, req *param.CommentCreateReq) (*resp.CommentCreateRes, error) {
	var articleUID string
	if err := l.App.DB.WithContext(c.Request.Context()).Model(&model.Comment{}).
		Where("uid = ?", parentUID).Pluck("article_uid", &articleUID).Error; err != nil || articleUID == "" {
		return nil, errors.New("评论不存在")
	}
	return l.CreateComment(c, articleUID, parentUID, req)
}

// ToggleArticleLike 原子切换文章点赞状态并返回最新状态与计数。
func (l *NewsLogic) ToggleArticleLike(c *gin.Context, articleIdentifier string) (*resp.InteractionToggleRes, error) {
	memberSN, err := requireMemberSN(c)
	if err != nil {
		return nil, err
	}
	if err := limitMemberAction(l.App, c, "like", memberSN, likeLimitPerMinute); err != nil {
		return nil, err
	}
	result := &resp.InteractionToggleRes{}
	err = l.App.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		article, err := publishedArticle(tx, articleIdentifier)
		if err != nil {
			return err
		}
		var record model.ArticleLike
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("article_uid = ? AND member_sn = ?", article.UID, memberSN).First(&record).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			uid, err := generateUID()
			if err != nil {
				return err
			}
			if err := tx.Create(&model.ArticleLike{UID: uid, ArticleUID: article.UID, MemberSN: memberSN, Status: newsEnums.InteractionStatusActive}).Error; err != nil {
				return err
			}
			result.Active = true
		} else if err != nil {
			return err
		} else {
			result.Active = record.Status != newsEnums.InteractionStatusActive
			target := newsEnums.InteractionStatusCanceled
			if result.Active {
				target = newsEnums.InteractionStatusActive
			}
			if err := tx.Model(&record).Update("status", target).Error; err != nil {
				return err
			}
		}
		expression := "GREATEST(like_count - 1, 0)"
		if result.Active {
			expression = "like_count + 1"
		}
		if err := tx.Model(&model.Article{}).Where("uid = ?", article.UID).UpdateColumn("like_count", gorm.Expr(expression)).Error; err != nil {
			return err
		}
		return tx.Model(&model.Article{}).Where("uid = ?", article.UID).Pluck("like_count", &result.Count).Error
	})
	return result, err
}

// ToggleArticleCollection 原子切换文章收藏状态并返回最新状态与计数。
func (l *NewsLogic) ToggleArticleCollection(c *gin.Context, articleIdentifier string) (*resp.InteractionToggleRes, error) {
	memberSN, err := requireMemberSN(c)
	if err != nil {
		return nil, err
	}
	if err := limitMemberAction(l.App, c, "collection", memberSN, collectLimitPerMinute); err != nil {
		return nil, err
	}
	result := &resp.InteractionToggleRes{}
	err = l.App.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		article, err := publishedArticle(tx, articleIdentifier)
		if err != nil {
			return err
		}
		var record model.ArticleCollection
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("article_uid = ? AND member_sn = ?", article.UID, memberSN).First(&record).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			uid, err := generateUID()
			if err != nil {
				return err
			}
			if err := tx.Create(&model.ArticleCollection{UID: uid, ArticleUID: article.UID, MemberSN: memberSN, Status: newsEnums.InteractionStatusActive}).Error; err != nil {
				return err
			}
			result.Active = true
		} else if err != nil {
			return err
		} else {
			result.Active = record.Status != newsEnums.InteractionStatusActive
			target := newsEnums.InteractionStatusCanceled
			if result.Active {
				target = newsEnums.InteractionStatusActive
			}
			if err := tx.Model(&record).Update("status", target).Error; err != nil {
				return err
			}
		}
		expression := "GREATEST(collection_count - 1, 0)"
		if result.Active {
			expression = "collection_count + 1"
		}
		if err := tx.Model(&model.Article{}).Where("uid = ?", article.UID).UpdateColumn("collection_count", gorm.Expr(expression)).Error; err != nil {
			return err
		}
		return tx.Model(&model.Article{}).Where("uid = ?", article.UID).Pluck("collection_count", &result.Count).Error
	})
	return result, err
}

// ToggleCommentLike 原子切换评论点赞状态并返回最新状态与计数。
func (l *NewsLogic) ToggleCommentLike(c *gin.Context, commentUID string) (*resp.InteractionToggleRes, error) {
	memberSN, err := requireMemberSN(c)
	if err != nil {
		return nil, err
	}
	if err := limitMemberAction(l.App, c, "comment-like", memberSN, likeLimitPerMinute); err != nil {
		return nil, err
	}
	result := &resp.InteractionToggleRes{}
	err = l.App.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var comment model.Comment
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("uid = ? AND status = ?", commentUID, newsEnums.CommentStatusVisible).First(&comment).Error; err != nil {
			return errors.New("评论不存在或已关闭")
		}
		var record model.CommentLike
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("comment_uid = ? AND member_sn = ?", commentUID, memberSN).First(&record).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			uid, err := generateUID()
			if err != nil {
				return err
			}
			if err := tx.Create(&model.CommentLike{UID: uid, CommentUID: commentUID, MemberSN: memberSN, Status: newsEnums.InteractionStatusActive}).Error; err != nil {
				return err
			}
			result.Active = true
		} else if err != nil {
			return err
		} else {
			result.Active = record.Status != newsEnums.InteractionStatusActive
			target := newsEnums.InteractionStatusCanceled
			if result.Active {
				target = newsEnums.InteractionStatusActive
			}
			if err := tx.Model(&record).Update("status", target).Error; err != nil {
				return err
			}
		}
		expression := "GREATEST(like_count - 1, 0)"
		if result.Active {
			expression = "like_count + 1"
		}
		if err := tx.Model(&comment).UpdateColumn("like_count", gorm.Expr(expression)).Error; err != nil {
			return err
		}
		return tx.Model(&model.Comment{}).Where("uid = ?", commentUID).Pluck("like_count", &result.Count).Error
	})
	return result, err
}

// RecordShare 幂等记录用户发起分享。
func (l *NewsLogic) RecordShare(c *gin.Context, articleUID string, req *param.ShareReq) error {
	memberSN, err := requireMemberSN(c)
	if err != nil {
		return err
	}
	if err := limitMemberAction(l.App, c, "share", memberSN, shareLimitPerMinute); err != nil {
		return err
	}
	uid, err := generateUID()
	if err != nil {
		return err
	}
	return l.App.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if _, err := publishedArticle(tx, articleUID); err != nil {
			return err
		}
		record := model.ShareRecord{UID: uid, ArticleUID: articleUID, MemberSN: memberSN, RequestUID: req.RequestUID, ShareChannel: req.ShareChannel, Platform: commonmiddleware.CtxPlatformSource(c)}
		if err := tx.Create(&record).Error; err != nil {
			if dberror.IsDuplicateKey(err) {
				return nil
			}
			return err
		}
		return tx.Model(&model.Article{}).Where("uid = ?", articleUID).UpdateColumn("share_count", gorm.Expr("share_count + 1")).Error
	})
}

// InteractionState 查询当前用户文章点赞状态。
func (l *NewsLogic) InteractionState(c *gin.Context, articleUID string) (*resp.InteractionStateRes, error) {
	memberSN, err := requireMemberSN(c)
	if err != nil {
		return nil, err
	}
	article, err := publishedArticle(l.App.DB.WithContext(c.Request.Context()), articleUID)
	if err != nil {
		return nil, err
	}
	var likeCount int64
	if err = l.App.DB.WithContext(c.Request.Context()).Model(&model.ArticleLike{}).Where("article_uid = ? AND member_sn = ? AND status = ?", article.UID, memberSN, newsEnums.InteractionStatusActive).Count(&likeCount).Error; err != nil {
		return nil, errors.New("查询互动状态失败")
	}
	var collectionCount int64
	if err = l.App.DB.WithContext(c.Request.Context()).Model(&model.ArticleCollection{}).Where("article_uid = ? AND member_sn = ? AND status = ?", article.UID, memberSN, newsEnums.InteractionStatusActive).Count(&collectionCount).Error; err != nil {
		return nil, errors.New("查询互动状态失败")
	}
	return &resp.InteractionStateRes{Liked: likeCount > 0, Collected: collectionCount > 0}, nil
}
