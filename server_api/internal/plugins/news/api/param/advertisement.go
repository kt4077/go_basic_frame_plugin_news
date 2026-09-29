package param

type AdvertisementListReq struct {
	Position int `form:"position" binding:"required,oneof=1 2" validate:"展示位置"`
}
