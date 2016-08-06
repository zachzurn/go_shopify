package shopify

type Option struct {
	Id        int64  `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Position  int64  `json:"position,omitempty"`
	ProductId int64  `json:"product_id,omitempty"`
}
