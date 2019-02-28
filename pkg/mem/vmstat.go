package mem

import (
	"os/exec"
	"strconv"
	"strings"
)

// VMStat of VM memory stats
type VMStat map[string]int64

func parseVMStat(txt string) VMStat {
	// assume page size is 4096
	pageSize := int64(4096)
	snapshot := make(VMStat)
	for _, line := range strings.Split(txt, "\n") {
		if strings.HasPrefix(line, "Mach Virtual Memory Statistics:") {
			continue
		}
		bits := strings.SplitN(line, ":", 2)
		if len(bits) != 2 {
			continue
		}
		key := bits[0]
		pages := bits[1]
		if strings.HasSuffix(pages, ".") {
			pages = pages[0 : len(pages)-1]
		}
		pages = strings.TrimSpace(pages)
		i, err := strconv.ParseInt(pages, 10, 64)
		if err != nil {
			continue
		}
		snapshot[key] = i * pageSize
	}
	return snapshot
}

// Get a memory snapshot.
func GetVMStat() (VMStat, error) {
	cmd := exec.Command("vm_stat")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return parseVMStat(string(out)), nil
}
