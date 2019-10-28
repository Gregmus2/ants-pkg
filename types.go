package pkg

type Algorithm interface {
	Start(anthill Pos, birth Pos)
	Do(fields [9]FieldType, round int) (field uint8, action Action)
}

type Pos [2]uint

func (p *Pos) X() uint {
	return p[0]
}

func (p *Pos) Y() uint {
	return p[1]
}
