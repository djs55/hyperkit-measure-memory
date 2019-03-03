package sample

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/djs55/hyperkit-measure-memory/pkg/gnuplot"
	"github.com/djs55/hyperkit-measure-memory/pkg/mem"
	"github.com/pkg/errors"
)

// Sample is a set of memory measurements made at the same time
type Sample struct {
	Time             time.Time
	Meminfo          mem.Meminfo
	PS               []mem.PS
	Footprint        mem.Footprint
	FirefoxFootprint mem.Footprint
	VMStat           mem.VMStat
}

// ReadDir reads a directory full of samples and converts into an array of points.
func ReadDir(dir string, f func(s Sample) int64) ([]*gnuplot.Point, error) {
	var points []*gnuplot.Point
	for count := 0; ; count++ {
		path := filepath.Join(dir, fmt.Sprintf("%d", count))
		input, err := os.Open(path)
		if os.IsNotExist(err) {
			break
		}
		if err != nil {
			return nil, errors.Wrapf(err, "failed to open %s", path)
		}
		dec := json.NewDecoder(input)
		var s Sample
		if err := dec.Decode(&s); err != nil {
			log.Printf("Failed to decode %s: %v", path, err)
			continue
		}
		points = append(points, &gnuplot.Point{
			Second: float64(s.Time.Unix()),
			Memory: f(s),
		})

		if err := input.Close(); err != nil {
			return nil, errors.Wrapf(err, "Failed to close %s", path)
		}
	}
	return points, nil
}
