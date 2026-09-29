package logic

import (
	"errors"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"server_api/internal/common/app"
	commonModel "server_api/internal/common/model"
	commonupload "server_api/internal/common/upload"
	"server_api/internal/plugins/news/admin/param"
	"server_api/internal/plugins/news/admin/resp"
	newsEnums "server_api/internal/plugins/news/enums"
	"server_api/internal/plugins/news/model"
	"server_api/pkg/mask"
	"server_api/pkg/pagination"
)

// CommentLogic 新闻评论管理业务逻辑。
type CommentLogic struct {
	App *app.App
}

// List 分页查询一级评论或指定一级评论下的回复。
func (l *CommentLogic) List(c *gin.Context, req *param.CommentListReq, replies bool) (*resp.CommentListRes, error) {
	var total int64
	var list []resp.CommentItem
	db := l.App.DB.WithContext(c.Request.Context()).Model(&model.Comment{})
	if replies {
		if strings.TrimSpace(req.RootUID) == "" {
			return nil, errors.New("父评论业务标识不能为空")
		}
		db = db.Where("parent_uid = ?", req.RootUID)
	} else {
		db = db.Where("parent_uid = ''")
	}
	if articleUID := strings.TrimSpace(req.ArticleUID); articleUID != "" {
		db = db.Where("article_uid = ?", articleUID)
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}
	if keyword := strings.TrimSpace(req.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		var memberSNs []string
		l.App.DB.WithContext(c.Request.Context()).Model(&commonModel.SysMember{}).
			Where("nickname LIKE ? OR mobile LIKE ? OR account LIKE ?", like, like, like).
			Limit(200).Pluck("sn", &memberSNs)
		if len(memberSNs) == 0 {
			db = db.Where("content LIKE ?", like)
		} else {
			db = db.Where("content LIKE ? OR member_sn IN ?", like, memberSNs)
		}
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, errors.New("查询评论失败")
	}
	page, size := pagination.Normalize(req.Page, req.PageSize)
	if size > 50 {
		size = 50
	}
	if err := db.Select("uid,article_uid,member_sn,parent_uid,root_uid,reply_member_sn,content,status,like_count,reply_count,platform,ip,created_at").
		Order("created_at DESC, id DESC").Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, errors.New("查询评论失败")
	}
	if err := l.completeProfiles(c, list); err != nil {
		return nil, err
	}
	if err := l.completeImages(c, list); err != nil {
		return nil, err
	}
	if err := l.completeReplyCounts(c, list); err != nil {
		return nil, err
	}
	return &resp.CommentListRes{List: list, Total: total}, nil
}

// completeReplyCounts 为管理端补充包含全部状态的实际回复数量。
func (l *CommentLogic) completeReplyCounts(c *gin.Context, list []resp.CommentItem) error {
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
		Where("parent_uid IN ?", uids).
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

func (l *CommentLogic) completeProfiles(c *gin.Context, list []resp.CommentItem) error {
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
	if err := l.App.DB.WithContext(c.Request.Context()).Select("sn,nickname,avatar,mobile,account").Where("sn IN ?", sns).Find(&members).Error; err != nil {
		return errors.New("查询评论用户失败")
	}
	memberMap := make(map[string]commonModel.SysMember, len(members))
	for _, member := range members {
		memberMap[member.SN] = member
	}
	for index := range list {
		member, exists := memberMap[list[index].MemberSN]
		if exists {
			list[index].Nickname = member.Nickname
			list[index].Account = mask.Mobile(member.Mobile)
			avatarURL, err := commonupload.FileURL(l.App, member.Avatar)
			if err != nil {
				return errors.New("头像地址解析失败")
			}
			list[index].AvatarURL = avatarURL
		} else {
			list[index].Nickname = "已注销用户"
		}
		if replyMember, ok := memberMap[list[index].ReplyMemberSN]; ok {
			list[index].ReplyNickname = replyMember.Nickname
		}
		list[index].IPLocation = maskIP(list[index].IP)
	}
	return nil
}

func (l *CommentLogic) completeImages(c *gin.Context, list []resp.CommentItem) error {
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
		imageURL, err := commonupload.FileURL(l.App, image.Path)
		if err != nil {
			return errors.New("评论图片地址解析失败")
		}
		imageMap[image.CommentUID] = append(imageMap[image.CommentUID], imageURL)
	}
	for index := range list {
		list[index].Images = imageMap[list[index].UID]
	}
	return nil
}

// UpdateStatus 修改评论展示状态，并原子维护文章与一级评论计数。
func (l *CommentLogic) UpdateStatus(c *gin.Context, req *param.CommentStatusReq) error {
	return l.App.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var comment model.Comment
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("uid = ?", req.CommentUID).First(&comment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("评论不存在")
			}
			return errors.New("查询评论失败")
		}
		if comment.Status == req.Status {
			return nil
		}
		if err := tx.Model(&comment).Update("status", req.Status).Error; err != nil {
			return errors.New("更新评论状态失败")
		}
		delta := 0
		if comment.Status == newsEnums.CommentStatusVisible && req.Status != newsEnums.CommentStatusVisible {
			delta = -1
		}
		if comment.Status != newsEnums.CommentStatusVisible && req.Status == newsEnums.CommentStatusVisible {
			delta = 1
		}
		if delta == 0 {
			return nil
		}
		if err := updateUnsignedCounter(tx, &model.Article{}, "uid = ?", comment.ArticleUID, "comment_count", delta); err != nil {
			return err
		}
		if comment.ParentUID != "" {
			return updateUnsignedCounter(tx, &model.Comment{}, "uid = ?", comment.ParentUID, "reply_count", delta)
		}
		return nil
	})
}

func updateUnsignedCounter(tx *gorm.DB, modelValue interface{}, condition string, value interface{}, column string, delta int) error {
	expression := column + " + 1"
	if delta < 0 {
		expression = "GREATEST(" + column + " - 1, 0)"
	}
	return tx.Model(modelValue).Where(condition, value).UpdateColumn(column, gorm.Expr(expression)).Error
}

func maskIP(value string) string {
	ip := net.ParseIP(strings.TrimSpace(value))
	if ip == nil {
		return ""
	}
	if ipv4 := ip.To4(); ipv4 != nil {
		return ipv4.String()[:strings.LastIndex(ipv4.String(), ".")+1] + "*"
	}
	parts := strings.Split(ip.String(), ":")
	if len(parts) > 3 {
		return strings.Join(parts[:3], ":") + ":*"
	}
	return "*"
}
