package example

import (
	pkg "github.com/gregmus2/ants-pkg"
	"testing"
)

func TestAnt_CalcStep(t *testing.T) {
	ant := &Ant{Pos: &pkg.Pos{X: 55, Y: 56}}
	pos := ant.CalcStep(60, 60)
	if pos.X != 1 || pos.Y != 1 {
		t.Errorf("Wrong calculation of step: %v, &{%d %d}", pos, 1, 1)
	}

	pos = ant.CalcStep(60, 56)
	if pos.X != 1 || pos.Y != 0 {
		t.Errorf("Wrong calculation of step: %v, &{%d %d}", pos, 1, 0)
	}

	pos = ant.CalcStep(54, 20)
	if pos.X != -1 || pos.Y != -1 {
		t.Errorf("Wrong calculation of step: %v, &{%d %d}", pos, -1, -1)
	}
}

func TestAnt_RelativePos(t *testing.T) {
	ant := &Ant{Pos: &pkg.Pos{X: 55, Y: 56}}
	pos := ant.RelativePos(57, 58)
	if pos.X != 4 || pos.Y != 4 {
		t.Errorf("Wrong calculation of step: %v, &{%d %d}", pos, 4, 4)
	}

	pos = ant.RelativePos(56, 56)
	if pos.X != 3 || pos.Y != 2 {
		t.Errorf("Wrong calculation of step: %v, &{%d %d}", pos, 3, 2)
	}

	pos = ant.RelativePos(54, 54)
	if pos.X != 1 || pos.Y != 0 {
		t.Errorf("Wrong calculation of step: %v, &{%d %d}", pos, 1, 0)
	}
}
