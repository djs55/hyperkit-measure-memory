package main

import (
	"log"

	"github.com/djs55/hyperkit-measure-memory/pkg/gnuplot"
	"github.com/djs55/hyperkit-measure-memory/pkg/sample"
)

var (
	hyperkit = "com.docker.hyperkit"
)

func main() {
	doHyperkit()
}

func doHyperkit() {
	dir := getDir()

	RSSPoints, err := sample.ReadDir(dir, func(s sample.Sample) int64 {
		for _, command := range s.PS {
			if command.Command == hyperkit {
				return command.RSS
			}
		}
		return int64(0)
	})
	if err != nil {
		log.Fatal(err)
	}
	m := "10.14"
	lines := []*gnuplot.Line{
		&gnuplot.Line{
			Label:  "modified com.docker.hyperkit \"Real Mem\" in Activity Monitor",
			Points: RSSPoints,
		},
	}
	footprintPoints, err := sample.ReadDir(dir, func(s sample.Sample) int64 {
		return int64(s.Footprint)
	})
	if err != nil {
		log.Fatal(err)
	}
	lines = append(lines, &gnuplot.Line{
		Label:  "modified com.docker.hyperkit \"Memory\" in Activity Monitor",
		Points: footprintPoints,
	})

	g := gnuplot.Graph{
		Title:  "Memory usage with VM set to 4GB on " + m,
		Lines:  lines,
		Format: gnuplot.SVG,
	}
	if err := g.Render("footprint-" + dir + ".svg"); err != nil {
		log.Fatalf("Failed to render: %v", err)
	}
}

func getDir() string {
	m := "10.14"
	return m + "-hyperkit"
}
