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

type baseOrder struct {
	ant      *Ant
	hasOrder bool
	pos      *pkg.Pos
	action   pkg.Action
	ai       *AI
}

type Explore struct {
	baseOrder
}

type Attack struct {
	baseOrder
}

type Defend struct {
	baseOrder
	target *pkg.Pos
}

func giveOrder(ant *Ant, greg *AI) (*pkg.Pos, pkg.Action) {
	if ant.Order == nil {
		base := baseOrder{ant: ant, ai: greg}
		switch ant.Role {
		case explorer:
			ant.Order = &Explore{baseOrder: base}
		case defender:
			ant.Order = &Defend{baseOrder: base}
		case attacker:
			ant.Order = &Attack{baseOrder: base}
		}
	}

	if pos, action, ok := ant.Order.urgent(); ok {
		return pos, action
	}

	if !ant.Order.hasGoal() {
		ant.Order.goal()
	}

	return ant.Order.follow()
}

func (o *baseOrder) urgent() (*pkg.Pos, pkg.Action, bool) {
	var foodPos *pkg.Pos
	var enemyPos *pkg.Pos
	for x := o.ant.Pos.X - 2; x <= o.ant.Pos.X+2; x++ {
		for y := o.ant.Pos.Y - 2; y <= o.ant.Pos.Y+2; y++ {
			switch o.ai.area[x][y] {
			// primary goal it's enemy anthill
			// todo add Enemy|Ally AnthillField logic to main app
			case pkg.EnemyAnthillField:
				return o.ant.RelativeNearestPos(x, y), pkg.AttackAction, true
			case pkg.EnemyField:
				enemyPos = o.ant.RelativeNearestPos(x, y)
			case pkg.FoodField:
				foodPos = o.ant.RelativeNearestPos(x, y)
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
func (o *baseOrder) follow() (*pkg.Pos, pkg.Action) {
	// todo check for obstacles
	diffX := o.pos.X - o.ant.Pos.X
	deltaX := diffX / int(math.Abs(float64(diffX)))
	diffY := o.pos.Y - o.ant.Pos.Y
	deltaY := diffY / int(math.Abs(float64(diffY)))

	return &pkg.Pos{X: deltaX, Y: deltaY}, pkg.MoveAction
}

func (o *baseOrder) hasGoal() bool {
	return o.pos != nil && o.pos != o.ant.Pos
}

func (o *Defend) goal() {
	if o.target == nil {
		positions := make([]*pkg.Pos, 0, len(o.ai.anthills))
		for _, anthill := range o.ai.anthills {
			positions = append(positions, anthill.Pos)
		}
		o.target = o.ant.Pos.CalcNearest(positions)
		/*
			Calc nearest guard position of target anthill
			p - guard position; t - target anthill

			p---p
			--t--
			p---p
		*/
		o.pos = o.ant.Pos.CalcNearest([]*pkg.Pos{
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
	o.pos = &pkg.Pos{
		X: (0 - o.target.X/o.target.Y) * o.target.X,
		Y: o.target.X / o.target.Y * o.target.Y,
	}
}

func (o *Explore) goal() {
	o.pos = o.ai.area.Closest(o.ant.Pos, unknownField)
}

// todo we need to create teams for attack
func (o *Attack) goal() {
	// todo small steps
	o.pos = o.ant.Pos.CalcNearest(o.ai.enemyAnthills)
}
