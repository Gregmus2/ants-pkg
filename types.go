package pkg

import "math"

type Algorithm interface {
	Start(antID int)
	Do(antID int, fields [5][5]FieldType, round int, posDiff *Pos) (target *Pos, action Action)
}

type Pos struct{ X, Y int }

func (p Pos) Add(pos *Pos) {
	p.X += pos.X
	p.Y += pos.Y
}

func (p Pos) Set(x, y int) {
	p.X = x
	p.Y = y
}

func (p *Pos) CalcNearest(targets []*Pos) *Pos {
	minDist := p.CalcDist(targets[0])
	minPos := targets[0]
	for i := 1; i < len(targets); i++ {
		dist := p.CalcDist(targets[i])
		if dist < minDist {
			minDist = dist
			minPos = targets[i]
		}
	}

	return minPos
}

func (p *Pos) CalcDist(b *Pos) int {
	dist := math.Abs(float64(p.X - b.X))
	// because we can move by diagonal and move by x and y in one round
	if yDist := math.Abs(float64(p.Y - b.Y)); yDist > dist {
		dist = yDist
	}

	return int(dist)
}
