package mem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var exampleMeminfo = `
MemTotal:        2047036 kB
MemFree:          422564 kB
MemAvailable:    1579656 kB
Buffers:          113716 kB
Cached:          1089172 kB
SwapCached:          220 kB
Active:           305780 kB
Inactive:        1115336 kB
Active(anon):      91168 kB
Inactive(anon):   128100 kB
Active(file):     214612 kB
Inactive(file):   987236 kB
Unevictable:           0 kB
Mlocked:               0 kB
SwapTotal:       1048572 kB
SwapFree:        1044676 kB
Dirty:               228 kB
Writeback:             0 kB
AnonPages:        214228 kB
Mapped:            86876 kB
Shmem:              1000 kB
Slab:             175680 kB
SReclaimable:     143360 kB
SUnreclaim:        32320 kB
KernelStack:        7232 kB
PageTables:         2152 kB
NFS_Unstable:          0 kB
Bounce:                0 kB
WritebackTmp:          0 kB
CommitLimit:     2072088 kB
Committed_AS:     672304 kB
VmallocTotal:   34359738367 kB
VmallocUsed:           0 kB
VmallocChunk:          0 kB
AnonHugePages:         0 kB
ShmemHugePages:        0 kB
ShmemPmdMapped:        0 kB
HugePages_Total:       1
HugePages_Free:        0
HugePages_Rsvd:        0
HugePages_Surp:        0
Hugepagesize:       2048 kB
DirectMap4k:       73132 kB
DirectMap2M:     2023424 kB
DirectMap1G:           0 kB
`

func TestMeminfoFirstEntry(t *testing.T) {
	meminfo := parseMeminfo(exampleMeminfo)
	assert.Equal(t, int64(2047036*1024), meminfo["MemTotal"])
}

func TestMeminfoLastEntry(t *testing.T) {
	meminfo := parseMeminfo(exampleMeminfo)
	assert.Equal(t, int64(0), meminfo["DirectMap1G"])
}

func TestMeminfoPages(t *testing.T) {
	meminfo := parseMeminfo(exampleMeminfo)
	assert.Equal(t, int64(4096), meminfo["HugePages_Total"])
}

func TestMeminfoKb(t *testing.T) {
	meminfo := parseMeminfo(exampleMeminfo)
	assert.Equal(t, int64(2023424*1024), meminfo["DirectMap2M"])
}
