package server

import (
	"encoding/json"
	"fmt"
	"sync"

	"insane/general/base/appconfig"
	"insane/utils"
)

type Report struct {
	RequestTime       uint64         `json:"requestTime"`       // 请求总时间
	MaxTime           uint64         `json:"maxTime"`           // 最大时长
	MinTime           uint64         `json:"minTime"`           // 最小时长
	SuccessNum        uint64         `json:"successNum"`        // 成功请求数
	FailureNum        uint64         `json:"failureNum"`        // 失败请求数
	NewOrderNum       uint64         `json:"newOrderNum"`       // tpmc
	ConCurrency       uint64         `json:"conCurrency"`       // 并发数
	ErrCode           map[int]int    `json:"errCode"`           // 错误码/错误个数
	ErrCodeMsg        map[int]string `json:"errCodeMsg"`        // 错误码描述
	Status            bool           `json:"status"`
	m                 sync.Mutex
}

func (report *Report) ReceivingResults(id string, conCurrency uint64, ch <-chan *Response, wgReceiving *sync.WaitGroup) {
	defer wgReceiving.Done()
	// 时间
	var (
		maxTime           uint64                 // 最大时长
		minTime           uint64                 // 最小时长
		successNum        uint64                 // 成功请求数
		failureNum        uint64                 // 失败请求数
		newOrderNum       uint64                 // tpmc
		errCode           = make(map[int]int)    // 错误码/错误个数
		errCodeMsg        = make(map[int]string) // 错误码描述
	)

	startTime := utils.Now()
	for data := range ch {
		report.m.Lock()

		if data.IsSuccess {
			successNum++
			if data.Transaction == 1 {
				newOrderNum++
			}
			if data.WasteTime > maxTime {
				maxTime = data.WasteTime
			}
			if minTime == 0 || data.WasteTime < minTime {
				minTime = data.WasteTime
			}
		} else {
			errCode[data.ErrCode]++
			if _, ok := errCodeMsg[data.ErrCode]; !ok {
				errCodeMsg[data.ErrCode] = data.ErrMsg
			} else {
				if errCodeMsg[data.ErrCode] != data.ErrMsg {
					errCodeMsg[data.ErrCode+1] = data.ErrMsg
				}
			}
			failureNum++
		}

		report.MaxTime = maxTime
		report.MinTime = minTime
		report.SuccessNum = successNum
		report.FailureNum = failureNum
		report.NewOrderNum = newOrderNum
		report.ConCurrency = conCurrency
		report.ErrCode = errCode
		report.ErrCodeMsg = errCodeMsg

		report.m.Unlock()
	}
	endTime := utils.Now()
	report.RequestTime = uint64((endTime - startTime) / 1000)
	report.Status = true

	content, err := json.Marshal(report)
	if err == nil {
		filename := fmt.Sprintf("%s/%s.json", appconfig.GetConfig().Log.Location, id)
		utils.FileWrite(filename, string(content))
	}
}

func (report *Report) Get() (content string) {
	report.m.Lock()
	defer report.m.Unlock()
	con, err := json.Marshal(report)
	if err != nil {
		return ""
	}
	return string(con)

	//filename := fmt.Sprintf("%s/%s.json", appconfig.GetConfig().Log.Location, id)
	//content, err := utils.FileGet(filename)
	//if err != nil {
	//	return ""
	//}
	//return

}