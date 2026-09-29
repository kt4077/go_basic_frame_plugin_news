package param

type AdvertisementSaveReq struct {
	ID              uint `json:"id"`
	AdvertisementID uint `json:"advertisement_id" binding:"required,gt=0" validate:"广告"`
	Position        int  `json:"position" binding:"required,oneof=1 2" validate:"展示位置"`
	Status          int  `json:"status" binding:"required,oneof=1 2" validate:"启用状态"`
}

type IDReq struct {
	ID uint `json:"id" binding:"required,gt=0" validate:"主键ID"`
}
