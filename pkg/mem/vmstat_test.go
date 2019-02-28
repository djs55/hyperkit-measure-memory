package mem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var exampleVMStat = `
Pages free:                               26027.
Pages active:                           1404521.
Pages inactive:                         1391220.
Pages speculative:                         8996.
Pages throttled:                              0.
Pages wired down:                        755063.
Pages purgeable:                          37920.
"Translation faults":                1347146608.
Pages copy-on-write:                   58342606.
Pages zero filled:                    333813383.
Pages reactivated:                     79160241.
Pages purged:                           3145959.
File-backed pages:                       901294.
Anonymous pages:                        1903443.
Pages stored in compressor:             4023450.
Pages occupied by compressor:            608294.
Decompressions:                        30571821.
Compressions:                          43062520.
Pageins:                               95889608.
Pageouts:                                447048.
Swapins:                               28904306.
Swapouts:                              30414707.
`

func TestVmstatFirst(t *testing.T) {
	vmstat := parseVMStat(exampleVMStat)
	assert.Equal(t, int64(26027*4096), vmstat["Pages free"])
}

func TestVmstatLast(t *testing.T) {
	vmstat := parseVMStat(exampleVMStat)
	assert.Equal(t, int64(30414707*4096), vmstat["Swapouts"])
}
