package example

import (
	pkg "github.com/gregmus2/ants-pkg"
	"testing"
)

func TestGiveOrder(t *testing.T) {
	ai := NewAI()
	ant := &Ant{
		Pos:  &pkg.Pos{X: defaultSize / 2, Y: defaultSize / 2},
		Role: attacker,
	}
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
	ai := NewAI()
	ant := &Ant{
		Pos:  &pkg.Pos{X: defaultSize / 2, Y: defaultSize / 2},
		Role: attacker,
	}
	ai.enemyAnthills = append(ai.enemyAnthills, &pkg.Pos{X: defaultSize - 10, Y: defaultSize / 2})
	ai.area.matrix[defaultSize/2-1][defaultSize/2+1] = pkg.FoodField

	pos, action := GiveOrder(ant, &ai)
	if pos.X != -1 || pos.Y != 1 || action != pkg.EatAction {
		t.Errorf("Wrong calculated pos or action: %v, %d", pos, action)
	}
}

func TestUrgent(t *testing.T) {
	ai := NewAI()
	ant := &Ant{
		Pos:  &pkg.Pos{X: defaultSize / 2, Y: defaultSize / 2},
		Role: attacker,
	}
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
	ai := NewAI()
	ant := &Ant{
		Pos:  &pkg.Pos{X: defaultSize / 2, Y: defaultSize / 2},
		Role: attacker,
	}

	ai.area.matrix[defaultSize/2-1][defaultSize/2+1] = pkg.FoodField
	ai.area.matrix[defaultSize/2+1][defaultSize/2] = pkg.EnemyField

	pos, action := GiveOrder(ant, &ai)
	if pos.X != 1 || pos.Y != 0 || action != pkg.AttackAction {
		t.Errorf("Wrong calculated pos or action: %v, %d", pos, action)
	}
}
