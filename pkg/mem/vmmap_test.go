package mem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const kib = int64(1024)
const mib = int64(1024) * kib
const gib = int64(1024) * mib

var exampleVMMapSummary = `
Process:         com.docker.hyperkit [68841]
Path:            /Users/djs/go/src/github.com/docker/pinata/mac/build/Docker.app/Contents/Resources/bin/com.docker.hyperkit
Load Address:    0x10d0d4000
Identifier:      com.docker.hyperkit
Version:         ???
Code Type:       X86-64
Parent Process:  com.docker.driver.amd64-linux [68833]

Date/Time:       2019-02-28 14:22:55.943 +0000
Launch Time:     2019-02-26 12:31:24.767 +0000
OS Version:      Mac OS X 10.14.1 (18B75)
Report Version:  7
Analysis Tool:   /Applications/Xcode.app/Contents/Developer/usr/bin/vmmap
Analysis Tool Version:  Xcode 10.1 (10B61)

Physical footprint:         3.7G
Physical footprint (peak):  3.7G
----

ReadOnly portion of Libraries: Total=333.3M resident=143.8M(43%) swapped_out_or_unallocated=189.5M(57%)
Writable regions: Total=2.2G written=2880K(0%) resident=2.1G(97%) swapped_out=96K(0%) unallocated=60.2M(3%)

                                VIRTUAL RESIDENT    DIRTY  SWAPPED VOLATILE   NONVOL    EMPTY   REGION
REGION TYPE                        SIZE     SIZE     SIZE     SIZE     SIZE     SIZE     SIZE    COUNT (non-coalesced)
===========                     ======= ========    =====  ======= ========   ======    =====  =======
Activity Tracing                   256K      20K      20K       0K       0K      20K       0K        2
Dispatch continuations            8192K       8K       8K       0K       0K       0K       0K        2
IOKit                              112K     112K     112K       0K       0K       0K       0K        3
Kernel Alloc Once                    8K       4K       4K       0K       0K       0K       0K        2
MALLOC guard page                   16K       0K       0K       0K       0K       0K       0K        5
MALLOC metadata                     44K      44K      44K       0K       0K       0K       0K        6
MALLOC_LARGE                       2.1G     2.1G     2.1G       0K       0K       0K       0K        8         see MALLOC ZONE table below
MALLOC_LARGE metadata                4K       4K       4K       0K       0K       0K       0K        2         see MALLOC ZONE table below
MALLOC_SMALL                      32.0M     480K     476K      36K       0K       0K       0K        2         see MALLOC ZONE table below
MALLOC_TINY                       4096K     136K     136K      16K       0K       0K       0K        4         see MALLOC ZONE table below
STACK GUARD                       56.1M       0K       0K       0K       0K       0K       0K       18
Stack                             16.1M     152K     152K      44K       0K       0K       0K       18
__DATA                            19.2M    13.9M    3744K       0K       0K       0K       0K      189
__FONT_DATA                          4K       0K       0K       0K       0K       0K       0K        2
__LINKEDIT                       217.5M    76.8M       0K       0K       0K       0K       0K        4
__TEXT                           115.7M    67.0M       0K       0K       0K       0K       0K      190
__UNICODE                          564K     448K       0K       0K       0K       0K       0K        2
__ctl_set                            4K       4K       4K       0K       0K       0K       0K        2
__inout_port_set                     4K       4K       4K       0K       0K       0K       0K        2
__lpc_dsdt_set                       4K       4K       4K       0K       0K       0K       0K        2
__lpc_sysres_set                     4K       4K       4K       0K       0K       0K       0K        2
__pci_devemu_set                     4K       4K       4K       0K       0K       0K       0K        2
shared memory                       12K      12K      12K       0K       0K       0K       0K        4
===========                     ======= ========    =====  ======= ========   ======    =====  =======
TOTAL                              2.6G     2.3G     2.1G      96K       0K      20K       0K      450

                                 VIRTUAL   RESIDENT      DIRTY    SWAPPED ALLOCATION      BYTES DIRTY+SWAP          REGION
MALLOC ZONE                         SIZE       SIZE       SIZE       SIZE      COUNT  ALLOCATED  FRAG SIZE  % FRAG   COUNT
===========                      =======  =========  =========  =========  =========  =========  =========  ======  ======
DefaultMallocZone_0x10d87a000       2.2G       2.1G       2.1G        52K       1210       2.1G         0K      0%      12
`

func TestVMMapSummary(t *testing.T) {
	footprint, err := parseVMMapSummary(exampleVMMapSummary)
	assert.Nil(t, err)
	// 3.7 GiB rounded down
	assert.Equal(t, Footprint(3972844748), footprint)
}

var exampleVMMapSummaryFirefox1014 = `
Process:         firefox [7675]
Path:            /Applications/Firefox.app/Contents/MacOS/firefox
Load Address:    0x101dd8000
Identifier:      org.mozilla.firefox
Version:         65.0.2 (6519.2.25)
Code Type:       X86-64
Parent Process:  ??? [1]

Date/Time:       2019-03-04 10:37:10.344 +0000
Launch Time:     2019-03-01 17:11:35.641 +0000
OS Version:      Mac OS X 10.14.1 (18B75)
Report Version:  7
Analysis Tool:   /Applications/Xcode.app/Contents/Developer/usr/bin/vmmap
Analysis Tool Version:  Xcode 10.1 (10B61)

Physical footprint:         307.2M
Physical footprint (peak):  323.7M
`

func TestVMMapSummaryFirefox(t *testing.T) {
	footprint, err := parseVMMapSummary(exampleVMMapSummaryFirefox1014)
	assert.Nil(t, err)
	// 307.2M rounded down
	assert.Equal(t, Footprint(322122547), footprint)
}

var exampleVMMapSummary1012 = `
Process:         com.docker.hyperkit [1387]
Path:            /Applications/Docker.app/Contents/Resources/bin/com.docker.hyperkit
Load Address:    0x105b55000
Identifier:      com.docker.hyperkit
Version:         ???
Code Type:       X86-64
Parent Process:  com.docker.driver.amd64-linux [1381]

Date/Time:       2019-03-01 11:11:31.829 +0000
Launch Time:     2019-03-01 11:08:41.499 +0000
OS Version:      Mac OS X 10.12 (16A323)
Report Version:  7
Analysis Tool:   /usr/bin/vmmap
----

ReadOnly portion of Libraries: Total=194.3M resident=0K(0%) swapped_out_or_unallocated=194.3M(100%)
Writable regions: Total=2.2G written=0K(0%) resident=0K(0%) swapped_out=0K(0%) unallocated=2.2G(100%)

                                VIRTUAL   REGION
REGION TYPE                        SIZE    COUNT (non-coalesced)
===========                     =======  =======
Activity Tracing                   256K        2
Dispatch continuations            8192K        2
IOKit                              112K        3
Kernel Alloc Once                    8K        2
MALLOC guard page                   16K        4
MALLOC metadata                    180K        6
MALLOC_LARGE                       2.2G        8         see MALLOC ZONE table below
MALLOC_LARGE metadata                8K        2         see MALLOC ZONE table below
MALLOC_SMALL                      24.0M        2         see MALLOC ZONE table below
MALLOC_SMALL (empty)              8192K        2         see MALLOC ZONE table below
MALLOC_TINY                       4096K        2         see MALLOC ZONE table below
STACK GUARD                       56.1M       18
Stack                             16.1M       18
__DATA                            12.1M      139
__LINKEDIT                       112.7M        6
__TEXT                            81.5M      142
__UNICODE                          556K        2
__ctl_set                            4K        2
__inout_port_set                     4K        2
__lpc_dsdt_set                       4K        2
__lpc_sysres_set                     4K        2
__pci_devemu_set                     4K        2
shared memory                       12K        4
===========                     =======  =======
TOTAL                              2.5G      351

                                 VIRTUAL ALLOCATION      BYTES          REGION
MALLOC ZONE                         SIZE      COUNT  ALLOCATED  % FULL   COUNT
===========                      =======  =========  =========  ======  ======
DefaultMallocZone_0x106233000       2.2G       1620       2.2G     98%      11
`

func TestVMMapSummary1012(t *testing.T) {
	_, err := parseVMMapSummary(exampleVMMapSummary1012)
	assert.Equal(t, ErrNoPhysicalFootprint, err)
}
