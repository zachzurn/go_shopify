package shopify

import (
	"net/url"
	"testing"
)

var app App

func init() {
	app = App{APIKey: "asdf", APISecret: "1234", RedirectURI: "http://localhost:4000"}
}

func TestAuthorizeURL(t *testing.T) {
	redir := app.AuthorizeURL("burnsmod.myshopify.com", "read_orders")

	expected := "https://burnsmod.myshopify.com/admin/oauth/authorize?client_id=asdf&redirect_uri=http%3A%2F%2Flocalhost%3A4000&scope=read_orders"

	if redir != expected {
		t.Errorf("Expected %s, got %s", expected, redir)
	}
}

func TestSignatureString(t *testing.T) {
	u, _ := url.Parse("https://app.com/?shop=burnsmod.myshopify.com&code=asdf&timestamp=1337178173&signature=31b9fcfbd98a3650b8523bcc92f8c5d2")
	expected := "code=asdf&shop=burnsmod.myshopify.com&timestamp=1337178173"
	expected_prepended := "1234code=asdf&shop=burnsmod.myshopify.com&timestamp=1337178173"

	if output := app.signatureString(u, true); output != expected_prepended {
		t.Errorf("expected %s output %s", expected, output)
	}

	if output := app.signatureString(u, false); output != expected {
		t.Errorf("expected %s output %s", expected, output)
	}
}

func TestVerifyHMACSignature(t *testing.T) {
	u, _ := url.Parse("https://app.com/?hmac=89f3e7c84239719b8f5dc2c7bf743e884d7dc49c9de2847d3e3c38d1b4ad7b93&shop=burnsmod.myshopify.com&code=asdf&timestamp=1337178173&signature=bd28a1a098688d8937e991aef3bc80ab")

	if app.VerifyHMACSignature(u) != true {
		t.Errorf("signature checking failed")
	}
}

func TestIgnoreSignature(t *testing.T) {

	a := App{APIKey: "asdf", APISecret: "1234", RedirectURI: "http://localhost:4000", IgnoreSignature: true}

	u, _ := url.Parse("https://app.com/?shop=burnsmod.myshopify.com&code=asdf&timestamp=1337178173&signature=ffff")

	if a.VerifyHMACSignature(u) != true {
		t.Errorf("IgnoreSignature didn't work for Admin")
	}

	if a.AppProxySignatureOk(u) != true {
		t.Errorf("IgnoreSignature didn't work for AppProxy")
	}
}
