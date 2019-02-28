package mem

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

// Footprint is the headline figure in Activity Monitor in 10.14, also known
// as the "physical footprint" in the "vmmap" utility.
type Footprint int64

func parseVMMapSummary(txt string) (Footprint, error) {
	for _, line := range strings.Split(txt, "\n") {
		line = strings.TrimSpace(line)
		prefix := "Physical footprint:"
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		line = line[len(prefix):]
		line = strings.TrimSpace(line)
		multiplier := 1
		if strings.HasSuffix(line, "G") {
			multiplier = 1024 * 1024 * 1024
			line = line[0 : len(line)-1]
		}
		if strings.HasSuffix(line, "M") {
			multiplier = 1024 * 1024
			line = line[0 : len(line)-1]
		}
		if strings.HasSuffix(line, "K") {
			multiplier = 1024
			line = line[0 : len(line)-1]
		}
		i, err := strconv.ParseFloat(line, 64)
		if err != nil {
			return Footprint(0), err
		}
		return Footprint(i * float64(multiplier)), nil
	}
	return Footprint(0), errors.New("Failed to read physical footprint")
}

// GetFootprint returns the "physical footprint" of a process. This is shown as the
// headline figure in Activity Monitor.
func GetFootprint(proc string) (Footprint, error) {
	cmd := exec.Command("vmmap", "-summary", proc)
	out, err := cmd.Output()
	if err != nil {
		return Footprint(0), err
	}
	return parseVMMapSummary(string(out))
}
