package main

import (
	pkg "github.com/gregmus2/ants-pkg"
	"math"
)

type Order interface {
	urgent() (pkg.Pos, pkg.Action, bool)
	follow() (pkg.Pos, pkg.Action)
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

func GiveOrder(ant *Ant, greg *AI) (pkg.Pos, pkg.Action) {
	if ant.Order == nil {
		ant.Order = GetOrder(ant, greg)
	}

	if pos, action, ok := ant.Order.urgent(); ok {
		Greg.log.Printf("[%v] urgent %d, %v", ant.Pos, action, pos)
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
	default:
		return &Explore{BaseOrder: base}
	}
}

func (o *BaseOrder) urgent() (pkg.Pos, pkg.Action, bool) {
	var foodPos, enemyPos *pkg.Pos
	var enemyIsGoal, foodIsGoal bool
	for x := o.Ant.Pos.X - 2; x <= o.Ant.Pos.X+2; x++ {
		for y := o.Ant.Pos.Y - 2; y <= o.Ant.Pos.Y+2; y++ {
			switch o.AI.area.matrix[x][y] {
			// primary goal it's enemy anthill
			// todo get nearest food or enemy
			case pkg.EnemyAnthillField:
				step, isGoal := o.Ant.CalcStep(x, y)
				return o.urgentResponse(*step, isGoal, pkg.AttackAction)
			case pkg.EnemyField:
				enemyPos, enemyIsGoal = o.Ant.CalcStep(x, y)
			case pkg.FoodField:
				foodPos, foodIsGoal = o.Ant.CalcStep(x, y)
			}
		}
	}

	if enemyPos != nil {
		return o.urgentResponse(*enemyPos, enemyIsGoal, pkg.AttackAction)
	}

	if foodPos != nil {
		return o.urgentResponse(*foodPos, foodIsGoal, pkg.EatAction)
	}

	return pkg.Pos{}, 0, false
}

func (o *BaseOrder) urgentResponse(pos pkg.Pos, isGoal bool, action pkg.Action) (pkg.Pos, pkg.Action, bool) {
	if isGoal {
		return pos, action, true
	} else {
		return pos, pkg.MoveAction, true
	}
}

// it's about long-term goal based on Role. Like explorer or go to capture enemy anthill
func (o *BaseOrder) follow() (pkg.Pos, pkg.Action) {
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

	Greg.log.Printf("[%v] follow %d, %d", o.Ant.Pos, deltaX, deltaY)
	return pkg.Pos{X: deltaX, Y: deltaY}, pkg.MoveAction
}

func (o *BaseOrder) hasGoal() bool {
	return o.Pos != nil && !o.Pos.Equal(o.Ant.Pos)
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
	var x, y int
	if relativeY == 0 {
		x = 0
		y = 0
	} else {
		x = 0 - relativeX/relativeY
		y = relativeX / relativeY
	}
	o.Pos = &pkg.Pos{
		X: x*relativeX + o.target.X,
		Y: y*relativeY + o.target.Y,
	}
}

func (o *Explore) goal() {
	o.Pos = o.AI.area.Closest(o.Ant.Pos, unknownField)
}

func (o *Attack) goal() {
	o.Pos = o.Ant.Pos.CalcNearest(o.AI.enemyAnthills)
}
