package slides

import (
	"github.com/gomarkdown/markdown/ast"
)

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
