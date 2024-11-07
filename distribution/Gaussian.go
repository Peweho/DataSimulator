package distribution

import (
	"errors"
	"math/rand"
	"simulator/etc"
	"simulator/util"
)

// 基于正态分布生成数据
type Gaussian struct {
	BaseModel
	mean   float64
	stddev float64
}

func (g *Gaussian) Next() (float64, error) {
	return rand.NormFloat64()*g.stddev + g.mean, nil
}

func NewGaussian(content *etc.Data) (Model, error) {
	if len(content.Params) > 2 {
		return nil, errors.New("gaussian模型参数过多")
	}
	g := &Gaussian{
		BaseModel: BaseModel{
			min: content.Min,
			max: content.Max,
		},
		mean:   0,
		stddev: 1,
	}
	s := util.ParamsToFloat64s(content.Params)
	if len(s) > 0 {
		g.mean = s[0]
	}
	if len(s) > 1 {
		g.stddev = s[1]
	}

	return Model(g), nil
}
