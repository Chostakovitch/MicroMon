//metric contains functions to aggregate MetaResponse into timed metrics and report them.
package metric

import (
	"urlwatch"
	"sync"
	"time"
)

type respMap map[string]*safeData

//safeData is slice of MetaResponse along with a mutex.
//As data can be processed from multiple threads (e.g. feeding, removing old data, reading, etc.,), sync is a must have.
type safeData struct {
	Datas []urlwatch.MetaResponse
	Mux sync.Mutex
}

//Mandatory constructor  to initialize map and get a valid pointer
func NewRespMap(size int) (respMap) {
	return make(map[string]*safeData, size)
}

func NewSafeData() (*safeData) {
	return &safeData{Datas: make([]urlwatch.MetaResponse, 0)}
}

//ComputeMetrics compute multiple metrics for a given timeframe and return the packed result.
//It operates on a respMap struct, so a set of websites associated with MetaResponse.
func (s *respMap) ComputeMetrics(metrics []Metric, minutes int) (map[string]map[Metric]Result) {
	res := make(map[string]map[Metric]Result)
	for k, v := range *s {
		v.Mux.Lock()
		datas := since(&v.Datas, minutes)
		v.Mux.Unlock()

		res[k] = make(map[Metric]Result)
		for _, m := range metrics {
			res[k][m] = m.Compute(datas)
		}
	}
	return res
}

func since(data *[]urlwatch.MetaResponse, minutes int) ([]urlwatch.MetaResponse) {
	ret := make([]urlwatch.MetaResponse, 0)
	duration := time.Duration(time.Duration(minutes) * time.Minute)
	now := time.Now()
	for _, m := range *data {
		if now.Sub(m.Timestamp) <= duration {
			ret = append(ret, m)
		}
	}
	return ret
}