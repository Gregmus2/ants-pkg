package pkg

type FieldType uint8

const NoField FieldType = 0
const EmptyField FieldType = 1
const FoodField FieldType = 2
const AllyField FieldType = 3
const EnemyField FieldType = 4
const WallField FieldType = 5
const AntField FieldType = 6
const AnthillField FieldType = 7

type Action uint8

const NoAction Action = 0
const AttackAction Action = 1
const EatAction Action = 2
const MoveAction Action = 3
const DieAction Action = 4
