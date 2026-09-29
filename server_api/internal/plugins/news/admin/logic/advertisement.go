package logic

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"server_api/internal/common/app"
	newsparam "server_api/internal/plugins/news/admin/param"
	newsresp "server_api/internal/plugins/news/admin/resp"
	newsEnums "server_api/internal/plugins/news/enums"
	newsmodel "server_api/internal/plugins/news/model"
	"server_api/pkg/dberror"
)

type AdvertisementLogic struct{ App *app.App }

type advertisementRow struct {
	ID              uint   `gorm:"column:id"`
	AdvertisementID uint   `gorm:"column:advertisement_id"`
	Name            string `gorm:"column:name"`
	AdID            string `gorm:"column:ad_id"`
	Format          int    `gorm:"column:format"`
	Platforms       string `gorm:"column:platforms"`
	Status          int    `gorm:"column:status"`
	Description     string `gorm:"column:description"`
	Position        int    `gorm:"column:position"`
	ConfigStatus    int    `gorm:"column:config_status"`
}

func decodePlatforms(value string) []int {
	result := make([]int, 0)
	for _, part := range strings.Split(value, ",") {
		if item, err := strconv.Atoi(part); err == nil {
			result = append(result, item)
		}
	}
	return result
}

func (l *AdvertisementLogic) List(c *gin.Context) ([]newsresp.AdvertisementItem, error) {
	var rows []advertisementRow
	err := l.App.DB.WithContext(c.Request.Context()).Table("plg_news_advertisement_config c").
		Select("c.id,c.advertisement_id,c.position,c.status AS config_status,a.name,a.ad_id,a.format,a.platforms,a.status,a.description").
		Joins("JOIN sys_advertisement a ON a.id=c.advertisement_id AND a.deleted_at IS NULL").Order("c.position ASC,c.id DESC").Scan(&rows).Error
	if err != nil {
		return nil, errors.New("查询新闻广告配置失败")
	}
	result := make([]newsresp.AdvertisementItem, 0, len(rows))
	for _, row := range rows {
		result = append(result, newsresp.AdvertisementItem{ID: row.ID, AdvertisementID: row.AdvertisementID, Name: row.Name, AdID: row.AdID, Format: row.Format, Platforms: decodePlatforms(row.Platforms), Status: row.Status, Description: row.Description, Position: row.Position, ConfigStatus: row.ConfigStatus})
	}
	return result, nil
}

func (l *AdvertisementLogic) Options(c *gin.Context) ([]newsresp.AdvertisementOption, error) {
	var rows []newsresp.AdvertisementOption
	err := l.App.DB.WithContext(c.Request.Context()).Table("sys_advertisement a").
		Select("a.id,a.name,a.ad_id,a.format,a.platforms,a.status,a.description").
		Joins("JOIN sys_advertisement_plugin ap ON ap.advertisement_id=a.id").
		Joins("JOIN sys_plugin p ON p.id=ap.plugin_id AND p.plugin_id=? AND p.deleted_at IS NULL", "news").
		Where("a.deleted_at IS NULL").Order("a.id DESC").Scan(&rows).Error
	if err != nil {
		return nil, errors.New("查询新闻插件可用广告失败")
	}
	for index := range rows {
		rows[index].PlatformIDs = decodePlatforms(rows[index].Platforms)
	}
	return rows, nil
}

func (l *AdvertisementLogic) Save(c *gin.Context, req *newsparam.AdvertisementSaveReq) error {
	if !newsEnums.IsValidAdvertisementPosition(req.Position) {
		return errors.New("展示位置不正确")
	}
	var count int64
	err := l.App.DB.WithContext(c.Request.Context()).Table("sys_advertisement a").
		Joins("JOIN sys_advertisement_plugin ap ON ap.advertisement_id=a.id").
		Joins("JOIN sys_plugin p ON p.id=ap.plugin_id AND p.plugin_id=? AND p.deleted_at IS NULL", "news").
		Where("a.id=? AND a.deleted_at IS NULL", req.AdvertisementID).Count(&count).Error
	if err != nil || count == 0 {
		return errors.New("所选广告未绑定新闻插件")
	}
	values := newsmodel.AdvertisementConfig{AdvertisementID: req.AdvertisementID, Position: req.Position, Status: req.Status}
	if req.ID == 0 {
		if err := l.App.DB.WithContext(c.Request.Context()).Create(&values).Error; err != nil {
			if dberror.IsDuplicateKey(err) {
				return errors.New("该广告已配置在此展示位置")
			}
			return errors.New("新增新闻广告配置失败")
		}
		return nil
	}
	result := l.App.DB.WithContext(c.Request.Context()).Model(&newsmodel.AdvertisementConfig{}).Where("id=?", req.ID).Updates(map[string]interface{}{"advertisement_id": req.AdvertisementID, "position": req.Position, "status": req.Status})
	if result.Error != nil {
		return errors.New("保存新闻广告配置失败")
	}
	return nil
}

func (l *AdvertisementLogic) Delete(c *gin.Context, id uint) error {
	result := l.App.DB.WithContext(c.Request.Context()).Unscoped().Delete(&newsmodel.AdvertisementConfig{}, id)
	if result.Error != nil || result.RowsAffected == 0 {
		return errors.New("新闻广告配置不存在或删除失败")
	}
	return nil
}
