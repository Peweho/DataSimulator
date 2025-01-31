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

type TimeDatas struct {
	DataBaseId string      `json:"database_id"`
	Time       int64       `json:"time"` // 时间
	Data       []TimeData2 `json:"data"`
}

type TimeData2 struct {
	Id    string  `json:"id"`
	Value float64 `json:"value"`
}
