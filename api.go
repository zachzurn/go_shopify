package shopify

import (
	"bytes"
	"fmt"
	"github.com/jpillora/backoff"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const REFILL_RATE = float64(0.5) // 2 per second
const BUCKET_LIMIT = 40
const BUCKET_SLOWDOWN = 30
const MAX_RETRIES = 3

type API struct {
	Shop        string // for e.g. demo-3.myshopify.com
	AccessToken string // permanent store access token
	Token       string // API client token
	Secret      string // API client secret for this shop
	client      *http.Client

	callLimit  int
	callsMade  int
	backoff    *backoff.Backoff
	retryCount int
}

type errorResponse struct {
	Errors map[string]interface{} `json:"errors"`
}

func (api *API) request(endpoint string, method string, params map[string]interface{}, body io.Reader) (result *bytes.Buffer, status int, err error) {
	if api.client == nil {
		api.client = &http.Client{}
	}
	if api.backoff == nil {
		api.backoff = &backoff.Backoff{
			//These are the defaults
			Min:    100 * time.Millisecond,
			Max:    2 * time.Second,
			Jitter: true,
		}
	}
	if api.callLimit == 0 {
		api.callLimit = BUCKET_LIMIT
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

	calls, total := parseAPICallLimit(resp.Header.Get("HTTP_X_SHOPIFY_SHOP_API_CALL_LIMIT"))
	api.callsMade = calls
	api.callLimit = total

	status = resp.StatusCode
	if status == 429 { // statusTooManyRequests
		if api.retryCount < MAX_RETRIES {
			api.retryCount = api.retryCount + 1
			b := api.backoff.Duration()
			time.Sleep(b)
			// try again
			return api.request(endpoint, method, params, body)
		}
		// else just return
	}

	result = &bytes.Buffer{}
	if _, err = io.Copy(result, resp.Body); err != nil {
		return
	}
	return
}

func parseAPICallLimit(str string) (int, int) {
	tokens := strings.Split(str, "/")
	if len(tokens) != 2 {
		return 0, 0
	}
	calls, err := strconv.Atoi(tokens[0])
	if err != nil {
		return 0, 0
	}
	total, err := strconv.Atoi(tokens[1])
	if err != nil {
		return 0, 0
	}
	return calls, total
}
