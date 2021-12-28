package example

import "context"

type SomeAwesomeInterface interface {
	GetHoge()
	GetFugaWithContext(ctx context.Context)
}

type someAwesomeStruct struct {
}

func (s *someAwesomeStruct) GetHoge() {
	panic("implement me")
}

func (s *someAwesomeStruct) GetFugaWithContext(ctx context.Context) {
	panic("implement me")
}

func NewSomeAwesomeStruct() SomeAwesomeInterface {
	return &someAwesomeStruct{}
}
