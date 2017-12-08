//metric contains functions to aggregate MetaResponse into timed metrics and report them.
package metric

import (
	"urlwatch"
	"sync"
	"time"
	"log"
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

//CalculateMetrics regularly produces metrics from MetaResponses which are added in safeData concurrently.
func CalculateMetrics(s *respMap, metrics []Metric) {
	i := 0
	for range time.Tick(10 * time.Second) {
		i++
		for k, v := range *s {
			v.Mux.Lock()
			datas := since(&v.Datas, 2)
			v.Mux.Unlock()
			for _, m := range metrics {
				log.Print(k + " : " + m.Description() + " : " + m.Compute(datas).Format(false))
			}
			if i % 6 == 0 {

			}
		}
	}
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