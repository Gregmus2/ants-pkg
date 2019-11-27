package example

import (
	"github.com/gregmus2/ants-pkg"
)

// todo tests
type AI struct {
	area          Area // prospective area
	ants          map[int]*Ant
	anthills      map[int]*Anthill
	mapSize       int        // prospective size of area
	enemyAnthills []*pkg.Pos // todo handle destroy of anthill
}

const unknownField pkg.FieldType = 255
const defaultSize int = 100

func main() {

}

// todo add obsolete fields in area
func init() {
	Greg.mapSize = defaultSize
	Greg.area = NewArea(defaultSize, defaultSize)
	Greg.ants = make(map[int]*Ant)
	// for the beginning I guess that my birth point in the center of prospective area
	birthPoint := defaultSize / 2
	Greg.area[birthPoint][birthPoint] = pkg.AllyField
	Greg.anthills = make(map[int]*Anthill)
	Greg.enemyAnthills = make([]*pkg.Pos, 0, 1)
}

func (ai *AI) Do(antID int, fields [5][5]pkg.FieldType, round int, posDiff *pkg.Pos) (target *pkg.Pos, action pkg.Action) {
	// todo handle wrong eating, when two ants eat one food. If you send baseOrder about eating it's no mean that you
	// 	get new Ant. But how I can catch when new Ant birth?

	// todo handle dead ants. Maybe I need one more func in Algorithm (by events? or just call Die())

	currentAnt := ai.ants[antID]
	currentAnt.Pos.Add(posDiff)
	ai.updateArea(fields, currentAnt)

	return giveOrder(currentAnt, ai)
}

func (ai *AI) OnAntDie(antID int) {
	delete(ai.ants, antID)
}

func (ai *AI) OnAnthillDie(anthillID int) {
	delete(ai.anthills, anthillID)
}

func (ai *AI) OnAntBirth(antID int, anthillID int) {
	ai.ants[antID] = &Ant{
		Pos:  ai.anthills[anthillID].BirthPos,
		Role: ai.getActualRole(),
	}
}

func (ai *AI) OnNewAnthill(invaderID int, birthPos *pkg.Pos) {

}

// update information about real area on my prospective area
func (ai *AI) updateArea(fields [5][5]pkg.FieldType, current *Ant) {
	for dx := range fields {
		for dy, t := range fields[dx] {
			x := current.Pos.X + dx - 3
			y := current.Pos.Y + dy - 3
			if rewriteMap(x, y, t) {
				x = current.Pos.X + dx - 3
				y = current.Pos.Y + dy - 3
			}

			if ai.area[x][y] != t && t == pkg.EnemyAnthillField {
				ai.enemyAnthills = append(ai.enemyAnthills, &pkg.Pos{X: x, Y: y})
			}

			ai.area[x][y] = t
		}
	}
}

// todo if Ant go to explorer (no enemies or food near here) and somewhere we found food, Ant have to go there

// when we go beyond the intended map or get wall as a edge, we need to update our idea of map size
func rewriteMap(x, y int, t pkg.FieldType) bool {
	//t == pkg.NoField || t == pkg.WallField || x < 0 || y < 0 || x > mapSize || y > mapSize
}

// todo when we explored half part of map, we need to reorder all ants
func (ai *AI) getActualRole() Role {
	antCount := len(ai.ants)

	primaryRole := explorer
	if len(ai.enemyAnthills) > 0 {
		primaryRole = attacker
	}

	switch {
	case antCount < 6:
		return primaryRole
	case antCount >= 6 && antCount < 10:
		return defender
	default:
		return attacker
	}
}

var Greg AI
