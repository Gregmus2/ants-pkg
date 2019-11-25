package example

import (
	"github.com/gregmus2/ants-pkg"
	"math"
)

type greg struct {
	area         [][]pkg.FieldType // prospective area
	ants         map[int]*ant
	anthills     []*pkg.Pos
	mapSize      int // prospective size of area
	exploredPart float32
}

const unknownField pkg.FieldType = 255
const defaultSize int = 100

func main() {

}

func (g *greg) Start(antID int) {
	g.mapSize = defaultSize

	g.area = make([][]pkg.FieldType, defaultSize)
	for x := range g.area {
		g.area[x] = make([]pkg.FieldType, defaultSize)
		for y := range g.area[x] {
			g.area[x][y] = unknownField
		}
	}

	g.ants = make(map[int]*ant)
	// for the beginning I guess that my birth point in the center of prospective area
	birthPoint := defaultSize / 2
	// first ant exactly will birth in my birth point
	g.ants[antID] = &ant{
		Pos:  &pkg.Pos{X: birthPoint, Y: birthPoint},
		Role: explorer,
	}
	g.area[birthPoint][birthPoint] = pkg.AllyField
}

func (g *greg) Do(antID int, fields [5][5]pkg.FieldType, round int, posDiff *pkg.Pos) (target *pkg.Pos, action pkg.Action) {
	// todo handle wrong eating, when two ants eat one food. If you send baseOrder about eating it's no mean that you
	// 	get new ant. But how I can catch when new ant birth?

	// todo handle dead ants. Maybe I need one more func in Algorithm

	currentAnt := g.ants[antID]
	currentAnt.Pos.Add(posDiff)
	g.updateArea(fields, currentAnt)

	return giveOrder(currentAnt, g)
}

// update information about real area on my prospective area
func (g *greg) updateArea(fields [5][5]pkg.FieldType, current *ant) {
	for dx := range fields {
		for dy, t := range fields[dx] {
			x := current.Pos.X + dx - 3
			y := current.Pos.Y + dy - 3
			if rewriteMap(x, y, t) {
				x = current.Pos.X + dx - 3
				y = current.Pos.Y + dy - 3
			}

			g.area[x][y] = t
		}
	}
}

// todo if ant go to explorer (no enemies or food near here) and somewhere we found food, ant have to go there

// when we go beyond the intended map or get wall as a edge, we need to update our idea of map size
func rewriteMap(x int, y int, t pkg.FieldType) bool {
	//t == pkg.NoField || t == pkg.WallField || x < 0 || y < 0 || x > mapSize || y > mapSize
}

// todo if ant birth, they need to get role
// todo when we explored half part of map, we need to reorder all ants
func (g *greg) giveRole(unit *ant) {
	antCount := len(g.ants)

	primaryRole := explorer
	if g.exploredPart > 50 {
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

func calcDist(a *pkg.Pos, b *pkg.Pos) int {
	dist := math.Abs(float64(a.X - b.X))
	// because we can move by diagonal and move by x and y in one round
	if yDist := math.Abs(float64(a.Y - b.Y)); yDist > dist {
		dist = yDist
	}

	return int(dist)
}

func calcNearest(pos *pkg.Pos, targets []*pkg.Pos) *pkg.Pos {
	minDist := calcDist(targets[0], pos)
	minPos := targets[0]
	for i := 1; i < len(targets); i++ {
		dist := calcDist(targets[i], pos)
		if dist < minDist {
			minDist = dist
			minPos = targets[i]
		}
	}

	return minPos
}

var Greg greg
