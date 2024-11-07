package distribution

import (
	"errors"
	"math/rand"
	"simulator/etc"
)

// 随机生成数据
type random struct {
	BaseModel
}

func (r *random) Next() (float64, error) {
	return rand.Float64()*(r.max-r.min) + r.min, nil
}

func NewRandom(content *etc.Data) (Model, error) {
	if len(content.Params) > 0 {
		return nil, errors.New("random模型不接受参数")
	}
	return &random{
		BaseModel: BaseModel{
			min: content.Min,
			max: content.Max,
		},
	}, nil
}
