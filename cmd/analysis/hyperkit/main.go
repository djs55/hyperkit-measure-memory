package main

import (
	"log"
	"path/filepath"

	"github.com/djs55/hyperkit-measure-memory/pkg/gnuplot"
	"github.com/djs55/hyperkit-measure-memory/pkg/sample"
)

var (
	hyperkit     = "com.docker.hyperkit"
	realMemLabel = "modified com.docker.hyperkit \"Real Mem\" in Activity Monitor"
	memLabel     = "modified com.docker.hyperkit \"Memory\" in Activity Monitor"

	macOS1012 = 0
	macOS1014 = 1
)

func main() {
	doHyperkit1012()
}

func doHyperkit1014() {
	dir := getDir(macOS1014)

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
	m := getMacOS(macOS1014)
	lines := []*gnuplot.Line{
		&gnuplot.Line{
			Label:  realMemLabel,
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
		Label:  memLabel,
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

func doHyperkit1012() {
	dir := getDir(macOS1012)

	RSSPoints, err := sample.ReadDir(filepath.Join(dir, "auto"), func(s sample.Sample) int64 {
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
	m := getMacOS(macOS1012)
	lines := []*gnuplot.Line{
		&gnuplot.Line{
			Label:  realMemLabel,
			Points: RSSPoints,
		},
	}
	footprintPoints, err := gnuplot.ReadPoints(filepath.Join(dir, "manual", "memory.dat"))
	if err != nil {
		log.Fatal(err)
	}
	// The footprint data is manually captured from a quicktime screen capture (!)
	// and has to be manually synchronised with the automatically gathered data.
	for _, point := range footprintPoints {
		point.Second = point.Second - 6
	}
	lines = append(lines, &gnuplot.Line{
		Label:  memLabel,
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

func getDir(macOS int) string {
	switch macOS {
	case macOS1014:
		return "10.14-hyperkit"
	case macOS1012:
		return "10.12-hyperkit"
	default:
		return "unknown"
	}
}

func getMacOS(macOS int) string {
	switch macOS {
	case macOS1012:
		return "10.12"
	case macOS1014:
		return "10.14"
	default:
		return "unknown"
	}
}
