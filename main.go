package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"math"
	"os"
	"simulator/distribution"
	"simulator/entity"
	"simulator/etc"
	"simulator/util"
	"time"
)

func main() {
	cfgPath := "./etc.yaml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}
	// 解析配置文件
	config := etc.GetConfig(cfgPath)

	// 计算生成的数据量
	start, err := time.Parse(time.DateTime, config.Time.Start)
	if err != nil {
		util.Log.Fatalf("起始时间格式错误：%v", err)
		return
	}
	end, err := time.Parse(time.DateTime, config.Time.End)
	if err != nil {
		util.Log.Fatalf("结束时间格式错误：%v", err)
		return
	}
	gap := end.Sub(start).Seconds()
	if gap <= 0 {
		util.Log.Fatalf("时间设置错误：%v", err)
		return
	}
	count := int(gap / float64(config.Time.Interval))

	models := SetModels(config)

	startUnix := start.Unix()
	internal := config.Time.Interval

	ch := make(chan *entity.FarmEnvData, count)

	// 启动生成excel文件的goroutine
	exitCreateExcel := make(chan struct{})
	go CreateExcel(ch, config, exitCreateExcel)

	for i := 0; i < count; i++ {
		fed := entity.CreateFarmEnvData(startUnix+int64(i*internal), &models)
		ch <- fed
	}
	close(ch)

	//等待生成excel文件的goroutine关闭
	<-exitCreateExcel
}

// 解析每项数据段对应的模型
func SetModels(config *etc.Config) []distribution.Model {
	models := make([]distribution.Model, len(config.Data))
	for i, c := range config.Data {
		// 得到字段对应模型
		funcNewModel := distribution.ModelMap[c.Model]
		model, err := funcNewModel(&c)
		if err != nil {
			util.Log.Fatalf("解析模型参数错误：%v", err)
			return nil
		}
		models[i] = model
	}
	return models
}

func CreateExcel(datach <-chan *entity.FarmEnvData, cfg *etc.Config, exit chan struct{}) {
	f := excelize.NewFile()
	sheetName := "fed"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		panic(err)
	}

	// 设置表头
	f.SetCellValue(sheetName, "A1", "时间")
	for i := 0; i < len(cfg.Data); i++ {
		name, _ := excelize.ColumnNumberToName(i + 2)
		f.SetCellValue(sheetName, fmt.Sprintf("%s%d", name, 1), cfg.Data[i].Title)
	}
	rowNum := 2
	for data := range datach {
		// 第一列设置时间
		//fmt.Printf("%v\n", *data)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), time.Unix(data.Time, 0).Format(time.DateTime))
		for i := 1; i < len(data.Data)+1; i++ {
			name, _ := excelize.ColumnNumberToName(i + 1)
			f.SetCellValue(sheetName, fmt.Sprintf("%s%d", name, rowNum), math.Round(data.Data[i-1]*10)/10.0)
		}
		rowNum++
	}
	f.SetActiveSheet(index)

	// 保存文件
	if err := f.SaveAs("example.xlsx"); err != nil {
		util.Log.Fatalf("保存文件失败：%v", err)
	} else {
		util.Log.Println("生成数据完毕")
	}
	exit <- struct{}{}
}
