package core

import (
	"fmt"
	"testing"
)

func TestNewAoiManager(t *testing.T) {
	aoiMgr := NewAOIManager(100, 300, 4, 200, 450, 5)
	fmt.Println(aoiMgr.String())
}
