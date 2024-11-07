package entity

import (
	"simulator/distribution"
)

type TimeData struct {
	Time int64 // 时间
	Data []float64
}

func CreateFarmEnvData(T int64, models *[]distribution.Model) *TimeData {
	fed := &TimeData{Time: T}
	fed.Data = make([]float64, len(*models))
	for i := 0; i < len(*models); i++ {
		fed.Data[i], _ = (*models)[i].Next()
	}
	return fed
}
