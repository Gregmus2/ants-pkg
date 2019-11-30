package example

import (
	"github.com/gregmus2/ants-pkg"
)

// todo tests
type AI struct {
	area          *Area // prospective area
	ants          map[int]*Ant
	anthills      map[int]*Anthill
	enemyAnthills pkg.PosCollection
}

const unknownField pkg.FieldType = 255
const defaultSize int = 100

func main() {

}

func init() {
	Greg.area = NewArea(defaultSize, defaultSize)
	Greg.ants = make(map[int]*Ant)
	// for the beginning I guess that my birth point in the center of prospective area
	birthPoint := defaultSize / 2
	Greg.area.matrix[birthPoint][birthPoint] = pkg.AllyField
	Greg.anthills = make(map[int]*Anthill)
	Greg.enemyAnthills = make(pkg.PosCollection, 0, 1)
}

func (ai *AI) Do(antID int, fields [5][5]pkg.FieldType, round int, posDiff *pkg.Pos) (target *pkg.Pos, action pkg.Action) {
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
	// todo
}

// update information about real area on my prospective area
func (ai *AI) updateArea(fields [5][5]pkg.FieldType, current *Ant) {
	for dx := range fields {
		for dy, t := range fields[dx] {
			x := current.Pos.X + dx - 2
			y := current.Pos.Y + dy - 2
			if ai.area.RewriteMap(x, y, t) {
				x = current.Pos.X + dx - 2
				y = current.Pos.Y + dy - 2
			}

			if ai.area.matrix[x][y] != t {
				if t == pkg.EnemyAnthillField {
					ai.enemyAnthills = append(ai.enemyAnthills, &pkg.Pos{X: x, Y: y})
				} else if ai.area.matrix[x][y] == pkg.EnemyAnthillField {
					ai.enemyAnthills = ai.enemyAnthills.Remove(x, y)
				}
			}

			ai.area.matrix[x][y] = t
		}
	}
}

/*
	todo when we explored half part of map, we need to reorder all ants. When we have not unknown fields,
		we need to explore edges. When we know all edges,
		we no need explorers or we need to give them another algorithm
*/
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
