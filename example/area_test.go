package example

import (
	pkg "github.com/gregmus2/ants-pkg"
	"testing"
)

func TestAreaClosest(t *testing.T) {
	area := NewArea(100, 100)

	pos := area.Closest(&pkg.Pos{X: 50, Y: 50}, pkg.EnemyField)
	if pos != nil {
		t.Errorf("Wrong calculation of closest: %v, nil", pos)
	}

	area.matrix[40][40] = pkg.EnemyField
	pos = area.Closest(&pkg.Pos{X: 50, Y: 50}, pkg.EnemyField)
	if pos.X != 40 || pos.Y != 40 {
		t.Errorf("Wrong calculation of closest: %v, &{%d %d}", pos, 40, 40)
	}

	area.matrix[41][40] = pkg.EnemyField
	pos = area.Closest(&pkg.Pos{X: 50, Y: 50}, pkg.EnemyField)
	if pos.X != 41 || pos.Y != 40 {
		t.Errorf("Wrong calculation of closest: %v, &{%d %d}", pos, 41, 40)
	}

	area.matrix[51][50] = pkg.EnemyField
	pos = area.Closest(&pkg.Pos{X: 50, Y: 50}, pkg.EnemyField)
	if pos.X != 51 || pos.Y != 50 {
		t.Errorf("Wrong calculation of closest: %v, &{%d %d}", pos, 51, 50)
	}

	pos = area.Closest(&pkg.Pos{X: 1, Y: 5}, pkg.EnemyField)
	if pos.X != 40 || pos.Y != 40 {
		t.Errorf("Wrong calculation of closest: %v, &{%d %d}", pos, 40, 40)
	}

	pos = area.Closest(&pkg.Pos{X: 5, Y: 1}, pkg.EnemyField)
	if pos.X != 41 || pos.Y != 40 {
		t.Errorf("Wrong calculation of closest: %v, &{%d %d}", pos, 41, 40)
	}

	pos = area.Closest(&pkg.Pos{X: 98, Y: 95}, pkg.EnemyField)
	if pos.X != 51 || pos.Y != 50 {
		t.Errorf("Wrong calculation of closest: %v, &{%d %d}", pos, 51, 50)
	}

	pos = area.Closest(&pkg.Pos{X: 98, Y: 50}, pkg.EnemyField)
	if pos.X != 51 || pos.Y != 50 {
		t.Errorf("Wrong calculation of closest: %v, &{%d %d}", pos, 51, 50)
	}

	pos = area.Closest(&pkg.Pos{X: 32, Y: 3}, pkg.EnemyField)
	if pos.X != 41 || pos.Y != 40 {
		t.Errorf("Wrong calculation of closest: %v, &{%d %d}", pos, 41, 40)
	}
}

func TestAreaRewriteMap(t *testing.T) {
	area := NewArea(100, 100)

	area.matrix[45][39] = pkg.AllyField
	ok := area.RewriteMap(0, 0, pkg.EnemyField)
	if ok || area.matrix[45][39] != pkg.AllyField {
		t.Error("Area Expansion False")
	}

	ok = area.RewriteMap(99, 99, pkg.AnthillField)
	if ok || area.matrix[45][39] != pkg.AllyField {
		t.Error("Area Expansion False")
	}

	ok = area.RewriteMap(-5, 2, pkg.AnthillField)
	if !ok || area.matrix[145][39] != pkg.AllyField || area.w != 200 {
		t.Error("Area Expansion False")
	}

	ok = area.RewriteMap(58, -5, pkg.EnemyField)
	if !ok || area.matrix[145][139] != pkg.AllyField || area.h != 200 {
		t.Error("Area Expansion False")
	}

	ok = area.RewriteMap(220, -5, pkg.EnemyField)
	if !ok || area.matrix[145][339] != pkg.AllyField || area.w != 400 || area.h != 400 {
		t.Error("Area Expansion False")
	}

	ok = area.RewriteMap(220, 405, pkg.EnemyField)
	if !ok || area.matrix[145][339] != pkg.AllyField || area.w != 400 || area.h != 800 {
		t.Errorf("Area Expansion False. target: %v, area.w: %d, area.h: %d", area.matrix[145][339], area.w, area.h)
	}
}
