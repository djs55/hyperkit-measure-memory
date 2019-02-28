package mem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var examplePs = `
RSS      VSZ COMM
1761620 111487192 com.apple.WebKit.WebContent
1251196  6559436 com.docker.hyperkit
650480  7720696 Safari
579532 110215776 com.apple.WebKit.WebContent
552260  8903848 Console
513400 91336008 com.apple.WebKit.WebContent
496644 107040856 com.apple.WebKit.WebContent
314200  6155928 Slack Helper
248628 108020528 com.apple.WebKit.WebContent
237084  5776604 Slack Helper
222668 107426936 com.apple.WebKit.WebContent
220640  6075504 Code Helper
175588  4870400 Slack Helper
152852  5122108 Spotify Helper
138572  5456020 iTerm2
124056  4691436 AppleSpell
106644  5938488 Slack
 98908 91176488 com.apple.WebKit.WebContent
 700  4270648 sh
 420  4280576 periodic-wrapper
 420  4288768 periodic-wrapper
 420  4296960 periodic-wrapper
   0        0 (uname)
`

func TestRSS(t *testing.T) {
	ps, err := parsePS(examplePs)
	assert.Nil(t, err)
	assert.Equal(t, 23, len(ps))
	assert.Equal(t, int64(1761620*1024), ps[0].RSS)
	assert.Equal(t, int64(111487192*1024), ps[0].VSZ)
	assert.Equal(t, "com.apple.WebKit.WebContent", ps[0].Command)
}
