package util

import "strconv"

func ParamsToFloat64s(params []string) []float64 {
	args := make([]float64, len(params))
	var err error
	for i, v := range params {
		args[i], err = strconv.ParseFloat(v, 64)
		if err != nil {
			Log.Println("参数转换失败：", v)
			return nil
		}
	}
	return args
}
