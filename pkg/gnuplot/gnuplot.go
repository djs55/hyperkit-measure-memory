package gnuplot

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Graph is a simple time-series gnuplot graph
type Graph struct {
	Title  string
	Lines  []*Line
	Format Format
	Time   Time
}

// Line represents the evolution of some labelled parameter over time
type Line struct {
	Points []*Point
	Label  string
}

// Point represents a single sample of a parameter
type Point struct {
	Second float64 // Time the sample was taken
	Memory int64   // Memory value
}

// Format of the rendered output
type Format int

const (
	PNG = Format(0)
	SVG = Format(1)
)

// Time unit to use on the x axis
type Time int

const (
	Seconds = Time(0)
	Hours   = Time(1)
)

const kib = int64(1024)
const mib = int64(1024) * kib
const gib = int64(1024) * mib

// Render renders a graph
func (g *Graph) Render(path string) error {
	dir, err := ioutil.TempDir("", "gnuplot")
	if err != nil {
		return err
	}
	//defer os.RemoveAll(dir)
	fmt.Printf("%s\n", dir)

	if err := writeGp(g, dir, path); err != nil {
		return err
	}
	if err := writeDats(g, dir); err != nil {
		return err
	}

	cmd := exec.Command("gnuplot", gpPath(dir))
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "gnuplot: %s", out)
	}
	return nil
}

func writeDats(g *Graph, dir string) error {
	if len(g.Lines) == 0 {
		return errors.New("There must be at least one line to plot")
	}
	for _, line := range g.Lines {
		if err := writeDat(g, *line, dir); err != nil {
			return errors.Wrapf(err, "while plotting %s", line.Label)
		}
	}
	return nil
}

func writeDat(g *Graph, l Line, dir string) error {
	dat, err := os.Create(datPath(l, dir))
	if err != nil {
		return err
	}
	defer dat.Close()
	if _, err := fmt.Fprintf(dat, "# %s\n", l.Label); err != nil {
		return err
	}
	for _, point := range l.Points {
		time := point.Second
		if g.Time == Hours {
			time = time / 3600.0
		}
		if _, err := fmt.Fprintf(dat, "%f %f\n", time, float64(point.Memory)/float64(gib)); err != nil {
			return err
		}
	}
	return nil
}

func datPath(l Line, dir string) string {
	return filepath.Join(dir, fmt.Sprintf("%s.dat", l.Label))
}

func gpPath(dir string) string {
	return filepath.Join(dir, "graph.gp")
}

func writeGp(g *Graph, dir, pngPath string) error {
	gp, err := os.Create(gpPath(dir))
	if err != nil {
		return err
	}
	defer gp.Close()
	var plots []string
	for _, line := range g.Lines {
		plots = append(plots, fmt.Sprintf("'%s' using 1:2 with lines title '%s'", filepath.Base(datPath(*line, dir)), line.Label))
	}
	path := pngPath
	if !filepath.IsAbs(pngPath) {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		path = filepath.Join(cwd, pngPath)
	}
	// round up to the next GiB
	maxMem := (maxMemoryValue(g) + gib - int64(1)) / gib
	var terminal string
	switch g.Format {
	case SVG:
		terminal = "svg"
	default:
		terminal = "png"
	}
	var xlabel string
	switch g.Time {
	case Hours:
		xlabel = "Time/hours"
	default:
		xlabel = "Time/seconds"
	}
	lines := []string{
		fmt.Sprintf("set terminal " + terminal),
		fmt.Sprintf("set output '%s'", path),
		fmt.Sprintf("set title '%s'", g.Title),
		fmt.Sprintf("set xlabel '%s'", xlabel),
		fmt.Sprintf("set ylabel 'Memory/GiB'"),
		fmt.Sprintf("set yrange [0:%d]", maxMem+1),
		//"set timefmt '%s'",
		//"set xdata time",
		fmt.Sprintf("plot %s", strings.Join(plots, ", ")),
	}
	for _, line := range lines {
		if _, err := fmt.Fprintf(gp, "%s\n", line); err != nil {
			return err
		}
	}
	return nil
}

func maxMemoryValue(g *Graph) int64 {
	m := int64(0)
	for _, line := range g.Lines {
		for _, point := range line.Points {
			if point.Memory > m {
				m = point.Memory
			}
		}
	}
	return m
}

// ReadPoints reads a gnuplot-format datafile
func ReadPoints(file string) ([]*Point, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b := bufio.NewReader(f)
	var results []*Point
	for {
		line, err := b.ReadString(byte('\n'))
		if err == io.EOF {
			return results, nil
		}
		// trim comment
		line = strings.Split(line, "#")[0]
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			log.Printf("Skipping line with len <> 2: %s", line)
			continue
		}
		second, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			log.Printf("Skipping non-float: %s", line)
			continue
		}
		memory, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			log.Printf("Skipping non-int64: %s", line)
			continue
		}
		results = append(results, &Point{Second: second, Memory: memory})
	}
}
