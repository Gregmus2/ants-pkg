package example

import (
	pkg "github.com/gregmus2/ants-pkg"
	"testing"
)

func TestAI_Do(t *testing.T) {
	Greg.Start(1, &pkg.Pos{X: 1})
	Greg.OnAntBirth(1, 1)
	fields := [5][5]pkg.FieldType{
		{pkg.FoodField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField},
		{pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField},
		{pkg.EmptyField, pkg.EmptyField, pkg.AllyField, pkg.EmptyField, pkg.EmptyField},
		{pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField},
		{pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField},
	}

	target, action := Greg.Do(1, fields, 1, &pkg.Pos{})
	if target.X != -1 || target.Y != -1 || action != pkg.EatAction {
		t.Error("Wrong behaviour of algorithm")
	}
}

func TestAI_Start(t *testing.T) {
	Greg.Start(1, &pkg.Pos{X: 1})
	anthill := Greg.anthills[1]

	if anthill.Pos.X != defaultSize/2 || anthill.Pos.Y != defaultSize/2 {
		t.Error("Something wrong with anthill pos")
	}

	if anthill.BirthPos.X != defaultSize/2+1 || anthill.BirthPos.Y != defaultSize/2 {
		t.Error("Something wrong with anthill birth pos")
	}

	if Greg.area.matrix[defaultSize/2][defaultSize/2] != pkg.AllyAnthillField ||
		Greg.area.matrix[defaultSize/2+1][defaultSize/2] != pkg.AllyField {
		t.Error("Something wrong with base area fields")
	}
}

func TestAI_OnAntBirth(t *testing.T) {
	ai := NewAI(&pkg.Pos{1, 0}, 1)
	ai.OnAntBirth(1, 1)
	if ai.ants[1].Pos.X != 51 || ai.ants[1].Pos.Y != 50 {
		t.Errorf("Wrong pos of new ant. Expected: %d %d. Actual: %v", 51, 50, ai.ants[1].Pos)
	}
}

func TestAI_OnAntDie(t *testing.T) {
	ai := NewAI(&pkg.Pos{1, 0}, 1)
	ai.OnAntBirth(1, 1)
	ai.OnAntDie(1)
	if _, ok := ai.ants[1]; ok {
		t.Error("Bad reaction on AntDie event")
	}
}
