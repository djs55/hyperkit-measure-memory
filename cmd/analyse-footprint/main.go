package main

import (
	"log"

	"github.com/djs55/hyperkit-measure-memory/pkg/gnuplot"
	"github.com/djs55/hyperkit-measure-memory/pkg/sample"
)

func main() {
	points, err := sample.ReadDir("10.14-idle", func(s sample.Sample) int64 {
		return int64(s.Footprint)
	})
	if err != nil {
		log.Fatal(err)
	}
	footprint := gnuplot.Line{
		Label:  "physical footprint",
		Points: points,
	}
	g := gnuplot.Graph{
		Title: "hyperkit physical footprint",
		Lines: []*gnuplot.Line{
			&footprint,
		},
	}
	if err := g.Render("/tmp/output.png"); err != nil {
		log.Fatalf("Failed to render: %v", err)
	}
}
