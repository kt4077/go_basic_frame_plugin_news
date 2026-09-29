package resp

type AdvertisementItem struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	AdID        string `json:"ad_id"`
	Format      int    `json:"format"`
	Description string `json:"description"`
}
