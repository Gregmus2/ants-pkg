package pkg

type Algorithm interface {
	Start(anthillID int, birthPos Pos)
	Do(antID int, fields [5][5]FieldType, round int, posDiff Pos) (target *Pos, action Action)
	OnAntDie(antID int)
	OnAnthillDie(anthillID int)
	OnAntBirth(antID int, anthillID int)
	OnNewAnthill(invaderID int, birthPos Pos, anthillID int) // antID; position relative anthill
}
