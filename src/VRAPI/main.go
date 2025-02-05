package VRAPI

import (
	"VRCHAT/src/ConfigManager"
	"net/http"
)

var defaultHeaders = map[string]string{
	"Host":                             "vrchat.com",
	"User-Agent":                       "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:132.0) Gecko/20100101 Firefox/132.0",
	"Access-Control-Allow-Credentials": "true",
	"Access-Control-Allow-Origin":      "*",
}

var config = ConfigManager.Config

func ChangeDefaultHeader(req *http.Request) {
	for key, value := range defaultHeaders {
		req.Header.Set(key, value)
	}

	req.Header.Set("Cookie", "auth="+config.AuthCookie+"; twoFactorAuth="+config.TwoFactorAuth)
}
