package handlers

import (
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
			label := parts[0]
			value := parts[1]
			fields = append(fields, Field{
				Label: label,
				Dots: makeDots(label),
				Value: value,
			})
		}
	}

	input := defaultInput(CardInput{
		Username:   q.Get("username"),
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