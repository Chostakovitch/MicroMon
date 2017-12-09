//metric contains methods to aggregate MetaResponse into timed metrics and report them.
//
//metric defines several types. First, Metric and Result are primitives to aggregate MetaResponse.
//Second, respMap associates a website's name with a safeData, which is a thread-compatible slice of MetaResponse.
//
//The most important method is respMap.ComputeMetrics, which compute metrics for each website within a given interval.
package metric

import (
	"sync"
	"time"
	"urlwatch"
)

//respMap is just a map of website names associated with a safeData struct.
type respMap map[string]*safeData

//safeData is slice of MetaResponse along with a mutex.
//As data can be processed from multiple threads (e.g. feeding, removing old data, reading, etc.,), sync is a must have.
type safeData struct {
	Datas []urlwatch.MetaResponse
	Mux   sync.Mutex
}

//NewRespMap initializes a new respMap suited for a given number of websites.
func NewRespMap(size int) respMap {
	return make(map[string]*safeData, size)
}

//NewSafeData initializes an empty safeData and returns a pointer to it.
func NewSafeData() *safeData {
	return &safeData{Datas: make([]urlwatch.MetaResponse, 0)}
}

//ComputeMetrics compute multiple metrics for a given timeframe and return the packed result.
//It operates on a respMap struct, basically a set of websites names associated with MetaResponse.
func (s *respMap) ComputeMetrics(metrics []Metric, minutes int) map[string]map[Metric]Result {
	res := make(map[string]map[Metric]Result)

	//Iterate over each website data
	for k, v := range *s {
		//Copy data within the given timeframe
		v.Mux.Lock()
		datas := since(&v.Datas, minutes)
		v.Mux.Unlock()

		//If not data is available we do not compute anything
		if len(datas) == 0 {
			continue
		}
		res[k] = make(map[Metric]Result)

		//For each metric asked, add result
		for _, m := range metrics {
			res[k][m] = m.Compute(datas)
		}
	}
	return res
}

//since selects and returns all MetaResponse produced in the last X minutes, X given in function parameters.
func since(data *[]urlwatch.MetaResponse, minutes int) []urlwatch.MetaResponse {
	ret := make([]urlwatch.MetaResponse, 0)
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
