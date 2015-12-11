package shopify

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

var api API
var remoteEnabled = false

func init() {

	if os.Getenv("SHOPIFY_API_PERM_TOKEN") != "" && os.Getenv("SHOPIFY_API_SHOP") != "" {
		remoteEnabled = true
		api = API{
			Shop:        os.Getenv("SHOPIFY_API_SHOP"),
			AccessToken: os.Getenv("SHOPIFY_API_PERM_TOKEN"),
		}
	} else if os.Getenv("SHOPIFY_API_TOKEN") != "" && os.Getenv("SHOPIFY_API_SECRET") != "" && os.Getenv("SHOPIFY_API_SHOP") != "" {
		remoteEnabled = true
		api = API{
			Shop:   os.Getenv("SHOPIFY_API_SHOP"),
			Token:  os.Getenv("SHOPIFY_API_TOKEN"),
			Secret: os.Getenv("SHOPIFY_API_SECRET"),
		}
	} else {
		log.Printf("Remote tests disabled, set SHOPIFY_API_KEY, SHOPIFY_API_SECRET, SHOPIFY_API_HOST, SHOPIFY_API_PERM_TOKEN")
	}
}

func TestProductsCount(t *testing.T) {
	if !remoteEnabled {
		return
	}

	_, err := api.ProductsCount(nil)
	if err != nil {
		t.Errorf("Error fetching products count: %v", err)
	}
}

func TestListCreateGetDeleteWebhook(t *testing.T) {
	if !remoteEnabled {
		return
	}

	// List and delete all
	webhooks, err := api.Webhooks()

	if err != nil {
		fmt.Printf("Err fetching webhooks: %v", err)
	}

	for _, v := range webhooks {
		fmt.Printf("Existing webhook: %#v", v)
		v.Delete()
	}

	// create
	newHook := api.NewWebhook()

	newHook.Address = "https://aaa.ngrok.com/service/hook"
	newHook.Format = "json"
	newHook.Topic = "orders/delete"
	err = newHook.Save()
	if err != nil {
		t.Fatalf("Error creating webhook: %v", err)
	}

	//get
	hook, err := api.Webhook(newHook.Id)
	if err != nil {
		t.Errorf("Error fetching webhook (%v): %v", newHook.Id, err)
	}

	if hook.Id != newHook.Id {
		t.Errorf("Expected retrieved webhook to have the same ID as newly created webhook")
	}

	// clean up
	err = newHook.Delete()
	if err != nil {
		t.Errorf("Error deleting webhook: %s", err)
	}
}

func TestListCreateGetDeleteProduct(t *testing.T) {
	if !remoteEnabled {
		return
	}

	// List and delete all
	products, err := api.Products()
	if err != nil {
		fmt.Printf("Err fetching products: %v", err)
	}

	for _, v := range products {
		fmt.Printf("Existing products: %#v", v)
		v.Delete()
	}

	// create
	newProduct := api.NewProduct()
	newProduct.Title = "T-shirt"
	newProduct.PublishedAt = time.Now()
	newProduct.ProductType = "shirts"
	err = newProduct.Save()
	if err != nil {
		t.Fatalf("Error saving product: %s", err)
	}
	if newProduct.Id == 0 {
		t.Errorf("Missing ID for newly created product")
	}

	// get new product by id
	product, err := api.Product(newProduct.Id)

	if err != nil {
		t.Errorf("Error fetching product (%v): %v", newProduct.Id, err)
	}

	if product.Id != newProduct.Id {
		t.Errorf("Expected retrieved product to have the same ID as newly created product")
	}

	// clean up
	err = newProduct.Delete()
	if err != nil {
		t.Errorf("Error deleting product: %s", err)
	}
}
