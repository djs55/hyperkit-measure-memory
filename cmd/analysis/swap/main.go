package main

import (
	"log"

	"github.com/djs55/hyperkit-measure-memory/pkg/gnuplot"
	"github.com/djs55/hyperkit-measure-memory/pkg/sample"
)

var (
	hyperkit = "com.docker.hyperkit"
	qcow2 = "qcow2"
	raw = "raw"
)

func main() {
	doSwap(qcow2)
	doSwap(raw)
}

func doSwap(disk string) {
	dir := "10.14-" + disk + "-swap"

	VSZPoints, err := sample.ReadDir(dir, func(s sample.Sample) int64 {
		for _, command := range s.PS {
			if command.Command == hyperkit {
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
			if command.Command == hyperkit {
				return command.RSS
			}
		}
		return int64(0)
	})
	if err != nil {
		log.Fatal(err)
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
		return int64(s.Footprint)
	})
	if err != nil {
		log.Fatal(err)
	}
	lines = append(lines, &gnuplot.Line{
		Label:  "physical footprint",
		Points: footprintPoints,
	})

	g := gnuplot.Graph{
		Title: "adding swap to hyperkit dynamically on 10.14 on " + disk,
		Lines: lines,
	}
	if err := g.Render("footprint-" + dir + ".png"); err != nil {
		log.Fatalf("Failed to render: %v", err)
	}
}

