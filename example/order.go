package example

import (
	pkg "github.com/gregmus2/ants-pkg"
	"math"
)

type Order interface {
	urgent() (*pkg.Pos, pkg.Action, bool)
	follow() (*pkg.Pos, pkg.Action)
	goal()
	hasGoal() bool
}

type BaseOrder struct {
	Ant    *Ant
	Pos    *pkg.Pos
	Action pkg.Action
	AI     *AI
}

type Explore struct {
	BaseOrder
}

type Attack struct {
	BaseOrder
}

type Defend struct {
	BaseOrder
	target *pkg.Pos
}

func GiveOrder(ant *Ant, greg *AI) (*pkg.Pos, pkg.Action) {
	if ant.Order == nil {
		ant.Order = GetOrder(ant, greg)
	}

	if pos, action, ok := ant.Order.urgent(); ok {
		return pos, action
	}

	if !ant.Order.hasGoal() {
		ant.Order.goal()
	}

	return ant.Order.follow()
}

func GetOrder(ant *Ant, ai *AI) Order {
	base := BaseOrder{Ant: ant, AI: ai}
	switch ant.Role {
	case explorer:
		return &Explore{BaseOrder: base}
	case defender:
		return &Defend{BaseOrder: base}
	case attacker:
		return &Attack{BaseOrder: base}
	}

	return nil
}

func (o *BaseOrder) urgent() (*pkg.Pos, pkg.Action, bool) {
	var foodPos *pkg.Pos
	var enemyPos *pkg.Pos
	for x := o.Ant.Pos.X - 2; x <= o.Ant.Pos.X+2; x++ {
		for y := o.Ant.Pos.Y - 2; y <= o.Ant.Pos.Y+2; y++ {
			switch o.AI.area.matrix[x][y] {
			// primary goal it's enemy anthill
			// todo add Enemy|Ally AnthillField logic to main app
			case pkg.EnemyAnthillField:
				return o.Ant.CalcStep(x, y), pkg.AttackAction, true
			case pkg.EnemyField:
				enemyPos = o.Ant.CalcStep(x, y)
			case pkg.FoodField:
				foodPos = o.Ant.CalcStep(x, y)
			}
		}
	}

	if enemyPos != nil {
		return enemyPos, pkg.AttackAction, true
	}

	if foodPos != nil {
		return foodPos, pkg.EatAction, true
	}

	return &pkg.Pos{}, 0, false
}

// it's about long-term goal based on Role. Like explorer or go to capture enemy anthill
func (o *BaseOrder) follow() (*pkg.Pos, pkg.Action) {
	// todo check for obstacles
	deltaX := 0
	if o.Pos.X != o.Ant.Pos.X {
		diffX := o.Pos.X - o.Ant.Pos.X
		deltaX = diffX / int(math.Abs(float64(diffX)))
	}

	deltaY := 0
	if o.Pos.Y != o.Ant.Pos.Y {
		diffY := o.Pos.Y - o.Ant.Pos.Y
		deltaY = diffY / int(math.Abs(float64(diffY)))
	}

	return &pkg.Pos{X: deltaX, Y: deltaY}, pkg.MoveAction
}

func (o *BaseOrder) hasGoal() bool {
	return o.Pos != nil && o.Pos != o.Ant.Pos
}

func (o *Defend) goal() {
	if o.target == nil {
		positions := make([]*pkg.Pos, 0, len(o.AI.anthills))
		for _, anthill := range o.AI.anthills {
			positions = append(positions, anthill.Pos)
		}
		o.target = o.Ant.Pos.CalcNearest(positions)
		/*
			Calc nearest guard position of target anthill
			p - guard position; t - target anthill

			p---p
			--t--
			p---p
		*/
		o.Pos = o.Ant.Pos.CalcNearest([]*pkg.Pos{
			{X: o.target.X - 2, Y: o.target.Y - 2},
			{X: o.target.X + 2, Y: o.target.Y - 2},
			{X: o.target.X + 2, Y: o.target.Y + 2},
			{X: o.target.X - 2, Y: o.target.Y + 2},
		})

		return
	}

	/*
		We need to calculate next guard position clockwise relative to the current
		-2|-2     +2|-2
			 p-->p
			 ^---|
			 |-t-|
			 |---v
			 p<--p
		-2|+2	  +2|+2
	*/

	relativeX := o.Ant.Pos.X - o.target.X
	relativeY := o.Ant.Pos.Y - o.target.Y
	o.Pos = &pkg.Pos{
		X: (0-relativeX/relativeY)*relativeX + o.target.X,
		Y: relativeX/relativeY*relativeY + o.target.Y,
	}
}

func (o *Explore) goal() {
	o.Pos = o.AI.area.Closest(o.Ant.Pos, unknownField)
}

// todo we need to create teams for attack
// todo test
func (o *Attack) goal() {
	// todo small steps
	o.Pos = o.Ant.Pos.CalcNearest(o.AI.enemyAnthills)
}
