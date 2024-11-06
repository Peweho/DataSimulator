package entity

import (
	"simulator/distribution"
)

type FarmEnvData struct {
	Time int64 // 时间
	Data []float64
}

func CreateFarmEnvData(T int64, models *[]distribution.Model) *FarmEnvData {
	fed := &FarmEnvData{Time: T}
	fed.Data = make([]float64, len(*models))
	for i := 0; i < len(*models); i++ {
		fed.Data[i], _ = (*models)[i].Next()
	}
	return fed
}

//func (f *FarmEnvData) String() string {
//	unix := time.Unix(f.Time, 0)
//	return fmt.Sprintf("时间：%s，"+
//		"温度：%.1f℃，"+
//		"湿度：%.1fKg/m3",
//		unix.Format(time.DateTime), f.Tempture, f.Humidity)
//}
