package pkg

type Algorithm interface {
	Start(anthill Pos, birth Pos)
	Do(fields [5][5]FieldType, round int) (target *Pos, action Action)
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
