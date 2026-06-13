package handlers

import (
	"net/http"
)

func CardLive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()

	username := q.Get("username")
	showStats := q.Get("showstats") != "false"

	fields := parseFields(q["field"])
	fields = appendGitHubStats(fields, username, showStats)

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