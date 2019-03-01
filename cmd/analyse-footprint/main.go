package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/djs55/hyperkit-measure-memory/pkg/gnuplot"
	"github.com/djs55/hyperkit-measure-memory/pkg/sample"
)

func main() {
	footprint := gnuplot.Line{
		Label: "physical footprint",
	}
	g := gnuplot.Graph{
		Title: "hyperkit physical footprint",
		Lines: []*gnuplot.Line{
			&footprint,
		},
	}

	for count := 0; ; count++ {
		path := filepath.Join("results", fmt.Sprintf("%d", count))
		input, err := os.Open(path)
		if os.IsNotExist(err) {
			break
		}
		if err != nil {
			log.Fatalf("Failed to open %s: %v", path, err)
		}
		dec := json.NewDecoder(input)
		var s sample.Sample
		if err := dec.Decode(&s); err != nil {
			log.Printf("Failed to decode %s: %v", path, err)
			continue
		}
		footprint.Points = append(footprint.Points, gnuplot.Point{
			Second: float64(s.Time.Unix()),
			Memory: int64(s.Footprint),
		})

		if err := input.Close(); err != nil {
			log.Fatalf("Failed to close %s: %v", path, err)
		}
	}
	if err := g.Render("/tmp/output.png"); err != nil {
		log.Fatalf("Failed to render: %v", err)
	}
}
