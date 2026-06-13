package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type GitHubStats struct {
	TotalRepos	 int
	TotalCommits int
	TotalLines   int
	LinesAdded   int
	LinesDeleted int
}

type contributionWeek struct {
	W int `json:"w"`
	A int `json:"a"`
	D int `json:"d"`
	C int `json:"c"`
}

type contributorStats struct {
	Total int					`json:"total"`
	Weeks []contributionWeek	`json:"weeks"`
}

type statsCacheEntry struct {
	stats   GitHubStats
	err     error
	expires time.Time
}

var (
	statsCacheMu sync.Mutex
	statsCache   = map[string]statsCacheEntry{}
)

const (
	statsSuccessTTL = 1 * time.Hour
	statsErrorTTL   = 1 * time.Minute
	statsPendingTTL = 30 * time.Second
)

func FetchGitHubStats(username string) (GitHubStats, error) {
	statsCacheMu.Lock()
	if e, ok := statsCache[username]; ok && time.Now().Before(e.expires) {
		statsCacheMu.Unlock()
		return e.stats, e.err
	}
	statsCacheMu.Unlock()

	stats, err := fetchGitHubStats(username)

	ttl := statsSuccessTTL
	if err != nil {
		ttl = statsErrorTTL
	} else if stats.TotalRepos > 0 && stats.TotalCommits == 0 && stats.LinesAdded == 0 {
		ttl = statsPendingTTL
	}

	statsCacheMu.Lock()
	statsCache[username] = statsCacheEntry{stats: stats, err: err, expires: time.Now().Add(ttl)}
	statsCacheMu.Unlock()

	return stats, err
}

func githubGet(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return http.DefaultClient.Do(req)
}

func fetchGitHubStats(username string) (GitHubStats, error) {
	var stats GitHubStats

	repos, err := fetchRepos(username)

	if err != nil {
		return stats, err
	}

	stats.TotalRepos = len(repos)

	for _, repo := range repos {
		repoStats, err := fetchRepoStats(username, repo)
		
		if err != nil {
			continue
		}

		stats.LinesAdded += repoStats.LinesAdded
		stats.LinesDeleted += repoStats.LinesDeleted
		stats.TotalCommits += repoStats.TotalCommits
	}

	stats.TotalLines = stats.LinesAdded - stats.LinesDeleted
	return stats, nil
}

func fetchRepos(username string) ([]string, error) {
	url := fmt.Sprintf(
		"https://api.github.com/users/%s/repos?per_page=100&type=owner",
		username,
	)

	resp, err := githubGet(url)

	if err != nil {
		return nil, err
	}
	
	defer resp.Body.Close()

	var repos []struct {
		Name string `json:"name"`
	}
	
	if err = json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}

	var names []string
	
	for _, r := range repos {
		names = append(names, r.Name)
	
	}
	return names, nil
}

func fetchRepoStats(username, repo string) (GitHubStats, error) {
	var stats GitHubStats

	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/stats/contributors",
		username, repo,
	)

	resp, err := githubGet(url)

	if err != nil {
		return stats, err
	}
	
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusAccepted {
		return stats, fmt.Errorf("stats not ready")
	}

	if resp.StatusCode != http.StatusOK {
		return stats, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var contributors []contributorStats
	
	if err = json.NewDecoder(resp.Body).Decode(&contributors); err != nil {
		return stats, err
	}

	for _, c := range contributors {
		for _, w := range c.Weeks {
			stats.LinesAdded += w.A
			stats.LinesDeleted += w.D
			stats.TotalCommits += w.C
		}
	}

	return stats, nil
}