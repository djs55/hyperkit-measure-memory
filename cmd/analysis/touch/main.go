package main

import (
	"flag"
	"log"

	"github.com/djs55/hyperkit-measure-memory/pkg/gnuplot"
	"github.com/djs55/hyperkit-measure-memory/pkg/sample"
)

const (
	macOS1012 = 0
	macOS1014 = iota
)

var (
	touch = "touch"

	footprint1012GiB = int64(0)
)

func main() {
	flag.Int64Var(&footprint1012GiB, "footprint-10.12", int64(0), "manually observed footprint on 10.12")
	flag.Parse()

	doTouch(macOS1012)
	doTouch(macOS1014)
}

func doTouch(macOS int) {
	dir := getDir(macOS)

	VSZPoints, err := sample.ReadDir(dir, func(s sample.Sample) int64 {
		for _, command := range s.PS {
			if command.Command == touch {
				return command.VSZ
			}
		}
		return int64(0)
	})
	if err != nil {
		log.Fatal(err)
	}
	RSSPoints, err := sample.ReadDir(dir, func(s sample.Sample) int64 {
		for _, command := range s.PS {
			if command.Command == touch {
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
			Label:  "Resident Memory (RSS)",
			Points: RSSPoints,
		},
		&gnuplot.Line{
			Label:  "Virtual Size (VSZ)",
			Points: VSZPoints,
		},
	}
	footprintPoints, err := sample.ReadDir(dir, func(s sample.Sample) int64 {
		if macOS == macOS1014 {
			return int64(s.TouchFootprint)
		}
		return footprint1012GiB * int64(1024) * int64(1024) * int64(1024)
	})
	if err != nil {
		log.Fatal(err)
	}
	lines = append(lines, &gnuplot.Line{
		Label:  "physical footprint",
		Points: footprintPoints,
	})

	g := gnuplot.Graph{
		Title: "minimal C program memory usage, " + m,
		Lines: lines,
	}
	if err := g.Render("footprint-touch-" + dir + ".png"); err != nil {
		log.Fatalf("Failed to render: %v", err)
	}
}

func getDir(macOS int) string {
	m := "10.12"
	if macOS == macOS1014 {
		m = "10.14"
	}
	return m + "-touch"
}
