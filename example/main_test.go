package example

import (
	pkg "github.com/gregmus2/ants-pkg"
	"testing"
)

func TestAI_OnAntBirth(t *testing.T) {
	ai := NewAI(&pkg.Pos{1, 0}, 1)
	ai.OnAntBirth(1, 1)
	if ai.ants[1].Pos.X != 51 || ai.ants[1].Pos.Y != 50 {
		t.Errorf("Wrong pos of new ant. Expected: %d %d. Actual: %v", 51, 50, ai.ants[1].Pos)
	}
}

func BenchmarkAI_OnAntDie(b *testing.B) {
	ai := NewAI(&pkg.Pos{1, 0}, 1)
}
