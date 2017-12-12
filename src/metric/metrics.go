package metric

import (
	"fmt"
	"strconv"
	"time"
	"urlwatch"
)

//Result defines what should do a metric result.
type Result interface {
	//Format returns a formatted string representing the Result.
	//The inline parameters defines whether the resulting string should contain linebreaks or not.
	Format(inline bool) string
}

//Metric defines how a metric should behave.
type Metric interface {
	//Compute takes a slice of MetaResponse and produce a single aggregated Result.
	Compute([]urlwatch.MetaResponse) Result

	//Description returns a string describing what does the Metric.
	Description() string

	//Name returns the name of the Metric.
	Name() string
}

//GetMetric allows to instantiate a Metric from a string and return it.
//If no corresponding Metric is found, a non-nil error is returned.
func GetMetric(name string) (Metric, error) {
	switch name {
	case "averageTime":
		return AvgRespTime{}, nil
	case "maxTime":
		return MaxRespTime{}, nil
	case "codeCount":
		return CodeCount{}, nil
	case "availability":
		return Availability{}, nil
	}
	return nil, fmt.Errorf("%s is not a known metric name", name)
}

//AvgRespTime implements Metric and compute the average response time.
type AvgRespTime struct{}

//MaxRespTime implements Metric and compute the maximum response time.
type MaxRespTime struct{}

//CodeCount implements Metric and counts occurrence of HTTP response code.
type CodeCount struct{}

//Availability implements Metric and compute the percentage of availability.
type Availability struct{}

func (AvgRespTime) Compute(data []urlwatch.MetaResponse) Result {
	sum := float64(0)
	for _, m := range data {
		sum += float64(m.RespDuration) / float64(time.Millisecond)
	}
	return MetricFloat(sum / float64(len(data)))
}

func (AvgRespTime) Description() string {
	return "Average response time (ms)"
}

func (AvgRespTime) Name() string {
	return "averageTime"
}

func (MaxRespTime) Compute(data []urlwatch.MetaResponse) Result {
	max := time.Duration(0)
	for _, m := range data {
		if m.RespDuration > max {
			max = m.RespDuration
		}
	}
	return MetricFloat(float64(max) / float64(time.Millisecond))
}

func (MaxRespTime) Description() string {
	return "Maximum response time (ms)"
}

func (MaxRespTime) Name() string {
	return "maxTime"
}

func (CodeCount) Compute(data []urlwatch.MetaResponse) Result {
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

func (CodeCount) Description() string {
	return "HTTP codes counts"
}

func (CodeCount) Name() string {
	return "codeCount"
}

func (Availability) Compute(data []urlwatch.MetaResponse) Result {
	count := 0
	for _, m := range data {
		if m.Available {
			count++
		}
	}
	return MetricFloat(float64(count) / float64(len(data)) * 100)
}

func (Availability) Description() string {
	return "Availability (%)"
}

func (Availability) Name() string {
	return "availability"
}

//MetricInt implements Result and represents an integer result.
type MetricInt int

//MetricFloat implements Result and represents a float result.
type MetricFloat float64

//MetricMap implements Result and is a composite type, i.e. a map of string associated with Result.
//Keys are allowed to be another MetricMap.
type MetricMap map[string]Result

func (m MetricInt) Format(inline bool) string {
	return strconv.Itoa(int(m))
}

func (m MetricFloat) Format(inline bool) string {
	return strconv.FormatFloat(float64(m), 'f', 3, 64)
}

func (m MetricMap) Format(inline bool) string {
	var res string
	//As keys can be composite, call Format for each Result.
	for k, v := range m {
		if inline {
			res += "{" + k + " : " + v.Format(inline) + "}"
		} else {
			res += k + " : " + v.Format(inline) + "\n"
		}
	}
	return res
}
