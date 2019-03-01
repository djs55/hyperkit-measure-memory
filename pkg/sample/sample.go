package sample

import (
	"time"

	"github.com/djs55/hyperkit-measure-memory/pkg/mem"
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
