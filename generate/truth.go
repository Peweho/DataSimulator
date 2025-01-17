package generate

import (
	"context"
	"encoding/json"
	"simulator/distribution"
	"simulator/entity"
	"simulator/etc"
	"simulator/util"
	"sync"
	"time"
)

// 按照采集频率分组
var (
	buckets    map[int64]*[]*etc.Data
	exit       chan struct{}
	wg         *sync.WaitGroup
	ctx        context.Context
	cancelFunc context.CancelFunc
)

func init() {
	buckets = make(map[int64]*[]*etc.Data, 0)
	wg = &sync.WaitGroup{}
	ctx, cancelFunc = context.WithCancel(context.Background())
}

func TruthEnvGenerateData() error {
	config := etc.GetConfig("")
	// 更具采集频率分组
	var bucket *[]*etc.Data
	for i := 0; i < len(config.Data); i++ {
		data := &config.Data[i]
		if data.Frequency == 0 {
			util.Log.Fatalf("采集频率不能为0")
		}
		if _, ok := buckets[data.Frequency]; !ok {
			datas := make([]*etc.Data, 0)
			buckets[data.Frequency] = &datas
		}
		bucket = buckets[data.Frequency]
		*bucket = append(*bucket, &config.Data[i])
	}
	// 启动定时人任务

	start, err := time.Parse(time.DateTime, config.Time.Start)
	if err != nil {
		util.Log.Fatalf("起始时间格式错误：%v", err)
	}

	exit = make(chan struct{}, len(buckets))
	for frequency, data := range buckets {
		wg.Add(1)
		go generateDataByTime(start.Unix(), frequency, data)
	}
	return nil
}

func generateDataByTime(start int64, frequency int64, data *[]*etc.Data) {
	// 每隔五秒循环一次，检查是否推出
	var loopTime int64 = 5
	// 离上一次产生数据时间
	countTime := frequency
	now := start
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		default:
			if countTime < frequency {
				countTime += loopTime
				break
			}
			countTime = loopTime
			datas := createTimeDatas(now, data)
			marshal, err := json.Marshal(datas)
			if err != nil {
				util.Log.Fatalf("序列化数据错误：%v", err)
			}
			util.Log.Println(string(marshal))
		}
		time.Sleep(time.Duration(loopTime) * time.Second)
		now += loopTime
	}
}

func createTimeDatas(T int64, data *[]*etc.Data) *entity.TimeDatas {
	td := &entity.TimeDatas{Time: T}
	models := setModelsBytruth(*data)
	td.Data = make([]entity.TimeData2, len(models))
	for i := 0; i < len(models); i++ {
		td.Data[i].Value, _ = models[i].Next()
		td.Data[i].Id = (*data)[i].Id
	}
	return td
}

func setModelsBytruth(data []*etc.Data) []distribution.Model {
	models := make([]distribution.Model, len(data))
	for i, c := range data {
		// 得到字段对应模型
		funcNewModel := distribution.ModelMap[c.Model]
		model, err := funcNewModel(c)
		if err != nil {
			util.Log.Fatalf("解析模型参数错误：%v", err)
			return nil
		}
		models[i] = model
	}
	return models
}

func Stop() {
	cancelFunc()
	wg.Wait()
	util.Log.Println("停止生成数据")
}
