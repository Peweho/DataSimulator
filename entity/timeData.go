package entity

import (
	"simulator/distribution"
)

type TimeData struct {
	Time int64 // 时间
	Data []float64
	Row  int
}

func CreateTimeData(T int64, row int, models *[]distribution.Model) *TimeData {
	td := &TimeData{Time: T, Row: row}
	td.Data = make([]float64, len(*models))
	for i := 0; i < len(*models); i++ {
		td.Data[i], _ = (*models)[i].Next()
	}
	return td
}
