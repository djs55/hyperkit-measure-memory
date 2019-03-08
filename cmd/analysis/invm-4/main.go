package main

import (
	"log"

	"github.com/djs55/hyperkit-measure-memory/pkg/gnuplot"
	"github.com/djs55/hyperkit-measure-memory/pkg/sample"
)

const (
	macOS1012 = 0
	macOS1014 = iota
)

var (
	hyperkit = "com.docker.hyperkit"
)

func main() {
	doHyperkit(macOS1014)
}

func doHyperkit(macOS int) {
	dir := getDir(macOS)

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
	m := "10.12"
	if macOS == macOS1014 {
		m = "10.14"
	}
	lines := []*gnuplot.Line{
		&gnuplot.Line{
			Label:  "com.docker.hyperkit \"Real Mem\" in Activity Monitor",
			Points: RSSPoints,
		},
	}
	if macOS == macOS1014 {
		footprintPoints, err := sample.ReadDir(dir, func(s sample.Sample) int64 {
			return int64(s.Footprint)
		})
		if err != nil {
			log.Fatal(err)
		}
		lines = append(lines, &gnuplot.Line{
			Label:  "com.docker.hyperkit \"Memory\" in Activity Monitor",
			Points: footprintPoints,
		})
	}

	g := gnuplot.Graph{
		Title: "Memory usage with VM set to 4GB on macOS " + m,
		Lines: lines,
	}
	if err := g.Render("footprint-" + dir + ".png"); err != nil {
		log.Fatalf("Failed to render: %v", err)
	}
}

func getDir(macOS int) string {
	m := "10.12"
	if macOS == macOS1014 {
		m = "10.14"
	}
	return m + "-invm-touch-4"
}
