module github.com/tamjidrahman/slidedown

go 1.22.4

require (
	github.com/gomarkdown/markdown v0.0.0-20240930133441-72d49d9543d8
	github.com/tamjidrahman/slidedown/slides v0.0.0-00010101000000-000000000000
)

require github.com/unidoc/unioffice v1.37.0 // indirect

replace github.com/tamjidrahman/slides => /slides

replace github.com/tamjidrahman/slidedown/slides => ./slides
