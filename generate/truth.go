package generate

import (
	"context"
	"encoding/json"
	"fmt"
	"simulator/distribution"
	"simulator/entity"
	"simulator/etc"
	"simulator/mq/kafka"
	"simulator/util"
	"sync"
	"time"
)

// 按照采集频率分组
var (
	buckets    map[int64]*[]*etc.Data
	wg         *sync.WaitGroup
	ctx        context.Context
	cancelFunc context.CancelFunc
	kq         *kafka.KafkaClient
	config     *etc.Config
	locks      *sync.Map
)

func init() {
	buckets = make(map[int64]*[]*etc.Data, 0)
	wg = &sync.WaitGroup{}
	locks = &sync.Map{} //make(map[int64]*sync.Mutex, 0)
	ctx, cancelFunc = context.WithCancel(context.Background())
}

func TruthEnvGenerateData() error {
	config = etc.GetConfig("")
	kq = kafka.NewKafkaClient(config.Mq.Addr, config.Mq.Topic, config.Mq.Partition, config.Mq.TimeOut)
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
			locks.Store(data.Frequency, &sync.Mutex{})
		}
		bucket = buckets[data.Frequency]
		*bucket = append(*bucket, &config.Data[i])
	}
	// 启动定时人任务
	for frequency, data := range buckets {
		wg.Add(1)
		go generateDataByTime(time.Now().Unix(), frequency, data)
	}
	return nil
}

func generateDataByTime(start int64, frequency int64, data *[]*etc.Data) {
	defer wg.Done()
	// 每隔五秒循环一次，检查是否推出
	var loopTime int64 = 5
	// 离上一次产生数据时间
	countTime := frequency
	now := start
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if countTime < frequency {
				countTime += loopTime
				break
			}
			countTime = loopTime
			// 检查切片中是否还有任务，如果有生成数据
			lock(frequency)
			if len(*data) == 0 {
				delete(buckets, frequency)
				unlock(frequency)
				locks.Delete(frequency)
				return
			}
			datas := createTimeDatas(now, data)
			unlock(frequency)
			// 处理生产的数据
			go func() {
				marshal, err := json.Marshal(datas)
				if err != nil {
					util.Log.Fatalf("序列化数据错误：%v", err)
				}
				// 发送到kafka
				if err = kq.Write([]byte(""), marshal); err != nil {
					fmt.Println(string(marshal))
				}
			}()
		}
		time.Sleep(time.Duration(loopTime) * time.Second)
		now += loopTime
	}
}

func createTimeDatas(T int64, data *[]*etc.Data) *entity.TimeDatas {
	td := &entity.TimeDatas{Time: T, DataBaseId: config.DataBaseId}
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

func UpdateFrequency(index int, id string, oldFre, fre int64) {
	lock(oldFre)
	// 移除原来的桶中的数据
	bucket := buckets[oldFre]
	for i := 0; i < len(*bucket); i++ {
		if (*bucket)[i].Id == id {
			*bucket = append((*bucket)[:i], (*bucket)[i+1:]...)
			break
		}
	}
	unlock(oldFre)
	// 加入新桶
	// 判断是否需要开启协程
	bucket, ok := buckets[fre]
	if ok {
		lock(fre)
		*bucket = append((*bucket), &config.Data[index])
		unlock(fre)
	} else {
		datas := make([]*etc.Data, 0)
		datas = append(datas, &config.Data[index])
		buckets[fre] = &datas
		locks.Store(fre, &sync.Mutex{})
		wg.Add(1)
		go generateDataByTime(time.Now().Unix(), fre, &datas)
	}
}

func Stop() {
	cancelFunc()
	wg.Wait()
	util.Log.Println("停止生成数据")
}

func lock(key any) {
	value, ok := locks.Load(key)
	if !ok {
		util.Log.Fatalf("锁不存在")
	}
	lockVal := value.(*sync.Mutex)
	lockVal.Lock()
}

func unlock(key any) {
	value, ok := locks.Load(key)
	if !ok {
		util.Log.Fatalf("锁不存在")
	}
	lockVal := value.(*sync.Mutex)
	lockVal.Unlock()
}
