package httpService

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simulator/etc"
	"strconv"
)

type DataList struct {
	Data []DataDetail `json:"data"`
}

type DataDetail struct {
	Id        string `json:"id"`
	Frequency int64  `json:"frequency"`
}

// 获取监测数据信息
func dataList(ctx *gin.Context) {
	config := etc.GetConfig("")
	resp := &DataList{}
	resp.Data = make([]DataDetail, len(config.Data))

	for i, val := range config.Data {
		resp.Data[i].Id = val.Id
		resp.Data[i].Frequency = val.Frequency
	}
	ctx.JSON(http.StatusOK, resp)
}

// 设置修改数据采集频率
func setDataFrequency(ctx *gin.Context) {
	// 参数校验
	id := ctx.PostForm("id")
	frequency := ctx.PostForm("frequency")
	if id == "" || frequency == "" {
		HttpMsg(ctx, http.StatusBadRequest, "参数不完整")
		return
	}

	fre, err := strconv.ParseInt(frequency, 10, 64)
	if err != nil {
		HttpMsg(ctx, http.StatusBadRequest, "数据采集频率格式错误，整数，单位秒")
		return
	} else {
		config := etc.GetConfig("")
		for i, val := range config.Data {
			if val.Id == id {
				config.Data[i].Frequency = fre
				HttpMsg(ctx, http.StatusOK, "修改成功")
				return
			}
		}
	}
	HttpMsg(ctx, http.StatusBadRequest, "无效id")
}

// 绑定数据库
func bindDb(ctx *gin.Context) {
	config := etc.GetConfig("")
	dataBaseId := ctx.PostForm("id")
	if dataBaseId == "" {
		HttpMsg(ctx, http.StatusBadRequest, "参数不完整")
		return
	}
	config.DataBaseId = dataBaseId
	HttpMsg(ctx, http.StatusOK, "绑定成功")
}
