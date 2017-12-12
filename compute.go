package micromon

import (
	"sync"
	"time"
)

//WebMetrics is a simple wrapper to associate a website name with its set of Metric and Result, for a given timeframe (in minutes)
type WebMetrics struct {
	Timeframe   int
	WebsiteName string
	Metrics     []WebMetric
}

//WebMetric associate a Metric with its Result.
type WebMetric struct {
	Source Metric
	Output Result
}

//respMap is just a map of website names associated with a safeData struct.
type respMap map[string]*safeData

//safeData is a slice of MetaResponse along with a mutex.
//As data can be processed from multiple threads (e.g. feeding, removing old data, reading, etc.,), sync is a must have.
type safeData struct {
	Datas []MetaResponse
	Mux   sync.Mutex
}

//NewRespMap initializes a new respMap suited for a given number of websites.
func NewRespMap(size int) respMap {
	return make(map[string]*safeData, size)
}

//NewSafeData initializes an empty safeData and returns a pointer to it.
func NewSafeData() *safeData {
	return &safeData{Datas: make([]MetaResponse, 0)}
}

//ComputeMetrics compute multiple metrics for a given timeframe and return the packed result (each element corresponds to a website with its metrics).
//It operates on a respMap struct, basically a set of websites names associated with multiple MetaResponse.
func (s *respMap) ComputeMetrics(metrics []Metric, minutes int) []WebMetrics {
	res := make([]WebMetrics, 0)

	//Iterate over each website data
	for k, v := range *s {
		//Copy data within the given timeframe
		v.Mux.Lock()
		datas := since(&v.Datas, minutes)
		v.Mux.Unlock()

		//If no data is available, do not compute
		if len(datas) == 0 {
			continue
		}
		tempRes := WebMetrics{minutes, k, make([]WebMetric, 0)}

		//For each metric asked, add result
		for _, m := range metrics {
			tempRes.Metrics = append(tempRes.Metrics, WebMetric{m, m.Compute(datas)})
		}

		res = append(res, tempRes)
	}
	return res
}

//since selects and returns all MetaResponse produced in the last X minutes, X given in function parameters.
func since(data *[]MetaResponse, minutes int) []MetaResponse {
	ret := make([]MetaResponse, 0)
	duration := time.Duration(time.Duration(minutes) * time.Minute)
	now := time.Now()
	for _, m := range *data {
		//Data is recent enough, select it
		if now.Sub(m.Timestamp) <= duration {
			ret = append(ret, m)
		}
	}
	return ret
}
