package pkg

type Algorithm interface {
	Do(fields [9]FieldType) (field uint8, action Action)
}
