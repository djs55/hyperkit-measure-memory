package mem

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// PS is the view of a process from `ps`
type PS struct {
	Command string // executable name
	RSS     int64  // Resident size
	VSZ     int64  // Virtual memory size
}

func parsePS(txt string) ([]PS, error) {
	all := make([]PS, 0)
	for _, line := range strings.Split(txt, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "RSS") {
			continue
		}
		if line == "" {
			continue
		}
		var bits []string
		for _, bit := range strings.Split(line, " ") {
			bit = strings.TrimSpace(bit)
			if bit == "" {
				continue
			}
			bits = append(bits, bit)
		}
		if len(bits) < 2 {
			fmt.Printf("Failed to parse %s", line)
			continue
		}
		rss, err := strconv.ParseInt(bits[0], 10, 64)
		if err != nil {
			fmt.Printf("Failed to parse %s", line)
			continue
		}
		vsz, err := strconv.ParseInt(bits[1], 10, 64)
		if err != nil {
			fmt.Printf("Failed to parse VSZ '%s'", bits[1])
			continue
		}
		command := bits[2]
		all = append(all, PS{Command: command, RSS: rss * 1024, VSZ: vsz * 1024})
	}
	return all, nil
}

// GetPS returns stats on all running processes
func GetPS() ([]PS, error) {
	cmd := exec.Command("ps", "-caxm", "-orss,vsz,comm")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return parsePS(string(out))
}
