package sakura

import (
	"strings"
	"testing"
)

func TestSakura(t *testing.T) {
	var s = "1+2; 3*4; -(10+2);(1+2)*3+4*5+6<<2; 3.1416926*1.618; +3*-4; +\"3.1415\"*2.0; let a=10; ^10;"
	Run(strings.NewReader(s))
}
