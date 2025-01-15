package httpService

import (
	"github.com/gin-gonic/gin"
)

func Start() error {
	r := gin.Default()

	// 获取监测数据信息
	r.GET("/list", dataList)

	// 设置修改数据采集频率
	r.POST("/set/frequency", setDataFrequency)

	// 绑定数据库
	r.POST("/bind/db", bindDb)

	err := r.Run(":36664")
	if err != nil {
		return err
	}
	return nil
}
