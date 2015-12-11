package shopify

import (
	"bytes"
	"encoding/json"
	"fmt"

	"strconv"
	"time"
)

type Product struct {
	BodyHtml       string        `json:"body_html"`
	CreatedAt      time.Time     `json:"created_at"`
	Handle         string        `json:"handle"`
	Id             int64         `json:"id"`
	ProductType    string        `json:"product_type"`
	PublishedAt    time.Time     `json:"published_at"`
	PublishedScope string        `json:"published_scope"`
	TemplateSuffix string        `json:"template_suffix"`
	Title          string        `json:"title"`
	UpdatedAt      time.Time     `json:"updated_at"`
	Vendor         string        `json:"vendor"`
	Tags           string        `json:"tags"`
	Variants       []Variant     `json:"variants"`
	Options        []Option      `json:"options"`
	Images         []interface{} `json:"images"`
	api            *API
}

func (api *API) Products() ([]*Product, error) {
	res, status, err := api.request("/admin/products.json", "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := &map[string][]*Product{}
	err = json.NewDecoder(res).Decode(r)

	result := (*r)["products"]

	if err != nil {
		return nil, err
	}

	for _, v := range result {
		v.api = api
	}

	return result, nil
}

type ProductsCountOptions struct {
	Vendor          string `url:"vendor"`
	ProductType     string `url:"product_type"`
	CollectionID    string `url:"collection_id"`
	CreatedAtMin    string `url:"created_at_min"`
	CreatedAtMax    string `url:"created_at_max"`
	UpdatedAtMin    string `url:"updated_at_min"`
	UpdatedAtMax    string `url:"updated_at_max"`
	PublishedAtMin  string `url:"published_at_min"`
	PublishedAtMax  string `url:"published_at_max"`
	PublishedStatus string `url:"published_status"`
}

func (api *API) ProductsCount(options *ProductsCountOptions) (int, error) {

	// TODO: marshall options using go-querystring
	res, status, err := api.request("/admin/products/count.json", "GET", nil, nil)

	if err != nil {
		return 0, err
	}

	if status != 200 {
		return 0, fmt.Errorf("Status returned: %d", status)
	}

	r := map[string]interface{}{}
	err = json.NewDecoder(res).Decode(&r)

	result, _ := strconv.Atoi(fmt.Sprintf("%v", r["count"]))
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (api *API) Product(id int64) (*Product, error) {
	endpoint := fmt.Sprintf("/admin/products/%d.json", id)

	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := map[string]Product{}
	err = json.NewDecoder(res).Decode(&r)

	result := r["product"]

	if err != nil {
		return nil, err
	}

	result.api = api

	return &result, nil
}

func (api *API) NewProduct() *Product {
	return &Product{api: api}
}

func (obj *Product) Save() error {
	endpoint := fmt.Sprintf("/admin/products/%d.json", obj.Id)
	method := "PUT"
	expectedStatus := 201

	if obj.Id == 0 {
		endpoint = fmt.Sprintf("/admin/products.json")
		method = "POST"
		expectedStatus = 201
	}

	body := map[string]*Product{}
	body["product"] = obj

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(body)

	if err != nil {
		return err
	}

	res, status, err := obj.api.request(endpoint, method, nil, buf)

	if err != nil {
		return err
	}

	if status != expectedStatus {
		r := errorResponse{}
		err = json.NewDecoder(res).Decode(&r)
		if err == nil {
			return fmt.Errorf("Status %d: %v", status, r.Errors)
		} else {
			return fmt.Errorf("Status %d, and error parsing body: %s", status, err)
		}
	}

	r := map[string]Product{}
	err = json.NewDecoder(res).Decode(&r)

	if err != nil {
		return err
	}

	api := obj.api
	*obj = r["product"]
	obj.api = api

	return nil
}

func (obj *Product) Delete() error {
	endpoint := fmt.Sprintf("/admin/products/%d.json", obj.Id)
	method := "DELETE"
	expectedStatus := 200

	res, status, err := obj.api.request(endpoint, method, nil, nil)

	if err != nil {
		return err
	}

	if status != expectedStatus {
		r := errorResponse{}
		err = json.NewDecoder(res).Decode(&r)
		if err == nil {
			return fmt.Errorf("Status %d: %v", status, r.Errors)
		} else {
			return fmt.Errorf("Status %d, and error parsing body: %s", status, err)
		}
	}

	return nil
}
