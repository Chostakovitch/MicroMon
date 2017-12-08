//This package contains types and methods to format and write metrics with context.
package report

import (
	"log"
	"os"
	"strconv"

	"metric"
)

type Formatter interface {
	Single(metric.Metric, metric.Result) (string)
	Multiple([]metric.Metric, map[metric.Metric]metric.Result, string, int) (string)
}

type Reporter struct {
	l *log.Logger
	f Formatter
}

type DefaultFormatter struct{}

func (DefaultFormatter) Single(m metric.Metric, r metric.Result) (string) {
	return m.Description() + " : " + r.Format(true)
}

func (f DefaultFormatter) Multiple(order []metric.Metric, metrics map[metric.Metric]metric.Result, name string, since int) (string) {
	res := "Metrics for " + name + " (last " + strconv.Itoa(since) + " minute(s)) :\n"
	for _, v := range order {
		res += "\t" + f.Single(v, metrics[v]) + "\n"
	}
	return res
}

func NewReporter(l *log.Logger, f Formatter) (Reporter) {
	return Reporter{l, f}
}

func (r Reporter) Report(order []metric.Metric, metrics map[metric.Metric]metric.Result, name string, since int) {
	r.l.Print(r.f.Multiple(order, metrics, name, since))
}

func DefaultLogger() (*log.Logger) {
	return log.New(os.Stdout, "[MicroMon]",  log.LstdFlags)
}