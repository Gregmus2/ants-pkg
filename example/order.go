package example

import (
	pkg "github.com/gregmus2/ants-pkg"
	"math"
)

type order struct {
	ant      *ant
	hasOrder bool
	pos      *pkg.Pos
	action   pkg.Action
}

func giveOrder(ant *ant) (*pkg.Pos, pkg.Action) {
	if ant.Order == nil {
		ant.Order = &order{ant: ant}
	}

	if pos, action, ok := ant.Order.urgent(); ok {
		return pos, action
	}

	return ant.Order.follow()
}

func (o *order) urgent() (*pkg.Pos, pkg.Action, bool) {
	var foodPos *pkg.Pos
	var enemyPos *pkg.Pos
	for x := o.ant.Pos.X - 1; x <= o.ant.Pos.X+1; x++ {
		for y := o.ant.Pos.Y - 1; y <= o.ant.Pos.Y+1; y++ {
			switch area[x][y] {
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

// it's about long-term goal based on role. Like explore or go to capture enemy anthill
func (o *order) follow() (*pkg.Pos, pkg.Action) {
	if o.pos == nil || o.pos == o.ant.Pos {
		o.goal()
	}

	diffX := o.pos.X - o.ant.Pos.X
	deltaX := diffX / int(math.Abs(float64(diffX)))
	diffY := o.pos.Y - o.ant.Pos.Y
	deltaY := diffY / int(math.Abs(float64(diffY)))

	return &pkg.Pos{X: deltaX, Y: deltaY}, pkg.MoveAction
}

func (o *order) goal() {

}
