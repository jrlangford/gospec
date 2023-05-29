package gospectest

type TestCase[IN any, OUT any] struct {
	Description    string
	Input          IN
	ExpectedOutput OUT
}
