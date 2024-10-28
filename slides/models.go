package slides

import (
	"fmt"
	"strings"
)

type Layout string

const (
	LayoutBody    Layout = "Body"
	LayoutTitle   Layout = "Title"
	LayoutSection Layout = "Section"
)

type GenerateConfig struct {
	HeaderMapping           HeaderMapping
	CreateTitleSlide        bool
	CreateSectionTitleSlide bool
}

type HeaderMapping struct {
	Title     int `json:"title"`
	Section   int `json:"section"`
	Slide     int `json:"slide"`
	Paragraph int `json:"paragraph"`
}

type Slides struct {
	Title    string
	Sections []*Section
}

type Section struct {
	Title  string
	Slides []*Slide
}

type Slide struct {
	Title      string
	Paragraphs []*Paragraph
	Layout     Layout
}

type Paragraph struct {
	Header string
	Image  string
	Text   string
}

func createUntitledSection() *Section {
	return &Section{
		Title:  "[untitled section]",
		Slides: []*Slide{},
	}
}

func createUntitledSlide() *Slide {
	return &Slide{
		Title:      "[untitled slide]",
		Layout:     LayoutBody,
		Paragraphs: []*Paragraph{},
	}
}

func createUntitledParagraph() *Paragraph {
	return &Paragraph{
		Header: "",
		Image:  "",
		Text:   "",
	}
}

func (slide *Slide) String() string {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("[%s] (%s)\n", slide.Title, slide.Layout))
	if len(slide.Paragraphs) > 0 {
		for k, paragraph := range slide.Paragraphs {
			result.WriteString(fmt.Sprintf("  Paragraph %d:\n", k+1))
			if paragraph.Header != "" {
				result.WriteString(fmt.Sprintf("    Header: %s\n", paragraph.Header))
			}
			if paragraph.Image != "" {
				result.WriteString(fmt.Sprintf("    Image: %s\n", paragraph.Image))
			}
			if paragraph.Text != "" {
				result.WriteString(fmt.Sprintf("    Text: %s\n", paragraph.Text))
			}
		}
	} else {
		result.WriteString("  No paragraphs\n")
	}

	return result.String()
}

func (slides *Slides) String() string {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("Slides Title: %s\n", slides.Title))
	slideCount := 1
	for i, section := range slides.Sections {
		result.WriteString(fmt.Sprintf("Section %d: %s\n", i+1, section.Title))
		for _, slide := range section.Slides {
			result.WriteString(fmt.Sprintf("  Slide %d: %s", slideCount, slide.String()))
			slideCount++
		}
	}

	return result.String()
}

func (slides *Slides) FlattenToSlides() []*Slide {
	var flattenedSlides []*Slide
	// Iterate through sections and their slides
	for _, section := range slides.Sections {
		// Add all slides in the section
		flattenedSlides = append(flattenedSlides, section.Slides...)
	}

	return flattenedSlides
}
