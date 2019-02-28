package mem

import (
	"io/ioutil"
	"strconv"
	"strings"
)

// Meminfo is the key=value memory data from /proc/meminfo
type Meminfo map[string]int64

func parseMeminfo(txt string) Meminfo {
	meminfo := make(Meminfo)
	for _, line := range strings.Split(txt, "\n") {
		line = strings.TrimSpace(line)
		bits := strings.SplitN(line, ":", 2)
		if len(bits) != 2 {
			continue
		}
		key := bits[0]
		val := strings.TrimSpace(bits[1])
		toBytes := int64(4096) // assume units are pages and pages are 4096 bytes
		if strings.HasSuffix(val, " kB") {
			toBytes = int64(1024)
			val = val[0 : len(val)-3]
		}
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			continue
		}
		meminfo[key] = i * toBytes
	}
	return meminfo
}

// GetMeminfo returns the memory information from /proc/meminfo
func GetMeminfo() (Meminfo, error) {
	out, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return nil, err
	}

	return parseMeminfo(string(out)), nil
}
