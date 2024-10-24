package main

import (
	"fmt"
	"github.com/gomarkdown/markdown/parser"
	"github.com/tamjidrahman/slidedown/slides"
	"os"
)

func main() {
	const filename = "test.md"
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	parser := parser.NewWithExtensions(extensions)
	root := parser.Parse(content)
	slides := slides.GenerateSlidesFromMarkdownAST(root, slides.GenerateConfig{
		HeaderMapping: slides.HeaderMapping{
			Title:     1,
			Section:   2,
			Slide:     3,
			Paragraph: 4,
		},
		CreateTitleSlide:        true,
		CreateSectionTitleSlide: false,
	})

	fmt.Println(slides)

}
