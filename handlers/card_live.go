package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

func CardLive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()

	var fields []Field
	for _, f := range q["field"] {
		parts := strings.SplitN(f, ":", 2)
		if len(parts) == 2 {
			fields = append(fields, makeField(parts[0], parts[1], ""))
		}
	}

	username := q.Get("username")
	if username != "" {
		stats, err := FetchGitHubStats(username)
		if err == nil {
			fields = append(fields,
				makeField("Repos",         fmt.Sprintf("%d", stats.TotalRepos),   ""),
				makeField("Commits",       fmt.Sprintf("%d", stats.TotalCommits), ""),
				makeField("Lines Added",   fmt.Sprintf("%d", stats.LinesAdded),   "green"),
				makeField("Lines Deleted", fmt.Sprintf("%d", stats.LinesDeleted), "red"),
			)
		}
	}

	input := defaultInput(CardInput{
		Username:   username,
		Hostname:   q.Get("hostname"),
		Background: q.Get("background"),
		KeyColor:   q.Get("keycolor"),
		TextColor:  q.Get("textcolor"),
		Fields:     fields,
	})

	var result string
	ok := run(w,
		func() error { return renderCard(input, &result) },
	)
	if !ok {
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Write([]byte(result))
}