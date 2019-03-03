package main

import (
	"log"

	"github.com/djs55/hyperkit-measure-memory/pkg/gnuplot"
	"github.com/djs55/hyperkit-measure-memory/pkg/sample"
)

func main() {
	footprintPoints, err := sample.ReadDir("10.14-idle", func(s sample.Sample) int64 {
		return int64(s.Footprint)
	})
	if err != nil {
		log.Fatal(err)
	}
	VSZPoints, err := sample.ReadDir("10.14-idle", func(s sample.Sample) int64 {
		for _, command := range s.PS {
			if command.Command == "com.docker.hyperkit" {
				return command.VSZ
			}
		}
		return int64(0)
	})
	if err != nil {
		log.Fatal(err)
	}
	RSSPoints, err := sample.ReadDir("10.14-idle", func(s sample.Sample) int64 {
		for _, command := range s.PS {
			if command.Command == "com.docker.hyperkit" {
				return command.RSS
			}
		}
		return int64(0)
	})
	if err != nil {
		log.Fatal(err)
	}
	g := gnuplot.Graph{
		Title: "hyperkit physical footprint vs RSS vs VSZ, 10.14, idle Docker",
		Lines: []*gnuplot.Line{
			&gnuplot.Line{
				Label:  "physical footprint",
				Points: footprintPoints,
			},
			&gnuplot.Line{
				Label:  "Resident Memory (RSS)",
				Points: RSSPoints,
			},
			&gnuplot.Line{
				Label:  "Virtual Size (VSZ)",
				Points: VSZPoints,
			},
		},
	}
	if err := g.Render("/tmp/output.png"); err != nil {
		log.Fatalf("Failed to render: %v", err)
	}
}
