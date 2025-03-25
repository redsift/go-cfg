package testdata

type SomeStruct struct {
	SomeField int
}

func (s *SomeStruct) Some() int {
	return s.SomeField
}

type SomeInterface interface {
	Some() int
}
