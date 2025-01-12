package platform

import (
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/phcarneirobc/profile-search/internal/client"
	"github.com/phcarneirobc/profile-search/internal/config"
	"github.com/phcarneirobc/profile-search/internal/model"
)

type Platform struct {
	Name    string
	URL     string
	Valid   func(*http.Response) (bool, error)
	Extract func(*http.Response) (*model.ProfileInfo, error)
}

func extractGithubInfo(resp *http.Response) (*model.ProfileInfo, error) {
	body, err := client.ReadResponseBody(resp)
	if err != nil {
		return nil, err
	}

	info := &model.ProfileInfo{}

	nameRegex := regexp.MustCompile(`itemprop="name">([^<]+)</span>`)
	if matches := nameRegex.FindStringSubmatch(body); len(matches) > 1 {
		info.Name = matches[1]
	}

	locRegex := regexp.MustCompile(`itemprop="homeLocation">([^<]+)</span>`)
	if matches := locRegex.FindStringSubmatch(body); len(matches) > 1 {
		info.Location = matches[1]
	}

	bioRegex := regexp.MustCompile(`itemprop="description">([^<]+)</div>`)
	if matches := bioRegex.FindStringSubmatch(body); len(matches) > 1 {
		info.Bio = matches[1]
	}

	return info, nil
}

func GetPlatforms() []Platform {
	return []Platform{
		{
			Name: "Instagram",
			URL:  "https://www.instagram.com/%s",
			Valid: func(resp *http.Response) (bool, error) {
				if resp.StatusCode == 404 {
					return false, nil
				}
				body, err := client.ReadResponseBody(resp)
				if err != nil {
					return false, err
				}
				return !strings.Contains(body, "Esta página não está disponível") &&
					!strings.Contains(body, "page isn't available") &&
					!strings.Contains(body, "Sorry, this page") &&
					strings.Contains(body, "\"@type\":\"ProfilePage\""), nil
			},
		},
		{
			Name: "Twitter/X",
			URL:  "https://twitter.com/%s",
			Valid: func(resp *http.Response) (bool, error) {
				body, err := client.ReadResponseBody(resp)
				if err != nil {
					return false, err
				}
				return !strings.Contains(body, "This account doesn't exist") &&
					!strings.Contains(body, "Esta conta não existe"), nil
			},
		},
		{
			Name: "Facebook",
			URL:  "https://www.facebook.com/%s",
			Valid: func(resp *http.Response) (bool, error) {
				body, err := client.ReadResponseBody(resp)
				if err != nil {
					return false, err
				}
				return !strings.Contains(body, "Página não encontrada") &&
					!strings.Contains(body, "Page not found"), nil
			},
		},
		{
			Name: "LinkedIn",
			URL:  "https://www.linkedin.com/in/%s",
			Valid: func(resp *http.Response) (bool, error) {
				if resp.StatusCode == 404 || resp.StatusCode == 999 {
					return false, nil
				}
				body, err := client.ReadResponseBody(resp)
				if err != nil {
					return false, err
				}
				return !strings.Contains(body, "Page not found") &&
					!strings.Contains(body, "this page doesn't exist"), nil
			},
		},
		{
			Name: "TikTok",
			URL:  "https://www.tiktok.com/@%s",
			Valid: func(resp *http.Response) (bool, error) {
				body, err := client.ReadResponseBody(resp)
				if err != nil {
					return false, err
				}
				return !strings.Contains(body, "Não foi possível localizar esta conta") &&
					!strings.Contains(body, "Couldn't find this account"), nil
			},
		},
		{
			Name: "GitHub",
			URL:  "https://github.com/%s",
			Valid: func(resp *http.Response) (bool, error) {
				if resp.StatusCode == 404 {
					return false, nil
				}
				body, err := client.ReadResponseBody(resp)
				if err != nil {
					return false, err
				}
				return !strings.Contains(body, "Not Found") &&
					!strings.Contains(body, "404"), nil
			},
			Extract: extractGithubInfo,
		},
		{
			Name: "GitLab",
			URL:  "https://gitlab.com/%s",
			Valid: func(resp *http.Response) (bool, error) {
				return resp.StatusCode != 404, nil
			},
		},
		{
			Name: "Reddit",
			URL:  "https://www.reddit.com/user/%s",
			Valid: func(resp *http.Response) (bool, error) {
				body, err := client.ReadResponseBody(resp)
				if err != nil {
					return false, err
				}
				return !strings.Contains(body, "Sorry, nobody on Reddit goes by that name"), nil
			},
		},
		{
			Name: "Pinterest",
			URL:  "https://www.pinterest.com/%s/",
			Valid: func(resp *http.Response) (bool, error) {
				if resp.StatusCode == 404 {
					return false, nil
				}
				return true, nil
			},
		},
		{
			Name: "Steam",
			URL:  "https://steamcommunity.com/id/%s",
			Valid: func(resp *http.Response) (bool, error) {
				body, err := client.ReadResponseBody(resp)
				if err != nil {
					return false, err
				}
				return !strings.Contains(body, "The specified profile could not be found"), nil
			},
		},
		{
			Name: "Twitch",
			URL:  "https://www.twitch.tv/%s",
			Valid: func(resp *http.Response) (bool, error) {
				body, err := client.ReadResponseBody(resp)
				if err != nil {
					return false, err
				}
				return !strings.Contains(body, "esse conteúdo está indisponível") &&
					!strings.Contains(body, "content is unavailable"), nil
			},
		},
		{
			Name: "DeviantArt",
			URL:  "https://www.deviantart.com/%s",
			Valid: func(resp *http.Response) (bool, error) {
				return resp.StatusCode != 404, nil
			},
		},
		{
			Name: "Behance",
			URL:  "https://www.behance.net/%s",
			Valid: func(resp *http.Response) (bool, error) {
				return resp.StatusCode != 404, nil
			},
		},
		{
			Name: "Medium",
			URL:  "https://medium.com/@%s",
			Valid: func(resp *http.Response) (bool, error) {
				body, err := client.ReadResponseBody(resp)
				if err != nil {
					return false, err
				}
				return !strings.Contains(body, "PAGE NOT FOUND") &&
					!strings.Contains(body, "404"), nil
			},
		},
		{
			Name: "Tumblr",
			URL:  "https://%s.tumblr.com",
			Valid: func(resp *http.Response) (bool, error) {
				return resp.StatusCode != 404, nil
			},
		},
		{
			Name: "SoundCloud",
			URL:  "https://soundcloud.com/%s",
			Valid: func(resp *http.Response) (bool, error) {
				return resp.StatusCode != 404, nil
			},
		},
		{
			Name: "Spotify",
			URL:  "https://open.spotify.com/user/%s",
			Valid: func(resp *http.Response) (bool, error) {
				return resp.StatusCode != 404, nil
			},
		},
		{
			Name: "Telegram",
			URL:  "https://t.me/%s",
			Valid: func(resp *http.Response) (bool, error) {
				body, err := client.ReadResponseBody(resp)
				if err != nil {
					return false, err
				}
				return !strings.Contains(body, "If you have Telegram, you can view and join") &&
					resp.StatusCode != 404, nil
			},
		},
		{
			Name: "VK",
			URL:  "https://vk.com/%s",
			Valid: func(resp *http.Response) (bool, error) {
				body, err := client.ReadResponseBody(resp)
				if err != nil {
					return false, err
				}
				return !strings.Contains(body, "Page not found") &&
					!strings.Contains(body, "404"), nil
			},
		},
		{
			Name: "Patreon",
			URL:  "https://www.patreon.com/%s",
			Valid: func(resp *http.Response) (bool, error) {
				return resp.StatusCode != 404, nil
			},
		},
		{
			Name: "BitBucket",
			URL:  "https://bitbucket.org/%s/",
			Valid: func(resp *http.Response) (bool, error) {
				return resp.StatusCode != 404, nil
			},
		},
		{
			Name: "WordPress",
			URL:  "https://%s.wordpress.com",
			Valid: func(resp *http.Response) (bool, error) {
				return resp.StatusCode != 404, nil
			},
		},
	}
}

func CheckPlatform(p Platform, username string, cfg *config.Config, results chan<- model.Result) {
	result := model.Result{
		Platform: p.Name,
		URL:      fmt.Sprintf(p.URL, username),
	}

	startTime := time.Now()

	for retry := 0; retry < cfg.MaxRetries; retry++ {
		cl := client.CreateClient(cfg)
		req, err := http.NewRequest("GET", result.URL, nil)
		if err != nil {
			result.Error = "Erro ao criar requisição"
			continue
		}

		userAgent := cfg.UserAgents[rand.Intn(len(cfg.UserAgents))]
		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")

		resp, err := cl.Do(req)
		if err != nil {
			result.Error = "Erro na requisição"
			time.Sleep(time.Second * time.Duration(retry+1))
			continue
		}

		exists, err := p.Valid(resp)

		if exists && p.Extract != nil {
			if info, err := p.Extract(resp); err == nil {
				result.Info = info
			}
		}

		resp.Body.Close()
		result.ResponseTime = time.Since(startTime)

		if err != nil {
			result.Error = "Erro ao validar"
			continue
		}

		result.Exists = exists
		break
	}

	results <- result
}
