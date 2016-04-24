package sakura

import (
	"strings"
	"testing"
)

func TestSakura(t *testing.T) {
	var s = "1+2; 3*4;"
	Run(strings.NewReader(s))
}
