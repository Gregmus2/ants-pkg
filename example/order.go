package example

import (
	pkg "github.com/gregmus2/ants-pkg"
	"math"
)

type order interface {
	urgent() (*pkg.Pos, pkg.Action, bool)
	follow() (*pkg.Pos, pkg.Action)
	goal()
	hasGoal() bool
}

type baseOrder struct {
	ant      *ant
	hasOrder bool
	pos      *pkg.Pos
	action   pkg.Action
	greg     *greg
}

type explore struct {
	baseOrder
}

type attack struct {
	baseOrder
}

type defend struct {
	baseOrder
	target *pkg.Pos
}

func giveOrder(ant *ant, greg *greg) (*pkg.Pos, pkg.Action) {
	if ant.Order == nil {
		base := baseOrder{ant: ant, greg: greg}
		switch ant.Role {
		case explorer:
			ant.Order = &explore{baseOrder: base}
		case defender:
			ant.Order = &defend{baseOrder: base}
		case attacker:
			ant.Order = &attack{baseOrder: base}
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
	for x := o.ant.Pos.X - 1; x <= o.ant.Pos.X+1; x++ {
		for y := o.ant.Pos.Y - 1; y <= o.ant.Pos.Y+1; y++ {
			switch o.greg.area[x][y] {
			// primary goal it's enemy anthill
			// todo add Enemy|Ally AnthillField logic to main app
			case pkg.EnemyAnthillField:
				return o.ant.RelativePos(x, y), pkg.AttackAction, true
			case pkg.EnemyField:
				enemyPos = o.ant.RelativePos(x, y)
			case pkg.FoodField:
				foodPos = o.ant.RelativePos(x, y)
			}
		}
	}

	if enemyPos != nil {
		return enemyPos, pkg.AttackAction, true
	}

	if foodPos != nil {
		// todo add birth handler (you should know about new ant)
		return foodPos, pkg.EatAction, true
	}

	return &pkg.Pos{}, 0, false
}

// it's about long-term goal based on role. Like explorer or go to capture enemy anthill
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

func (o *defend) goal() {
	if o.target == nil {
		o.target = calcNearest(o.ant.Pos, o.greg.anthills)
		/*
			Calc nearest guard position of target anthill
			p - guard position; t - target anthill

			p---p
			--t--
			p---p
		*/
		o.pos = calcNearest(o.ant.Pos, []*pkg.Pos{
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

func (o *explore) goal() {
	var matrix [20][20]bool
	matrix[8][9] = true
	fromFrame := 1
	toFrame := 5
	centerX := 6
	centerY := 6

	isChangeX := true
	for frame := fromFrame; frame <= toFrame; frame++ {
		from := -(frame - 1)
		X := from + centerX
		Y := -frame + centerY
		for polarity := 1; polarity >= -1; polarity -= 2 {
			for axis := 0; axis <= 1; axis++ {
				for j := from; j <= frame; j++ {
					if isChangeX {
						X = j*polarity + centerX
					} else {
						Y = j*polarity + centerY
					}

					if matrix[X][Y] == true {
						return X, Y
					}
				}
				isChangeX = !isChangeX
			}
		}
	}
}

func (o *attack) goal() {

}
