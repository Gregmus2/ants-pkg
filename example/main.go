package example

import (
	"github.com/gregmus2/ants-pkg"
)

// todo tests
type AI struct {
	area          Area // prospective area
	ants          map[int]*ant
	anthills      []*pkg.Pos
	mapSize       int        // prospective size of area
	enemyAnthills []*pkg.Pos // todo handle destroy of anthill
}

const unknownField pkg.FieldType = 255
const defaultSize int = 100

func main() {

}

// todo add obsolete fields in area
func (ai *AI) Start(antID int) {
	ai.mapSize = defaultSize
	ai.area = NewArea(defaultSize, defaultSize)
	ai.ants = make(map[int]*ant)
	// for the beginning I guess that my birth point in the center of prospective area
	birthPoint := defaultSize / 2
	// first ant exactly will birth in my birth point
	ai.ants[antID] = &ant{
		Pos:  &pkg.Pos{X: birthPoint, Y: birthPoint},
		Role: explorer,
	}
	ai.area[birthPoint][birthPoint] = pkg.AllyField
	ai.anthills = make([]*pkg.Pos, 0, 1)
	ai.enemyAnthills = make([]*pkg.Pos, 0, 1)
}

func (ai *AI) Do(antID int, fields [5][5]pkg.FieldType, round int, posDiff *pkg.Pos) (target *pkg.Pos, action pkg.Action) {
	// todo handle wrong eating, when two ants eat one food. If you send baseOrder about eating it's no mean that you
	// 	get new ant. But how I can catch when new ant birth?

	// todo handle dead ants. Maybe I need one more func in Algorithm

	currentAnt := ai.ants[antID]
	currentAnt.Pos.Add(posDiff)
	ai.updateArea(fields, currentAnt)

	return giveOrder(currentAnt, ai)
}

// update information about real area on my prospective area
func (ai *AI) updateArea(fields [5][5]pkg.FieldType, current *ant) {
	for dx := range fields {
		for dy, t := range fields[dx] {
			x := current.Pos.X + dx - 3
			y := current.Pos.Y + dy - 3
			if rewriteMap(x, y, t) {
				x = current.Pos.X + dx - 3
				y = current.Pos.Y + dy - 3
			}

			if ai.area[x][y] != t {
				switch t {
				case pkg.AllyAnthillField:
					ai.anthills = append(ai.anthills, &pkg.Pos{X: x, Y: y})
				case pkg.EnemyAnthillField:
					ai.enemyAnthills = append(ai.enemyAnthills, &pkg.Pos{X: x, Y: y})
				}
			}

			ai.area[x][y] = t
		}
	}
}

// todo if ant go to explorer (no enemies or food near here) and somewhere we found food, ant have to go there

// when we go beyond the intended map or get wall as a edge, we need to update our idea of map size
func rewriteMap(x, y int, t pkg.FieldType) bool {
	//t == pkg.NoField || t == pkg.WallField || x < 0 || y < 0 || x > mapSize || y > mapSize
}

// todo if ant birth, they need to get Role
// todo when we explored half part of map, we need to reorder all ants
func (ai *AI) giveRole(unit *ant) {
	antCount := len(ai.ants)

	primaryRole := explorer
	if len(ai.enemyAnthills) > 0 {
		primaryRole = attacker
	}

	switch {
	case antCount < 6:
		unit.Role = primaryRole
	case antCount >= 6 && antCount < 10:
		unit.Role = defender
	default:
		unit.Role = attacker
	}
}

var Greg AI
