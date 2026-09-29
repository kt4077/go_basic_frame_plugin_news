package model

import basemodel "server_api/internal/common/model"

// AdvertisementConfig 配置核心广告在新闻列表或详情页中的展示位置。
type AdvertisementConfig struct {
	basemodel.Base
	AdvertisementID uint `gorm:"not null;uniqueIndex:uk_news_ad_position;index;comment:核心广告配置ID" json:"advertisement_id"`
	Position        int  `gorm:"not null;uniqueIndex:uk_news_ad_position;index;comment:展示位置：1新闻列表，2新闻详情" json:"position"`
	Status          int  `gorm:"not null;default:1;index;comment:启用状态：1启用，2停用" json:"status"`
}

func (AdvertisementConfig) TableName() string { return "plg_news_advertisement_config" }
