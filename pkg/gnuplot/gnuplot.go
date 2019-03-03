package gnuplot

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Graph is a simple time-series gnuplot graph
type Graph struct {
	Title string
	Lines []*Line
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

const kib = int64(1024)
const mib = int64(1024) * kib
const gib = int64(1024) * mib

// Render renders a graph to a .png
func (g *Graph) Render(pngPath string) error {
	dir, err := ioutil.TempDir("", "gnuplot")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	if err := writeGp(g, dir, pngPath); err != nil {
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
		if err := writeDat(*line, dir); err != nil {
			return errors.Wrapf(err, "while plotting %s", line.Label)
		}
	}
	return nil
}

func writeDat(l Line, dir string) error {
	dat, err := os.Create(datPath(l, dir))
	if err != nil {
		return err
	}
	defer dat.Close()
	if _, err := fmt.Fprintf(dat, "# %s\n", l.Label); err != nil {
		return err
	}
	for _, point := range l.Points {
		if _, err := fmt.Fprintf(dat, "%f %f\n", point.Second, float64(point.Memory)/float64(gib)); err != nil {
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
		plots = append(plots, fmt.Sprintf("'%s' using 1:2 with points title '%s'", filepath.Base(datPath(*line, dir)), line.Label))
	}
	path := pngPath
	if !filepath.IsAbs(pngPath) {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		path = filepath.Join(cwd, pngPath)
	}
	lines := []string{
		fmt.Sprintf("set terminal png"),
		fmt.Sprintf("set output '%s'", path),
		fmt.Sprintf("set title '%s'", g.Title),
		fmt.Sprintf("set xlabel 'Time/s'"),
		fmt.Sprintf("set ylabel 'Memory/GiB'"),
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
