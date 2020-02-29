package main

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

	area.matrix[51][50] = pkg.EnemyField
	pos = area.Closest(&pkg.Pos{X: 50, Y: 50}, pkg.EnemyField)
	if pos.X != 51 || pos.Y != 50 {
		t.Errorf("Wrong calculation of closest: %v, &{%d %d}", pos, 51, 50)
	}

	pos = area.Closest(&pkg.Pos{X: 1, Y: 5}, pkg.EnemyField)
	if pos.X != 40 || pos.Y != 40 {
		t.Errorf("Wrong calculation of closest: %v, &{%d %d}", pos, 40, 40)
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
	if (pos.X != 41 || pos.Y != 40) && (pos.X != 40 || pos.Y != 40) {
		t.Errorf("Wrong calculation of closest: %v, &{%d %d}", pos, 41, 40)
	}
}

func TestAreaRewriteMap(t *testing.T) {
	ai := NewAI(&pkg.Pos{X: 1}, 1)
	ai.ants[1] = &Ant{
		Pos: &pkg.Pos{X: 51, Y: 50},
	}
	ai.enemyAnthills = append(ai.enemyAnthills, &pkg.Pos{X: 23, Y: 21})

	ai.area.matrix[45][39] = pkg.AllyField
	ok := ai.area.RewriteMap(0, 0, &ai)
	if ok || ai.area.matrix[45][39] != pkg.AllyField ||
		ai.ants[1].Pos.X != 51 || ai.ants[1].Pos.Y != 50 ||
		ai.anthills[1].Pos.X != 50 || ai.anthills[1].Pos.Y != 50 ||
		ai.enemyAnthills[0].X != 23 || ai.enemyAnthills[0].Y != 21 {
		t.Error("Area Expansion False")
	}

	ok = ai.area.RewriteMap(99, 99, &ai)
	if ok || ai.area.matrix[45][39] != pkg.AllyField ||
		ai.ants[1].Pos.X != 51 || ai.ants[1].Pos.Y != 50 ||
		ai.anthills[1].Pos.X != 50 || ai.anthills[1].Pos.Y != 50 ||
		ai.enemyAnthills[0].X != 23 || ai.enemyAnthills[0].Y != 21 {
		t.Error("Area Expansion False")
	}

	ok = ai.area.RewriteMap(-5, 2, &ai)
	if !ok || ai.area.matrix[145][39] != pkg.AllyField || ai.area.w != 200 ||
		ai.ants[1].Pos.X != 151 || ai.ants[1].Pos.Y != 50 ||
		ai.anthills[1].Pos.X != 150 || ai.anthills[1].Pos.Y != 50 ||
		ai.enemyAnthills[0].X != 123 || ai.enemyAnthills[0].Y != 21 {
		t.Error("Area Expansion False")
	}

	ok = ai.area.RewriteMap(58, -5, &ai)
	if !ok || ai.area.matrix[145][139] != pkg.AllyField || ai.area.h != 200 ||
		ai.ants[1].Pos.X != 151 || ai.ants[1].Pos.Y != 150 ||
		ai.anthills[1].Pos.X != 150 || ai.anthills[1].Pos.Y != 150 ||
		ai.enemyAnthills[0].X != 123 || ai.enemyAnthills[0].Y != 121 {
		t.Error("Area Expansion False")
	}

	ok = ai.area.RewriteMap(220, -5, &ai)
	if !ok || ai.area.matrix[145][339] != pkg.AllyField || ai.area.w != 400 || ai.area.h != 400 ||
		ai.ants[1].Pos.X != 151 || ai.ants[1].Pos.Y != 350 ||
		ai.anthills[1].Pos.X != 150 || ai.anthills[1].Pos.Y != 350 ||
		ai.enemyAnthills[0].X != 123 || ai.enemyAnthills[0].Y != 321 {
		t.Error("Area Expansion False")
	}

	ok = ai.area.RewriteMap(220, 405, &ai)
	if !ok || ai.area.matrix[145][339] != pkg.AllyField || ai.area.w != 400 || ai.area.h != 800 ||
		ai.ants[1].Pos.X != 151 || ai.ants[1].Pos.Y != 350 ||
		ai.anthills[1].Pos.X != 150 || ai.anthills[1].Pos.Y != 350 ||
		ai.enemyAnthills[0].X != 123 || ai.enemyAnthills[0].Y != 321 {
		t.Errorf("Area Expansion False. target: %v, ai.area.w: %d, ai.area.h: %d", ai.area.matrix[145][339], ai.area.w, ai.area.h)
	}
}

func TestArea_CutArea(t *testing.T) {
	ai := NewAI(&pkg.Pos{X: 1}, 1)
	ai.ants[1] = &Ant{
		Pos: &pkg.Pos{X: 51, Y: 50},
	}
	ai.enemyAnthills = append(ai.enemyAnthills, &pkg.Pos{X: 23, Y: 21})

	fields := [5][5]pkg.FieldType{
		{pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField},
		{pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField},
		{pkg.FoodField, pkg.EmptyField, pkg.AllyField, pkg.EmptyField, pkg.EmptyField},
		{pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField},
		{pkg.WallField, pkg.WallField, pkg.WallField, pkg.WallField, pkg.WallField},
	}
	ai.area.CutArea(fields, &pkg.Pos{X: 65, Y: 58}, &ai)

	if ai.area.w != 69 || ai.area.h != 100 ||
		ai.ants[1].Pos.X != 51 || ai.ants[1].Pos.Y != 50 ||
		ai.enemyAnthills[0].X != 23 || ai.enemyAnthills[0].Y != 21 {
		t.Errorf(
			"Area Expansion False. w: %d/%d, h: %d/%d, ant: %v/%v",
			ai.area.w, 69,
			ai.area.h, 100,
			ai.ants[1].Pos, pkg.Pos{51, 50},
		)
	}

	fields = [5][5]pkg.FieldType{
		{pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.WallField},
		{pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.WallField},
		{pkg.FoodField, pkg.EmptyField, pkg.AllyField, pkg.EmptyField, pkg.WallField},
		{pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.WallField},
		{pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.WallField},
	}
	ai.area.CutArea(fields, &pkg.Pos{X: 65, Y: 58}, &ai)

	if ai.area.w != 69 || ai.area.h != 62 ||
		ai.ants[1].Pos.X != 51 || ai.ants[1].Pos.Y != 50 ||
		ai.enemyAnthills[0].X != 23 || ai.enemyAnthills[0].Y != 21 {
		t.Errorf(
			"Area Expansion False. w: %d/%d, h: %d/%d, ant: %v/%v",
			ai.area.w, 69,
			ai.area.h, 62,
			ai.ants[1].Pos, pkg.Pos{51, 50},
		)
	}

	fields = [5][5]pkg.FieldType{
		{pkg.WallField, pkg.WallField, pkg.WallField, pkg.WallField, pkg.WallField},
		{pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField},
		{pkg.FoodField, pkg.EmptyField, pkg.AllyField, pkg.EmptyField, pkg.EmptyField},
		{pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField},
		{pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField},
	}
	ai.area.CutArea(fields, &pkg.Pos{X: 20, Y: 20}, &ai)

	if ai.area.w != 52 || ai.area.h != 62 ||
		ai.ants[1].Pos.X != 34 || ai.ants[1].Pos.Y != 50 ||
		ai.enemyAnthills[0].X != 6 || ai.enemyAnthills[0].Y != 21 {
		t.Errorf(
			"Area Expansion False. w: %d/%d, h: %d/%d, ant: %v/%v",
			ai.area.w, 52,
			ai.area.h, 62,
			ai.ants[1].Pos, pkg.Pos{34, 50},
		)
	}
}
