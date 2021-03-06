package main

import (
	"github.com/gregmus2/ants-pkg"
	"log"
	"os"
)

// todo tests
type AI struct {
	area           *Area // prospective area
	ants           map[int]*Ant
	anthills       map[int]*Anthill
	enemyAnthills  pkg.PosCollection
	deviationTable [3][3][2]pkg.Pos
	log            *log.Logger
}

const unknownField pkg.FieldType = 255
const defaultSize int = 100

func main() {

}

func NewAI(birthRelativePos *pkg.Pos, anthillID int) AI {
	deviationTable := [3][3][2]pkg.Pos{
		{
			{{0, -1}, {-1, 0}},
			{{-1, -1}, {-1, 1}},
			{{0, 1}, {-1, 0}},
		},
		{
			{{-1, -1}, {1, -1}},
			{{0, 0}, {0, 0}},
			{{1, 1}, {-1, 1}},
		},
		{
			{{0, -1}, {1, 0}},
			{{1, 1}, {1, -1}},
			{{0, 1}, {1, 0}},
		},
	}

	anthillPos := &pkg.Pos{X: defaultSize / 2, Y: defaultSize / 2}
	birthRelativePos.Add(anthillPos)
	// for the beginning I guess that my birth point in the center of prospective area
	ai := AI{
		area:           NewArea(defaultSize, defaultSize),
		ants:           make(map[int]*Ant),
		anthills:       make(map[int]*Anthill),
		enemyAnthills:  make(pkg.PosCollection, 0, 1),
		deviationTable: deviationTable,
	}
	ai.area.SetByPos(anthillPos, pkg.AllyAnthillField)
	ai.area.SetByPos(birthRelativePos, pkg.AllyField)

	ai.anthills[anthillID] = &Anthill{
		Pos:      &pkg.Pos{X: defaultSize / 2, Y: defaultSize / 2},
		BirthPos: birthRelativePos,
	}

	f, _ := os.OpenFile("log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	ai.log = log.New(f, "", log.LstdFlags|log.Lshortfile)

	return ai
}

func (ai *AI) Start(anthillID int, birthPos pkg.Pos) {
	Greg = NewAI(&birthPos, anthillID)
}

func (ai *AI) Do(antID int, fields [5][5]pkg.FieldType, round int, posDiff pkg.Pos) (target pkg.Pos, action pkg.Action) {
	currentAnt := ai.ants[antID]
	currentAnt.Pos.Add(&posDiff)
	ai.updateArea(fields, currentAnt)

	target, action = GiveOrder(currentAnt, ai)
	if t := ai.area.GetRelative(currentAnt.Pos, target); (t == pkg.AllyAnthillField || t == pkg.AllyField) && !target.IsZero() {
		Greg.log.Printf("[%v] deviation before %v", currentAnt.Pos, target)
		target = ai.getDeviation(target)
		Greg.log.Printf("[%v] deviation %v", currentAnt.Pos, target)
	}

	return
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
		ID:   antID,
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
	ai.log.Printf("%d, %v updateArea", current.ID, current.Pos)
	for dx := range fields {
		for dy, t := range fields[dx] {
			x := current.Pos.X + dx - 2
			y := current.Pos.Y + dy - 2
			if ai.area.RewriteMap(x, y, ai) {
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

	ai.area.CutArea(fields, current.Pos, ai)
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

func (ai *AI) getDeviation(curDirection pkg.Pos) pkg.Pos {
	return ai.deviationTable[curDirection.X+1][curDirection.Y+1][ai.area.r.Intn(2)]
}

func (ai *AI) moveObjects(diff *pkg.Pos) {
	for _, ant := range ai.ants {
		ant.Pos.Add(diff)
	}
	for _, anthill := range ai.anthills {
		anthill.Pos.Add(diff)
	}
	for _, pos := range ai.enemyAnthills {
		pos.Add(diff)
	}
}

var Greg AI
