package handlers

import (
	"encoding/base64"
	"encoding/json"
)

type cardConfig struct {
	Username   string   `json:"username"`
	Hostname   string   `json:"hostname"`
	Background string   `json:"background"`
	KeyColor   string   `json:"keycolor"`
	TextColor  string   `json:"textcolor"`
	ShowStats  bool     `json:"showstats"`
	Fields     []string `json:"fields"`
	Ascii      string   `json:"ascii"`
}

func encodeConfig(c cardConfig) string {
	data, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(data)
}
