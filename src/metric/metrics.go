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
type MaxRespTime struct{}
type CodeCount struct{}
type Availibility struct{}

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

func(MaxRespTime) Compute(data []urlwatch.MetaResponse) (Result) {
	max := time.Duration(0)
	for _, m := range data {
		if m.RespDuration > max {
			max = m.RespDuration
		}
	}
	return MetricFloat(float64(max) / float64(time.Millisecond))
}

func (MaxRespTime) Description() (string) {
	return "Maximum response time (ms)"
}

func (CodeCount) Compute(data []urlwatch.MetaResponse) (Result) {
	codes := make(map[string]int)
	for _, m := range data {
		codes[strconv.Itoa(m.Code)] += 1
	}
	res := make(MetricMap)
	for k, v := range codes {
		res[k] = MetricInt(v)
	}
	return res
}

func (CodeCount) Description() (string) {
	return "HTTP codes counts"
}

func (Availibility) Compute(data []urlwatch.MetaResponse) (Result) {
	count := 0
	for _, m := range(data) {
		if m.Available {
			count++
		}
	}
	return MetricFloat(float64(count) / float64(len(data)) * 100)
}

func (Availibility) Description() (string) {
	return "Availibility (%)"
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