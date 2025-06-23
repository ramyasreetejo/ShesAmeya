package util

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/biter777/countries"
)

func GetCountryFromIP(ip string) string {
	// Use a default mock Indian Cloudflare’s public IP for local development to simulate geo-based responses
	if ip == "::1" || ip == "127.0.0.1" || ip == "" {
		ip = "103.21.244.1"
	}
	res, _ := http.Get(fmt.Sprintf("https://ipinfo.io/%s/country", ip))
	body, _ := io.ReadAll(res.Body)
	code := strings.TrimSpace(string(body))
	return countries.ByName(code).Info().Name
}
