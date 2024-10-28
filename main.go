package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gomarkdown/markdown/parser"
	"github.com/tamjidrahman/slidedown/slides"
)

func main() {
	http.HandleFunc("/", serveHTML)
	http.HandleFunc("/generate-slides", generateSlides)
	http.HandleFunc("/generate-pptx", generatePowerpoint)

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func generateSlidesFromMarkdownContent(content []byte, config slides.GenerateConfig) *slides.Slides {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	parser := parser.NewWithExtensions(extensions)
	root := parser.Parse(content)
	return slides.GenerateSlidesFromMarkdownAST(root, config)
}

func generateSlides(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		Content string                `json:"content"`
		Config  slides.GenerateConfig `json:"config,omitempty"`
		Debug   bool                  `json:"debug"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	config := requestBody.Config
	if config.HeaderMapping == (slides.HeaderMapping{}) {
		config = slides.GenerateConfig{
			HeaderMapping: slides.HeaderMapping{
				Title:     1,
				Section:   2,
				Slide:     3,
				Paragraph: 4,
			},
			CreateTitleSlide:        true,
			CreateSectionTitleSlide: true,
		}
	}

	generatedSlides := generateSlidesFromMarkdownContent([]byte(requestBody.Content), config)

	response := map[string]interface{}{
		"slides": generatedSlides.FlattenToSlides(),
	}

	if requestBody.Debug {
		response["debug"] = generatedSlides
	}

	json.NewEncoder(w).Encode(response)
}

func generatePowerpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		Content string                `json:"content"`
		Config  slides.GenerateConfig `json:"config,omitempty"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	config := requestBody.Config
	if config.HeaderMapping == (slides.HeaderMapping{}) {
		config = slides.GenerateConfig{
			HeaderMapping: slides.HeaderMapping{
				Title:     1,
				Section:   2,
				Slide:     3,
				Paragraph: 4,
			},
			CreateTitleSlide:        true,
			CreateSectionTitleSlide: true,
		}
	}

	generatedSlides := generateSlidesFromMarkdownContent([]byte(requestBody.Content), config)
	outputPath := "slides/output2.pptx"
	err = slides.GeneratePPTXFromSlides(generatedSlides, "slides/template2.pptx", outputPath)
	if err != nil {
		http.Error(w, "Error generating PowerPoint file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=output.pptx")
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.presentationml.presentation")
	http.ServeFile(w, r, outputPath)
}
