package slides

import (
	"fmt"
	"log"
	"os"

	"github.com/unidoc/unioffice/common/license"
	"github.com/unidoc/unioffice/presentation"
	// "github.com/unidoc/unioffice/schema/soo/pml"
)

func init() {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		log.Fatalf("Error setting metered key: %v", err)
	}
}

func ConvertSlidesToPPTX(slides *Slides, templatePath, outputPath string) error {
	ppt, err := presentation.OpenTemplate(templatePath)
	if err != nil {
		return fmt.Errorf("unable to open template: %v", err)
	}
	defer ppt.Close()

	// Remove any existing slides
	for _, s := range ppt.Slides() {
		ppt.RemoveSlide(s)
	}

	for _, section := range slides.Sections {
		for _, slide := range section.Slides {
			var layout presentation.SlideLayout
			var err error

			switch slide.Layout {
			case LayoutTitle:
				layout, err = ppt.GetLayoutByName("Title Slide")
				if err != nil {
					layout, err = ppt.GetLayoutByName("Title and Caption")
				}
			case LayoutBody:
				layout, err = ppt.GetLayoutByName("Title and Content")
			case LayoutSection:
				layout, err = ppt.GetLayoutByName("Section Header")
			default:
				layout, err = ppt.GetLayoutByName("Title and Content")
			}

			if err != nil {
				return fmt.Errorf("error retrieving layout for slide: %v", err)
			}

			sld, err := ppt.AddDefaultSlideWithLayout(layout)
			if err != nil {
				return fmt.Errorf("error adding slide: %v", err)
			}

			// Set slide title
			titlePh, _ := sld.GetPlaceholderByIndex(0)
			titlePh.SetText(slide.Title)

			// Add content
			contentPh, _ := sld.GetPlaceholderByIndex(1)
			contentPh.ClearAll()

			for _, paragraph := range slide.Paragraphs {
				para := contentPh.AddParagraph()

				if paragraph.Header != "" {
					run := para.AddRun()
					run.SetText(paragraph.Header)
					run.Properties().SetBold(true)
					para.AddBreak()
				}

				if paragraph.Text != "" {
					run := para.AddRun()
					run.SetText(paragraph.Text)
				}

				if paragraph.Image != "" {
					// Note: Image handling is not implemented in this example
					// You would need to use the appropriate UniDoc functions to add images
					log.Printf("Image handling not implemented: %s", paragraph.Image)
				}

				para.AddBreak()
			}
		}
	}

	// Save the presentation
	err = ppt.SaveToFile(outputPath)
	if err != nil {
		return fmt.Errorf("error saving presentation: %v", err)
	}

	return nil
}
