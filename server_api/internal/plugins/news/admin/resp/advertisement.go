package resp

type AdvertisementItem struct {
	ID              uint   `json:"id"`
	AdvertisementID uint   `json:"advertisement_id"`
	Name            string `json:"name"`
	AdID            string `json:"ad_id"`
	Format          int    `json:"format"`
	Platforms       []int  `json:"platforms"`
	Status          int    `json:"status"`
	Description     string `json:"description"`
	Position        int    `json:"position"`
	ConfigStatus    int    `json:"config_status"`
}

type AdvertisementOption struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	AdID        string `json:"ad_id"`
	Format      int    `json:"format"`
	Platforms   string `json:"-"`
	PlatformIDs []int  `json:"platforms" gorm:"-"`
	Status      int    `json:"status"`
	Description string `json:"description"`
}
