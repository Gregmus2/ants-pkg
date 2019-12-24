package main

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

	fields = [5][5]pkg.FieldType{
		{pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField},
		{pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField},
		{pkg.EmptyField, pkg.EmptyField, pkg.AllyField, pkg.EmptyField, pkg.EmptyField},
		{pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField},
		{pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField, pkg.EmptyField},
	}

	target = &pkg.Pos{}
	for i := 0; i < 10; i++ {
		target, action = Greg.Do(1, fields, 1, target)
		if target.X == 0 && target.Y == 0 {
			t.Error(target)
		}
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

func TestAI_OnNewAnthill(t *testing.T) {
	ai := NewAI(&pkg.Pos{1, 0}, 1)
	ai.area.matrix[40][40] = pkg.EnemyAnthillField
	ai.OnAntBirth(23, 1)
	ai.OnNewAnthill(23, &pkg.Pos{1, 1}, 2)
	if ai.anthills[2].Pos.X != 40 || ai.anthills[2].Pos.Y != 40 {
		t.Errorf("Wrong pos of new anthill. Expected: %d %d. Actual: %v", 40, 40, ai.ants[2].Pos)
	}
}

func TestAI_OnAnthillDie(t *testing.T) {
	ai := NewAI(&pkg.Pos{1, 0}, 1)
	ai.OnAnthillDie(1)
	if _, ok := ai.anthills[1]; ok {
		t.Error("Bad reaction on AnthillDie event")
	}
}
