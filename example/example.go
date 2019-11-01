package example

import (
	"github.com/gregmus2/ants-pkg"
)

type greg string

type field struct {
	Type pkg.FieldType
}

type order struct {
	ant       *ant
	hasOrder  bool
	pos       *pkg.Pos
	action    pkg.Action
	blackZone *pkg.Pos
}

type ant struct {
	Pos  *pkg.Pos
	Role role
}

func (a *ant) RelativePos(x int, y int) *pkg.Pos {
	return &pkg.Pos{
		X: x - a.Pos.X + 3,
		Y: y - a.Pos.Y + 3,
	}
}

type role uint8

const explore role = 0
const defend role = 1
const attack role = 2

var area [][]pkg.FieldType
var ants []*ant
var anthills []*pkg.Pos
var roundCounter = 0
var antOrder uint

var mapSize = defaultSize

// exploredPart identify that we know edges and can calculate the explored part
var exploredPart float32

const unknownField pkg.FieldType = 255
const defaultSize int = 100

func main() {

}

func init() {
	area = make([][]pkg.FieldType, defaultSize)
	for x := range area {
		area[x] = make([]pkg.FieldType, defaultSize)
		for y := range area[x] {
			area[x][y] = unknownField
		}
	}

	ants = make([]*ant, 1, 10)
	birthPoint := int(defaultSize / 2)
	ants[0] = &ant{
		Pos:  &pkg.Pos{X: birthPoint, Y: birthPoint},
		Role: explore,
	}
	area[birthPoint][birthPoint] = pkg.AllyField
}

func (g greg) Do(fields [5][5]pkg.FieldType, round int) (target *pkg.Pos, action pkg.Action) {
	if round != roundCounter {
		roundCounter = round
		antOrder = 0
	}

	// todo catch first anthill, birth point

	// update map
	currentAnt := ants[antOrder]
	antOrder++
	for dx := range fields {
		for dy, t := range fields[dx] {
			x := currentAnt.Pos.X + dx - 3
			y := currentAnt.Pos.Y + dy - 3
			if rewriteMap(x, y, t) {
				x = currentAnt.Pos.X + dx - 3
				y = currentAnt.Pos.Y + dy - 3
			}

			area[x][y] = t
		}
	}

	return giveOrder(currentAnt).get()
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

// todo if ant go to explore (no enemies or food near here) and somewhere we found food, ant have to go there

// when we go beyond the intended map or get wall as a edge, we need to update our idea of map size
func rewriteMap(x int, y int, t pkg.FieldType) bool {
	//t == pkg.NoField || t == pkg.WallField || x < 0 || y < 0 || x > mapSize || y > mapSize
}

// todo if ant birth, they need to get role
// todo when we explored half part of map, we need to reorder all ants
func giveRole(unit *ant) {
	antCount := len(ants)

	primaryRole := explore
	if exploredPart > 50 {
		primaryRole = attack
	}

	switch {
	case antCount < 6:
		unit.Role = primaryRole
	case antCount >= 6 && antCount < 10:
		unit.Role = defend
	default:
		unit.Role = attack
	}
}

var Greg greg
