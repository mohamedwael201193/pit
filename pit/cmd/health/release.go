package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/mohamedwael201193/pit/internal/obs"
)

const githubLatest = "https://api.github.com/repos/mohamedwael201193/pit/releases/latest"

type publicRelease struct {
	Tag      string `json:"tag"`
	Name     string `json:"name"`
	HTML     string `json:"html"`
	SHA      string `json:"sha,omitempty"`
	Unsigned bool   `json:"unsigned"`
	Asset    string `json:"asset,omitempty"`
	File     string `json:"file,omitempty"`
	Sums     string `json:"sums,omitempty"`
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
		"ok":        rel.Tag != "",
		"tag":       rel.Tag,
		"name":      rel.Name,
		"html":      rel.HTML,
		"sha":       rel.SHA,
		"installer": rel.Asset,
		"filename":  rel.File,
		"checksums": rel.Sums,
		"unsigned":  true,
		"sign":      false,
		"trade":     false,
	}
}

func windowsAsset(rel publicRelease) string {
	return strings.TrimSpace(rel.Asset)
}

func checksumsAsset(rel publicRelease) string {
	return strings.TrimSpace(rel.Sums)
}

func redirectLatestAsset(w http.ResponseWriter, r *http.Request, url, filename string) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	_, ok := cachedRelease()
	if !ok {
		obs.WriteJSON(w, http.StatusServiceUnavailable, map[string]any{
			"ok": false, "error": "release_unavailable", "sign": false, "trade": false,
		})
		return
	}
	dest := strings.TrimSpace(url)
	if dest == "" {
		obs.WriteJSON(w, http.StatusNotFound, map[string]any{
			"ok": false, "error": "installer_missing", "sign": false, "trade": false,
		})
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	if filename != "" {
		w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	}
	w.Header().Set("Location", dest)
	w.WriteHeader(http.StatusFound)
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
	req.Header.Set("User-Agent", "pit-health/0.9.13")
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
		name := strings.ToLower(strings.TrimSpace(a.Name))
		url := strings.TrimSpace(a.URL)
		if strings.HasSuffix(name, "x64-setup.exe") {
			out.Asset = url
			out.File = strings.TrimSpace(a.Name)
			continue
		}
		if strings.Contains(name, "sha256") {
			out.Sums = url
			out.SHA = fetchSumsSHA(url)
		}
	}
	return out, nil
}

func fetchSumsSHA(url string) string {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("User-Agent", "pit-health/0.9.13")
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
