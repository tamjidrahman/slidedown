package main

import (
	"fmt"
	"io"
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

func generateSlidesFromMarkdownContent(content []byte) *slides.Slides {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	parser := parser.NewWithExtensions(extensions)
	root := parser.Parse(content)
	return slides.GenerateSlidesFromMarkdownAST(root, slides.GenerateConfig{
		HeaderMapping: slides.HeaderMapping{
			Title:     1,
			Section:   2,
			Slide:     3,
			Paragraph: 4,
		},
		CreateTitleSlide:        true,
		CreateSectionTitleSlide: true,
	})
}

func generateSlides(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	content, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	generatedSlides := generateSlidesFromMarkdownContent(content)

	fmt.Fprint(w, generatedSlides)
}

func generatePowerpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	content, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	generatedSlides := generateSlidesFromMarkdownContent(content)
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
