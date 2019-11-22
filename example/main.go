package example

import (
	"github.com/gregmus2/ants-pkg"
)

type greg struct {
}

type field struct {
	Type pkg.FieldType
}

type ant struct {
	Pos   *pkg.Pos
	Role  role
	Order *order
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

// prospective area
var area [][]pkg.FieldType
var ants map[int]*ant
var anthills []*pkg.Pos

// prospective size of area
var mapSize = defaultSize

// exploredPart identify that we know edges and can calculate the explored part
var exploredPart float32

const unknownField pkg.FieldType = 255
const defaultSize int = 100

func main() {

}

func (g greg) Start(antID int) {
	area = make([][]pkg.FieldType, defaultSize)
	for x := range area {
		area[x] = make([]pkg.FieldType, defaultSize)
		for y := range area[x] {
			area[x][y] = unknownField
		}
	}

	ants = make(map[int]*ant)
	// for the beginning I guess that my birth point in the center of prospective area
	birthPoint := defaultSize / 2
	// first ant exactly will birth in my birth point
	ants[antID] = &ant{
		Pos:  &pkg.Pos{X: birthPoint, Y: birthPoint},
		Role: explore,
	}
	area[birthPoint][birthPoint] = pkg.AllyField
}

func (g greg) Do(antID int, fields [5][5]pkg.FieldType, round int, posDiff *pkg.Pos) (target *pkg.Pos, action pkg.Action) {
	// todo handle wrong eating, when two ants eat one food. If you send order about eating it's no mean that you
	// 	get new ant. But how I can catch when new ant birth?

	// todo handle dead ants. Maybe I need one more func in Algorithm

	currentAnt := ants[antID]
	currentAnt.Pos.Add(posDiff)
	updateArea(fields, currentAnt)

	return giveOrder(currentAnt)
}

// update information about real area on my prospective area
func updateArea(fields [5][5]pkg.FieldType, current *ant) {
	for dx := range fields {
		for dy, t := range fields[dx] {
			x := current.Pos.X + dx - 3
			y := current.Pos.Y + dy - 3
			if rewriteMap(x, y, t) {
				x = current.Pos.X + dx - 3
				y = current.Pos.Y + dy - 3
			}

			area[x][y] = t
		}
	}
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
