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

var (
	cfgPath = "./etc.yaml"
)

func main() {
	start := time.Now()
	//ConcurrencyRun()
	NotConcurrencyRun()
	util.Log.Printf("耗时：%v", time.Now().Sub(start).Nanoseconds())
}

func NotConcurrencyRun() {
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}
	// 解析配置文件
	config := etc.GetConfig(cfgPath)
	// 获取数据量和起始时间
	count, start := GetCountAndStart(config)
	// 生成数据字段对应的数据模型
	models := SetModels(config)
	// 创建传输数据管道
	ch := make(chan *entity.TimeData, count)
	// 启动生成excel文件的goroutine
	exitCreateExcel := make(chan struct{})
	go CreateExcel(ch, config, exitCreateExcel, start.Location())

	startUnix := start.Unix()
	internal := config.Time.Interval

	for i := 0; i < count; i++ {
		td := entity.CreateTimeData(startUnix+int64(i*internal), i+2, &models)
		ch <- td
	}
	close(ch)

	//等待生成excel文件的goroutine关闭
	<-exitCreateExcel
}

// 解析配置信息，得到生成数据量，开始的时间
func GetCountAndStart(config *etc.Config) (int, *time.Time) {
	// 计算生成的数据量
	start, err := time.Parse(time.DateTime, config.Time.Start)
	if err != nil {
		util.Log.Fatalf("起始时间格式错误：%v", err)
	}
	end, err := time.Parse(time.DateTime, config.Time.End)
	if err != nil {
		util.Log.Fatalf("结束时间格式错误：%v", err)
	}
	gap := end.Sub(start).Seconds()
	if gap <= 0 {
		util.Log.Fatalf("时间设置错误：%v", err)
	}
	count := int(gap/float64(config.Time.Interval)) + 1

	return count, &start
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

// 生成excel文件
func CreateExcel(datach <-chan *entity.TimeData, cfg *etc.Config, exit chan struct{}, loc *time.Location) {
	f := excelize.NewFile()
	sheetName := "simulator"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		panic(err)
	}

	// 设置表头
	f.SetCellValue(sheetName, "A1", "time")
	for i := 0; i < len(cfg.Data); i++ {
		name, _ := excelize.ColumnNumberToName(i + 2)
		f.SetCellValue(sheetName, fmt.Sprintf("%s%d", name, 1), cfg.Data[i].Title)
	}
	for data := range datach {
		// 第一列设置时间
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", data.Row), time.Unix(data.Time, 0).In(loc).Format(time.DateTime))
		for i := 1; i < len(data.Data)+1; i++ {
			name, _ := excelize.ColumnNumberToName(i + 1)
			f.SetCellValue(sheetName, fmt.Sprintf("%s%d", name, data.Row), math.Round(data.Data[i-1]*10)/10.0)
		}
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
