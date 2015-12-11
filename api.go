package shopify

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type API struct {
	Shop        string // for e.g. demo-3.myshopify.com
	AccessToken string // permanent store access token
	Token       string // API client token
	Secret      string // API client secret for this shop
	client      *http.Client
}

type errorResponse struct {
	Errors map[string]interface{} `json:"errors"`
}

func (api *API) request(endpoint string, method string, params map[string]interface{}, body io.Reader) (result *bytes.Buffer, status int, err error) {
	if api.client == nil {
		api.client = &http.Client{}
	}

	uri := fmt.Sprintf("https://%s%s", api.Shop, endpoint)
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return
	}

	if api.AccessToken != "" {
		req.Header.Set("X-Shopify-Access-Token", api.AccessToken)
	} else {
		req.SetBasicAuth(api.Token, api.Secret)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := api.client.Do(req)
	if err != nil {
		return
	}

	status = resp.StatusCode

	result = &bytes.Buffer{}
	if _, err = io.Copy(result, resp.Body); err != nil {
		return
	}

	return
}
