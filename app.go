package shopify

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type App struct {
	APIKey          string
	APISecret       string
	RedirectURI     string
	IgnoreSignature bool
}

func (s *App) AuthorizeURL(shop string, scopes string) string {
	var u url.URL
	u.Scheme = "https"
	u.Host = shop
	u.Path = "/admin/oauth/authorize"
	q := u.Query()
	q.Set("client_id", s.APIKey)
	q.Set("scope", scopes)
	q.Set("redirect_uri", s.RedirectURI)
	u.RawQuery = q.Encode()

	return u.String()
}

func verifyHMAC(expectedHMAC, message, sharedSecret string) bool {
	h := hmac.New(sha256.New, []byte(sharedSecret))
	h.Write([]byte(message))

	value := hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(expectedHMAC), []byte(value))
}

func (s *App) VerifyHMACSignature(u *url.URL) bool {
	if s.IgnoreSignature {
		return true
	}
	params := u.Query()
	hmac := params.Get("hmac")
	//	if hmac == "" {
	//		return false
	//	}
	params.Del("hmac")
	params.Del("signature")
	message := s.signatureString(u, false)
	return verifyHMAC(hmac, message, s.APISecret)
}

// Deprecated and removed on 1st Jun 2015
//func (s *App) AdminSignatureOk(u *url.URL) bool {
//	if s.IgnoreSignature {
//		return true
//	}
//
//	params := u.Query()
//	signature := params.Get("signature")
//	if signature == "" {
//		return false
//	}
//
//	raw := md5.Sum([]byte(s.signatureString(u, true)))
//	encrypted := hex.EncodeToString(raw[:])
//
//	return 1 == subtle.ConstantTimeCompare([]byte(encrypted), []byte(signature))
//}

func (s *App) AppProxySignatureOk(u *url.URL) bool {
	if s.IgnoreSignature {
		return true
	}

	params := u.Query()
	signature := params.Get("signature")
	if signature == "" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(s.APISecret))
	mac.Write([]byte(s.signatureString(u, false)))
	calculated := hex.EncodeToString(mac.Sum(nil))

	return 1 == subtle.ConstantTimeCompare([]byte(signature), []byte(calculated))
}

func (s *App) signatureString(u *url.URL, prependSig bool) string {
	params := u.Query()

	keys := []string{}
	for k, _ := range params {
		if k != "signature" && k != "hmac" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	inputs := []string{}
	for _, k := range keys {
		inputs = append(inputs, fmt.Sprintf("%s=%s", k, params.Get(k)))
	}

	if prependSig {
		return fmt.Sprintf("%v%v", s.APISecret, strings.Join(inputs, "&"))
	}
	return strings.Join(inputs, "&")
}

func (s *App) AccessToken(shop string, code string) (string, error) {
	url := fmt.Sprintf("https://%s/admin/oauth/access_token.json", shop)

	data := map[string]string{
		"client_id":     s.APIKey,
		"client_secret": s.APISecret,
		"code":          code,
	}

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	transport := &http.Transport{}
	response, err := transport.RoundTrip(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	token := map[string]string{}
	err = json.NewDecoder(response.Body).Decode(&token)

	if err != nil {
		return "", err
	}

	if _, ok := token["error"]; ok {
		return "", fmt.Errorf("%s", token["error"])
	}

	if _, ok := token["access_token"]; !ok {
		return "", fmt.Errorf("access_token not found in response")
	}

	return token["access_token"], nil
}
