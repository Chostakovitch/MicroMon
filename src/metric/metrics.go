package metric

import (
	"urlwatch"
	"strconv"
	"time"
)

type Result interface {
	Format(inline bool) (string)
}

type Metric interface {
	Compute([]urlwatch.MetaResponse) (Result)
	Description() (string)
}

type AvgRespTime struct{}

func (AvgRespTime) Compute(data []urlwatch.MetaResponse) (Result) {
	sum := float64(0)
	for _, m := range data {
		sum += float64(m.RespDuration) / float64(time.Millisecond)
	}
	return MetricFloat(sum / float64(len(data)))
}

func (AvgRespTime) Description() (string) {
	return "Average response time (ms)"
}

type MetricInt int
type MetricFloat float64
type MetricMap map[string]Result

func (m MetricInt) Format(inline bool) (string) {
	return strconv.Itoa(int(m))
}

func (m MetricFloat) Format(inline bool) (string) {
	return strconv.FormatFloat(float64(m), 'f', 3, 64)
}

func (m MetricMap) Format(inline bool) (string) {
	var res string
	for k, v := range m {
		if inline {
			res += "{" + k + " : " + v.Format(inline) + "}"
		} else {
			res += k + " : " + v.Format(inline) + "\n"
		}
	}
	return res
}