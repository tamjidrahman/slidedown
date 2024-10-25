package slides

import (
	"fmt"
)

func GeneratePPTXFromSlides(slides *Slides, templatePath, outputPath string) error {
	err := ConvertSlidesToPPTX(slides, templatePath, outputPath)
	if err != nil {
		return fmt.Errorf("error generating ODP: %v", err)
	}
	return nil
}
