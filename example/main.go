package main

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

func NewAI(birthRelativePos *pkg.Pos, anthillID int) AI {
	anthillPos := &pkg.Pos{X: defaultSize / 2, Y: defaultSize / 2}
	birthRelativePos.Add(anthillPos)
	// for the beginning I guess that my birth point in the center of prospective area
	ai := AI{
		area:          NewArea(defaultSize, defaultSize),
		ants:          make(map[int]*Ant),
		anthills:      make(map[int]*Anthill),
		enemyAnthills: make(pkg.PosCollection, 0, 1),
	}
	ai.area.SetByPos(anthillPos, pkg.AllyAnthillField)
	ai.area.SetByPos(birthRelativePos, pkg.AllyField)

	ai.anthills[anthillID] = &Anthill{
		Pos:      &pkg.Pos{X: defaultSize / 2, Y: defaultSize / 2},
		BirthPos: birthRelativePos,
	}

	return ai
}

func (ai *AI) Start(anthillID int, birthPos pkg.Pos) {
	Greg = NewAI(&birthPos, anthillID)
}

func (ai *AI) Do(antID int, fields [5][5]pkg.FieldType, round int, posDiff pkg.Pos) (target *pkg.Pos, action pkg.Action) {
	currentAnt := ai.ants[antID]
	currentAnt.Pos.Add(&posDiff)
	ai.updateArea(fields, currentAnt)

	return GiveOrder(currentAnt, ai)
}

func (ai *AI) OnAntDie(antID int) {
	ai.area.SetByPos(ai.ants[antID].Pos, pkg.EmptyField)
	delete(ai.ants, antID)
}

func (ai *AI) OnAnthillDie(anthillID int) {
	ai.area.SetByPos(ai.anthills[anthillID].Pos, pkg.EnemyAnthillField)
	delete(ai.anthills, anthillID)
}

func (ai *AI) OnAntBirth(antID int, anthillID int) {
	ai.area.SetByPos(ai.anthills[anthillID].BirthPos, pkg.AllyField)
	ai.ants[antID] = &Ant{
		Pos:  &pkg.Pos{X: ai.anthills[anthillID].BirthPos.X, Y: ai.anthills[anthillID].BirthPos.Y},
		Role: ai.getActualRole(),
	}
}

func (ai *AI) OnNewAnthill(invaderID int, birthPos pkg.Pos, anthillID int) {
	pos := ai.area.Closest(ai.ants[invaderID].Pos, pkg.EnemyAnthillField)
	ai.area.SetByPos(pos, pkg.AllyAnthillField)
	birthPos.Add(pos)
	ai.anthills[anthillID] = &Anthill{
		Pos:      pos,
		BirthPos: &birthPos,
	}
}

// update information about real area on my prospective area
func (ai *AI) updateArea(fields [5][5]pkg.FieldType, current *Ant) {
	for dx := range fields {
		for dy, t := range fields[dx] {
			x := current.Pos.X + dx - 2
			y := current.Pos.Y + dy - 2
			if ai.area.RewriteMap(x, y, t, ai) {
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
		return explorer
	}
}

var Greg AI
