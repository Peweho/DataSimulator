package distribution

import (
	"errors"
	"simulator/etc"
)

var (
	errArgsNotSupport = errors.New("参数不足")
	ModelMap          = map[string]func(*etc.Data) (Model, error){
		"random":   NewRandom,
		"gaussian": NewGaussian,
		"exp":      NewExp,
	}
	// 模型对应的最长参数
	ArgsMap = map[string]int{
		"random":   0,
		"gaussian": 2,
		"exp":      1,
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
