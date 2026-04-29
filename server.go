package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"ascii-art/ascii"
)

type GenerateRequest struct {
	Text  string `json:"text"`
	Font  string `json:"font"`
	Color string `json:"color"`
	Align string `json:"align"`
}

type GenerateResponse struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func startServer(port string) {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/generate", handleGenerate)

	fmt.Printf("🚀 ASCII Art Web Server running at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(GenerateResponse{Error: "Method not allowed"})
		return
	}

	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(GenerateResponse{Error: "Invalid request body"})
		return
	}

	if strings.TrimSpace(req.Text) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(GenerateResponse{Error: "Text is required"})
		return
	}

	font := req.Font
	if font == "" {
		font = "standard"
	}

	align := req.Align
	if align == "" {
		align = ascii.AlignLeft
	}
	if align != "" && !ascii.IsValidAlignment(align) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(GenerateResponse{Error: "Invalid alignment: " + align})
		return
	}

	colorRules := make(map[string]string)
	if req.Color != "" {
		colorRules[""] = req.Color
	}

	result, err := ascii.GenerateAsciiArt(req.Text, font, colorRules, align, 80)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(GenerateResponse{Error: err.Error()})
		return
	}

	// Strip ANSI codes for web display since browser doesn't render them
	result = stripANSI(result)

	json.NewEncoder(w).Encode(GenerateResponse{Result: result})
}

func stripANSI(s string) string {
	var b strings.Builder
	inEscape := false
	for _, r := range s {
		if r == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
