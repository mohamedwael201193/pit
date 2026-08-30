package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const githubLatest = "https://api.github.com/repos/mohamedwael201193/pit/releases/latest"

type publicRelease struct {
	Tag      string `json:"tag"`
	Name     string `json:"name"`
	HTML     string `json:"html"`
	SHA      string `json:"sha,omitempty"`
	Unsigned bool   `json:"unsigned"`
}

var (
	relMu   sync.Mutex
	relAt   time.Time
	relVal  publicRelease
	relOK   bool
	fetchGH = fetchGitHubLatest
)

func releaseBody(rel publicRelease) map[string]any {
	return map[string]any{
		"ok":       rel.Tag != "",
		"tag":      rel.Tag,
		"name":     rel.Name,
		"html":     rel.HTML,
		"sha":      rel.SHA,
		"unsigned": true,
		"sign":     false,
		"trade":    false,
	}
}

func cachedRelease() (publicRelease, bool) {
	relMu.Lock()
	defer relMu.Unlock()
	if relOK && time.Since(relAt) < 5*time.Minute {
		return relVal, true
	}
	got, err := fetchGH()
	if err != nil {
		if relOK {
			return relVal, true
		}
		return publicRelease{}, false
	}
	relVal = got
	relAt = time.Now()
	relOK = true
	return got, true
}

func fetchGitHubLatest() (publicRelease, error) {
	req, err := http.NewRequest(http.MethodGet, githubLatest, nil)
	if err != nil {
		return publicRelease{}, err
	}
	req.Header.Set("User-Agent", "pit-health/0.9.3")
	req.Header.Set("Accept", "application/vnd.github+json")
	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return publicRelease{}, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return publicRelease{}, err
	}
	if resp.StatusCode != 200 {
		return publicRelease{}, fmt.Errorf("github_%d", resp.StatusCode)
	}
	var body struct {
		TagName string `json:"tag_name"`
		Name    string `json:"name"`
		HTML    string `json:"html_url"`
		Assets  []struct {
			Name string `json:"name"`
			URL  string `json:"browser_download_url"`
		} `json:"assets"`
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		return publicRelease{}, err
	}
	out := publicRelease{
		Tag:      strings.TrimSpace(body.TagName),
		Name:     strings.TrimSpace(body.Name),
		HTML:     strings.TrimSpace(body.HTML),
		Unsigned: true,
	}
	for _, a := range body.Assets {
		if !strings.Contains(strings.ToLower(a.Name), "sha256") {
			continue
		}
		out.SHA = fetchSumsSHA(a.URL)
		break
	}
	return out, nil
}

func fetchSumsSHA(url string) string {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("User-Agent", "pit-health/0.9.3")
	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<16))
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(raw), "\n") {
		if !strings.Contains(strings.ToLower(line), "setup.exe") && !strings.Contains(strings.ToLower(line), "x64") {
			continue
		}
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) > 0 && len(fields[0]) >= 64 {
			return strings.ToUpper(fields[0])
		}
	}
	return ""
}
