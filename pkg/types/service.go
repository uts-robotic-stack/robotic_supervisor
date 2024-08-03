package types

type Service interface {
	Name() string
	Status() string
}
