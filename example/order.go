package example

import pkg "github.com/gregmus2/ants-pkg"

type order struct {
	ant       *ant
	hasOrder  bool
	pos       *pkg.Pos
	action    pkg.Action
	blackZone *pkg.Pos
}

func giveOrder(ant *ant) *order {
	o := &order{ant: ant}
	if o.urgent() {
		return o
	}

	o.calcBlackArea()
	o.goal()

	return o
}

func (o *order) get() (*pkg.Pos, pkg.Action) {
	return o.pos, o.action
}

func (o *order) urgent() bool {
	var foodPos *pkg.Pos
	var enemyPos *pkg.Pos
	for x := o.ant.Pos.X - 1; x <= o.ant.Pos.X+1; x++ {
		for y := o.ant.Pos.Y - 1; y <= o.ant.Pos.Y+1; y++ {
			switch area[x][y] {
			// primary goal it's enemy anthill
			// todo add Enemy|Ally AnthillField logic to main app
			case pkg.EnemyAnthillField:
				o.pos = o.ant.RelativePos(x, y)
				o.action = pkg.AttackAction
				return true
			case pkg.EnemyField:
				enemyPos = o.ant.RelativePos(x, y)
			case pkg.FoodField:
				foodPos = o.ant.RelativePos(x, y)
			}
		}
	}

	if enemyPos != nil {
		o.pos = enemyPos
		o.action = pkg.AttackAction
		return true
	}

	if foodPos != nil {
		// todo add birth handler (you should know about new ant)
		o.pos = foodPos
		o.action = pkg.EatAction
		return true
	}

	return false
}

// calculate the black area, your ant should avoid
func (o *order) calcBlackArea() {
	y := o.ant.Pos.Y - 3
	for x := o.ant.Pos.X - 3; x <= o.ant.Pos.X+3; x++ {
		// todo handle possibility of two enemies
		if area[x][y] == pkg.EnemyField {
			o.blackZone = &pkg.Pos{X: x, Y: y}
		}
	}
}

// it's about long-term goal based on role. Like explore or go to capture enemy anthill
func (o *order) goal() {

}
