package main

import (
	pkg "github.com/gregmus2/ants-pkg"
	"testing"
)

func setup() (AI, *Ant) {
	ai := NewAI(&pkg.Pos{X: 1}, 1)
	ant := &Ant{
		ID:   1,
		Pos:  &pkg.Pos{X: defaultSize / 2, Y: defaultSize / 2},
		Role: attacker,
	}
	return ai, ant
}

func TestGiveOrder(t *testing.T) {
	ai, ant := setup()
	ai.enemyAnthills = append(ai.enemyAnthills, &pkg.Pos{X: defaultSize - 10, Y: defaultSize / 2})
	pos, action := GiveOrder(ant, &ai)
	if _, ok := ant.Order.(*Attack); !ok {
		t.Error("Wrong order was set")
	}
	if pos.X != 1 || pos.Y != 0 || action != pkg.MoveAction {
		t.Errorf("Wrong calculation of pos or action: %v, %d", pos, action)
	}

	ant.Order = nil
	ant.Role = explorer
	pos, action = GiveOrder(ant, &ai)
	if _, ok := ant.Order.(*Explore); !ok {
		t.Error("Wrong order was set")
	}

	ant.Order = nil
	ant.Role = defender
	pos, action = GiveOrder(ant, &ai)
	if _, ok := ant.Order.(*Defend); !ok {
		t.Error("Wrong order was set")
	}

	ant.Order = nil
	ant.Role = 9
	pos, action = GiveOrder(ant, &ai)
	if _, ok := ant.Order.(*Explore); !ok {
		t.Error("Wrong order was set")
	}
}

func TestGiveOrderWithUrgent(t *testing.T) {
	ai, ant := setup()
	ai.enemyAnthills = append(ai.enemyAnthills, &pkg.Pos{X: defaultSize - 10, Y: defaultSize / 2})
	ai.area.matrix[defaultSize/2-1][defaultSize/2+1] = pkg.FoodField

	pos, action := GiveOrder(ant, &ai)
	if pos.X != -1 || pos.Y != 1 || action != pkg.EatAction {
		t.Errorf("Wrong calculation of pos or action: %v, %d", pos, action)
	}
}

func TestUrgent(t *testing.T) {
	ai, ant := setup()
	ant.Order = GetOrder(ant, &ai)

	pos, action, ok := ant.Order.urgent()
	if pos.X != 0 || pos.Y != 0 || action != 0 || ok == true {
		t.Errorf("Wrong calculation of pos or action: %v, %d, %v", pos, action, ok)
	}

	ai.area.matrix[defaultSize/2-1][defaultSize/2+1] = pkg.FoodField

	pos, action, ok = ant.Order.urgent()
	if pos.X != -1 || pos.Y != 1 || action != pkg.EatAction || ok == false {
		t.Errorf("Wrong calculation of pos or action: %v, %d", pos, action)
	}

	ai.area.matrix[defaultSize/2+1][defaultSize/2] = pkg.EnemyAnthillField
	pos, action, ok = ant.Order.urgent()
	if pos.X != 1 || pos.Y != 0 || action != pkg.AttackAction || ok == false {
		t.Errorf("Wrong calculation of pos or action: %v, %d", pos, action)
	}
}

func TestUrgentEnemyPriority(t *testing.T) {
	ai, ant := setup()

	ai.area.matrix[defaultSize/2-1][defaultSize/2+1] = pkg.FoodField
	ai.area.matrix[defaultSize/2+1][defaultSize/2] = pkg.EnemyField

	pos, action := GiveOrder(ant, &ai)
	if action != pkg.AttackAction {
		t.Errorf("Wrong calculation of action: %d", action)
	}

	if pos.X != 1 || pos.Y != 0 {
		t.Errorf("Wrong calculation of pos or action. Expected: %v, actual: %v", &pkg.Pos{X: 1}, pos)
	}
}

func TestFollow(t *testing.T) {
	ai, ant := setup()
	explore := &Explore{BaseOrder{Ant: ant, AI: &ai}}
	ant.Order = explore

	dataProvider := map[*pkg.Pos]*pkg.Pos{
		{X: defaultSize / 2, Y: defaultSize/2 + 1}:     {X: 0, Y: 1},
		{X: defaultSize/2 + 35, Y: defaultSize/2 + 1}:  {X: 1, Y: 1},
		{X: defaultSize/2 + 35, Y: defaultSize/2 - 20}: {X: 1, Y: -1},
		{X: defaultSize/2 - 35, Y: defaultSize / 2}:    {X: -1, Y: 0},
	}

	for goal, step := range dataProvider {
		explore.Pos = goal
		pos, action := ant.Order.follow()
		if action != pkg.MoveAction {
			t.Errorf("Wrong calculation of action: %d", action)
		}

		if pos.X != step.X || pos.Y != step.Y {
			t.Errorf("Wrong calculation of pos. Expected: %v, actual: %v", step, pos)
		}
	}
}

func TestExploreGoal(t *testing.T) {
	ai, ant := setup()
	explore := &Explore{BaseOrder{Ant: ant, AI: &ai}}
	ant.Order = explore

	for x := range ai.area.matrix {
		for y := range ai.area.matrix[x] {
			ai.area.matrix[x][y] = pkg.EmptyField
		}
	}

	nearestPos := &pkg.Pos{X: defaultSize/2 - 11, Y: defaultSize/2 - 5}
	ai.area.matrix[defaultSize/2-10][defaultSize/2-5] = unknownField
	ai.area.matrix[nearestPos.X][nearestPos.Y] = unknownField

	ant.Order.goal()
	if explore.Pos.X != nearestPos.X && explore.Pos.Y != nearestPos.Y {
		t.Errorf("Wrong calculation of goal. Expected: %v, actual: %v", nearestPos, explore.Pos)
	}

	ai.area.matrix[defaultSize/2+5][defaultSize/2-5] = unknownField

	ant.Order.goal()
	if explore.Pos.X != nearestPos.X && explore.Pos.Y != nearestPos.Y {
		t.Errorf("Wrong calculation of goal. Expected: %v, actual: %v", &pkg.Pos{X: 5, Y: -5}, explore.Pos)
	}
}

func TestDefendGoal(t *testing.T) {
	ai, ant := setup()
	defend := &Defend{BaseOrder: BaseOrder{Ant: ant, AI: &ai}, target: nil}
	ant.Order = defend
	ant.Pos.X = defaultSize/2 + 3
	ant.Pos.Y = defaultSize/2 - 2

	ant.Order.goal()
	if defend.Pos.X != defaultSize/2+2 && defend.Pos.Y != defaultSize/2-2 {
		t.Errorf("Wrong calculation of goal. Expected: %v, actual: %v", &pkg.Pos{X: defaultSize/2 + 3, Y: defaultSize/2 - 2}, defend.Pos)
	}
	if defend.target.X != defaultSize/2 && defend.target.Y != defaultSize/2 {
		t.Errorf("Wrong calculation of target. Expected: %v, actual: %v", &pkg.Pos{X: defaultSize / 2, Y: defaultSize / 2}, defend.target)
	}

	ant.Pos.X = ant.Pos.X - 1

	ant.Order.goal()
	if defend.Pos.X != defaultSize/2+2 && defend.Pos.Y != defaultSize/2+2 {
		t.Errorf("Wrong calculation of goal. Expected: %v, actual: %v", &pkg.Pos{X: defaultSize/2 + 2, Y: defaultSize/2 + 2}, defend.Pos)
	}
}
