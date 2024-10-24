package slides

import (
	"fmt"
	"github.com/gomarkdown/markdown/ast"
	"strings"
)

type Layout int

const (
	LayoutBody Layout = iota
	LayoutTitle
	LayoutSection
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

func GenerateSlidesFromMarkdownAST(root ast.Node, config GenerateConfig) *Slides {

	mapping := config.HeaderMapping

	slides := &Slides{}
	var currentSection *Section
	var currentSlide *Slide
	var currentParagraph *Paragraph

	ast.Walk(root, ast.NodeVisitorFunc(func(node ast.Node, entering bool) ast.WalkStatus {
		if !entering {
			return ast.GoToNext
		}

		switch n := node.(type) {
		case *ast.Heading:
			headingText := string(n.Children[0].(*ast.Text).Literal)
			if n.Level == mapping.Title {
				slides.Title = headingText

				if currentSection == nil {
					currentSection = &Section{
						Title:  "Intro Section",
						Slides: []*Slide{},
					}
					slides.Sections = append(slides.Sections, currentSection)
				}

				currentSlide = nil
				currentParagraph = nil

				if config.CreateTitleSlide {
					currentSlide = &Slide{
						Title:      headingText,
						Layout:     LayoutTitle,
						Paragraphs: []*Paragraph{},
					}
					currentSection.Slides = append(currentSection.Slides, currentSlide)
				}
			}
			if n.Level == mapping.Section {
				currentSection = &Section{
					Title:  headingText,
					Slides: []*Slide{},
				}

				currentSlide = nil
				currentParagraph = nil

				if config.CreateSectionTitleSlide {
					currentSlide = &Slide{
						Title:      headingText,
						Layout:     LayoutSection,
						Paragraphs: []*Paragraph{},
					}
					currentSection.Slides = append(currentSection.Slides, currentSlide)
				}
				slides.Sections = append(slides.Sections, currentSection)
			}
			if n.Level == mapping.Slide {
				fmt.Println(headingText)
				currentSlide = &Slide{
					Title:      headingText,
					Layout:     LayoutBody,
					Paragraphs: []*Paragraph{},
				}

				if currentSection == nil {
					currentSection = createUntitledSection()
					slides.Sections = append(slides.Sections, currentSection)
				}

				currentParagraph = nil
				currentSection.Slides = append(currentSection.Slides, currentSlide)
			}
			if n.Level == mapping.Paragraph {
				currentParagraph = &Paragraph{
					Header: headingText,
				}
				if currentSection == nil {
					currentSection = createUntitledSection()
					slides.Sections = append(slides.Sections, currentSection)
				}

				if currentSlide == nil {
					currentSlide = createUntitledSlide()
					currentSection.Slides = append(currentSection.Slides, currentSlide)
				}
				currentSlide.Paragraphs = append(currentSlide.Paragraphs, currentParagraph)
			}
		case *ast.Image:
			if currentParagraph == nil {
				currentParagraph = createUntitledParagraph()
				currentSlide.Paragraphs = append(currentSlide.Paragraphs, currentParagraph)
			}
			currentParagraph.Image = string(n.Destination)

		case *ast.Paragraph:
			if currentSection == nil {
				currentSection = createUntitledSection()
				slides.Sections = append(slides.Sections, currentSection)
			}
			if currentSlide == nil {
				currentSlide = createUntitledSlide()
				currentSection.Slides = append(currentSection.Slides, currentSlide)
			}
			if currentParagraph == nil {
				currentParagraph = createUntitledParagraph()
				currentSlide.Paragraphs = append(currentSlide.Paragraphs, currentParagraph)
			}
			for _, child := range n.Children {
				if t, ok := child.(*ast.Text); ok {
					currentParagraph.Text += string(t.Literal)
				}
			}
		}

		return ast.GoToNext
	}))

	return slides
}

func (slides *Slides) String() string {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("Slides Title: %s\n", slides.Title))
	slideCount := 1
	for i, section := range slides.Sections {
		result.WriteString(fmt.Sprintf("Section %d: %s\n", i+1, section.Title))
		for _, slide := range section.Slides {
			result.WriteString(fmt.Sprintf("  Slide %d: [%s] (%s)\n", slideCount, slide.Title, layoutToString(slide.Layout)))
			slideCount++
			if len(slide.Paragraphs) > 0 {
				for k, paragraph := range slide.Paragraphs {
					result.WriteString(fmt.Sprintf("    Paragraph %d:\n", k+1))
					if paragraph.Header != "" {
						result.WriteString(fmt.Sprintf("      Header: %s\n", paragraph.Header))
					}
					if paragraph.Image != "" {
						result.WriteString(fmt.Sprintf("      Image: %s\n", paragraph.Image))
					}
					if paragraph.Text != "" {
						result.WriteString(fmt.Sprintf("      Text: %s\n", paragraph.Text))
					}
				}
			} else {
				result.WriteString("    No paragraphs\n")
			}
		}
	}

	return result.String()
}

func layoutToString(layout Layout) string {
	switch layout {
	case LayoutBody:
		return "Body"
	case LayoutTitle:
		return "Title"
	case LayoutSection:
		return "Section"
	default:
		return "Unknown"
	}
}
