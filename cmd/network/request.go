package network

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// GetRedirectLink performs a GET request to the given URL with browser-like headers
// and returns the redirect download link from response headers.
// Renamed from MakeBrowserLikeRequest, made public.
func GetRedirectLink(downloadPageLink string) (string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	id := ""
	if u, err := url.Parse(downloadPageLink); err == nil {
		parts := strings.Split(u.Path, "/")
		if len(parts) > 2 {
			id = parts[1]
		}
	}

	req, err := http.NewRequest("GET", downloadPageLink, nil)
	if err != nil {
		return "", fmt.Errorf("could not create GET request: %w", err)
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("HX-Request", "true")
	if id != "" {
		req.Header.Set("HX-Current-URL", "https://buzzheavier.com/"+id+"/")
		req.Header.Set("Referer", "https://buzzheavier.com/"+id+"/")
	}
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-CH-UA", `"Not.A/Brand";v="99", "Chromium";v="136"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"Linux"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not fetch download page: %w", err)
	}
	defer resp.Body.Close()

	downloadLink := resp.Header.Get("Hx-Redirect")
	if downloadLink == "" {
		downloadLink = resp.Header.Get("HX-Redirect")
	}
	if downloadLink == "" {
		downloadLink = resp.Header.Get("Location")
	}
	if downloadLink == "" {
		downloadLink = resp.Header.Get("X-Redirect-To")
	}

	if downloadLink == "" {
		return "", fmt.Errorf("no redirect header found in response")
	}

	return downloadLink, nil
}
