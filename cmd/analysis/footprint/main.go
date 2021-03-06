package main

import (
	"log"

	"github.com/djs55/hyperkit-measure-memory/pkg/gnuplot"
	"github.com/djs55/hyperkit-measure-memory/pkg/sample"
)

const (
	docker = 0
	k8s    = iota

	macOS1012 = 0
	macOS1014 = iota
)

var (
	hyperkit = "com.docker.hyperkit"
	firefox  = "firefox"
)

func main() {
	doHyperkit(docker, macOS1014)
	doFirefox(docker, macOS1014)
	doHyperkit(k8s, macOS1014)
	doFirefox(k8s, macOS1014)
	doHyperkit(docker, macOS1012)
	doFirefox(docker, macOS1012)
	doHyperkit(k8s, macOS1012)
	doFirefox(k8s, macOS1012)
}

func doHyperkit(running, macOS int) {
	dir := getDir(running, macOS)

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
	t := "Docker"
	if running == k8s {
		t = "Docker and Kubernetes"
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
		Title: "Memory usage with VM set to 2GB on macOS " + m + " with " + t,
		Lines: lines,
		Time:  gnuplot.Hours,
	}
	if err := g.Render("footprint-hyperkit-" + dir + ".png"); err != nil {
		log.Fatalf("Failed to render: %v", err)
	}
}

func doFirefox(running, macOS int) {
	dir := getDir(running, macOS)
	VSZPoints, err := sample.ReadDir(dir, func(s sample.Sample) int64 {
		for _, command := range s.PS {
			if command.Command == firefox {
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
			if command.Command == firefox {
				return command.RSS
			}
		}
		return int64(0)
	})
	if err != nil {
		log.Fatal(err)
	}
	t := "Docker"
	if running == k8s {
		t = "k8s"
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
	if macOS == macOS1014 {
		footprintPoints, err := sample.ReadDir(dir, func(s sample.Sample) int64 {
			return int64(s.FirefoxFootprint)
		})
		if err != nil {
			log.Fatal(err)
		}
		lines = append(lines, &gnuplot.Line{
			Label:  "physical footprint",
			Points: footprintPoints,
		})
	}
	g := gnuplot.Graph{
		Title: "firefox memory usage, " + m + " + " + t,
		Lines: lines,
	}
	if err := g.Render("footprint-firefox-" + dir + ".png"); err != nil {
		log.Fatalf("Failed to render: %v", err)
	}
}

func getDir(running, macOS int) string {
	m := "10.12"
	if macOS == macOS1014 {
		m = "10.14"
	}
	r := "idle"
	if running == k8s {
		r = "k8s"
	}
	return m + "-" + r
}
