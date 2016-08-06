package shopify

type Variant struct {
	Barcode              string      `json:"barcode,omitempty"`
	CompareAtPrice       string      `json:"compare_at_price,omitempty"`
	CreatedAt            string      `json:"created_at,omitempty"`
	FulfillmentService   string      `json:"fulfillment_service,omitempty"`
	Grams                float64     `json:"grams,omitempty"`
	Weight               float64     `json:"weight,omitempty"`
	WeightUnit           string      `json:"weight_unit,omitempty"`
	Id                   int64       `json:"id,omitempty"`
	InventoryManagement  string      `json:"inventory_management,omitempty"`
	InventoryPolicy      string      `json:"inventory_policy,omitempty"`
	InventoryQuantity    int64       `json:"inventory_quantity,omitempty"`
	OldInventoryQuantity int64       `json:"old_inventory_quantity,omitempty"`
	Metafield            interface{} `json:"metafield,omitempty"`
	Option1              string      `json:"option1,omitempty"`
	Option2              string      `json:"option2,omitempty"`
	Option3              string      `json:"option3,omitempty"`
	Position             int64       `json:"position,omitempty"`
	Price                string      `json:"price,omitempty"`
	ProductId            int64       `json:"product_id,omitempty"`
	RequiresShipping     bool        `json:"requires_shipping,omitempty"`
	Sku                  string      `json:"sku,omitempty"`
	Taxable              bool        `json:"taxable,omitempty"`
	Title                string      `json:"title,omitempty"`
	UpdatedAt            string      `json:"updated_at,omitempty"`
	ImageId              int64       `json:"image_id,omitempty"`
}
