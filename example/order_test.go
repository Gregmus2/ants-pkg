package example

import (
	pkg "github.com/gregmus2/ants-pkg"
	"testing"
)

func setup() (AI, *Ant) {
	ai := NewAI()
	ant := &Ant{
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
		t.Errorf("Wrong calculated pos or action: %v, %d", pos, action)
	}
}

func TestGiveOrderWithUrgent(t *testing.T) {
	ai, ant := setup()
	ai.enemyAnthills = append(ai.enemyAnthills, &pkg.Pos{X: defaultSize - 10, Y: defaultSize / 2})
	ai.area.matrix[defaultSize/2-1][defaultSize/2+1] = pkg.FoodField

	pos, action := GiveOrder(ant, &ai)
	if pos.X != -1 || pos.Y != 1 || action != pkg.EatAction {
		t.Errorf("Wrong calculated pos or action: %v, %d", pos, action)
	}
}

func TestUrgent(t *testing.T) {
	ai, ant := setup()
	ant.Order = GetOrder(ant, &ai)

	pos, action, ok := ant.Order.urgent()
	if pos.X != 0 || pos.Y != 0 || action != 0 || ok == true {
		t.Errorf("Wrong calculated pos or action: %v, %d, %v", pos, action, ok)
	}

	ai.area.matrix[defaultSize/2-1][defaultSize/2+1] = pkg.FoodField

	pos, action, ok = ant.Order.urgent()
	if pos.X != -1 || pos.Y != 1 || action != pkg.EatAction || ok == false {
		t.Errorf("Wrong calculated pos or action: %v, %d", pos, action)
	}

	ai.area.matrix[defaultSize/2+1][defaultSize/2] = pkg.EnemyAnthillField
	pos, action, ok = ant.Order.urgent()
	if pos.X != 1 || pos.Y != 0 || action != pkg.AttackAction || ok == false {
		t.Errorf("Wrong calculated pos or action: %v, %d", pos, action)
	}
}

func TestUrgentEnemyPriority(t *testing.T) {
	ai, ant := setup()

	ai.area.matrix[defaultSize/2-1][defaultSize/2+1] = pkg.FoodField
	ai.area.matrix[defaultSize/2+1][defaultSize/2] = pkg.EnemyField

	pos, action := GiveOrder(ant, &ai)
	if action != pkg.AttackAction {
		t.Errorf("Wrong calculated action: %d", action)
	}

	if pos.X != 1 || pos.Y != 0 {
		t.Errorf("Wrong calculated pos or action. Expected: %v, actual: %v", &pkg.Pos{X: 1}, pos)
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
			t.Errorf("Wrong calculated action: %d", action)
		}

		if pos.X != step.X || pos.Y != step.Y {
			t.Errorf("Wrong calculated pos. Expected: %v, actual: %v", step, pos)
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

	ai.area.matrix[defaultSize/2-10][defaultSize/2-5] = unknownField
	ai.area.matrix[defaultSize/2-11][defaultSize/2-5] = unknownField

	ant.Order.goal()
	if explore.Pos.X != 39 && explore.Pos.Y != 45 {
		t.Errorf("Wrong calculated goal. Expected: %v, actual: %v", &pkg.Pos{X: 39, Y: 45}, explore.Pos)
	}

	ai.area.matrix[defaultSize/2+5][defaultSize/2-5] = unknownField

	ant.Order.goal()
	if explore.Pos.X != 39 && explore.Pos.Y != 45 {
		t.Errorf("Wrong calculated goal. Expected: %v, actual: %v", &pkg.Pos{X: 5, Y: -5}, explore.Pos)
	}
}
