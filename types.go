package pkg

type Algorithm interface {
	Start(anthill Pos, birth Pos)
	Do(fields [5][5]FieldType, round int) (target Pos, action Action)
}

type Pos [2]int

func (p *Pos) X() int {
	return p[0]
}

func (p *Pos) Y() int {
	return p[1]
}

func (p Pos) Add(pos Pos) {
	p[0] += pos.X()
	p[1] += pos.Y()
}
