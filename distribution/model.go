package distribution

import (
	"simulator/etc"
)

var (
	ModelMap = map[string]func(*etc.Data) (Model, error){
		"random":   NewRandom,
		"gaussian": NewGaussian,
		"exp":      NewExp,
	}
)

type Model interface {
	Next() (float64, error)
}

type BaseModel struct {
	min float64
	max float64
}

func NewBaseModel() BaseModel {
	return BaseModel{
		min: 0,
		max: 1,
	}
}
