package distribution

import (
	"errors"
	"math"
	"math/rand"
	"simulator/etc"
	"simulator/util"
)

type Exp struct {
	BaseModel
	lambda float64
}

// Next 返回指数分布数据
func (e *Exp) Next() (float64, error) {
	U := rand.Float64()
	if e.lambda < 0 {
		return e.max - math.Log(1-U)/e.lambda, nil
	}
	return -math.Log(1-U)/e.lambda + e.min, nil
}

func NewExp(content *etc.Data) (Model, error) {
	if len(content.Params) > 1 {
		return nil, errors.New("exp参数过多")
	}
	s := util.ParamsToFloat64s(content.Params)
	exp := &Exp{
		BaseModel: NewBaseModel(),
		lambda:    -0.5,
	}

	if len(s) > 0 {
		exp.lambda = s[0]
	}
	exp.min = content.Min
	exp.max = content.Max

	return Model(exp), nil
}
