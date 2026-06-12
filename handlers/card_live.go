package handlers

import "net/http"

func CardLive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()

	fields := []Field{
		{Label: "OS",     Dots: makeDots("OS"),     Value: q.Get("os")},
		{Label: "Kernel", Dots: makeDots("Kernel"), Value: q.Get("kernel")},
		{Label: "Shell",  Dots: makeDots("Shell"),  Value: q.Get("shell")},
		{Label: "Editor", Dots: makeDots("Editor"), Value: q.Get("editor")},
		{Label: "Uptime", Dots: makeDots("Uptime"), Value: q.Get("uptime")},
	}

	// filter out empty fields
	var filtered []Field
	for _, f := range fields {
		if f.Value != "" {
			filtered = append(filtered, f)
		}
	}

	input := defaultInput(CardInput{
		Username:   q.Get("username"),
		Hostname:   q.Get("hostname"),
		Background: q.Get("background"),
		KeyColor:   q.Get("keycolor"),
		TextColor:  q.Get("textcolor"),
		Fields:     filtered,
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