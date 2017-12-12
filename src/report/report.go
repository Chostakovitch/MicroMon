//This package contains types and methods to format and write metrics with context.
package report

import (
	"fmt"
	"log"
	"os"

	"metric"
)

//Formatter defines how a metric formatter, i.e. a type which operates on Metric and associated Results, should behave.
type Formatter interface {
	//Single takes a Metric and its associated Result and produces a formatted string containing informations on both of them.
	Single(metric.WebMetric) string
	//Multiple takes a WebMetrics, i.e. Metrics with their Results for a website and a timeframe.
	//Multiple produces a formatted string which reports all theses informations.
	Multiple(metric.WebMetrics) string
}

//A Reporter is the association of a Logger and a Formatter.
//This type is meant to be generic, i.e. to allow formatting in any fashion and writing everywhere.
type Reporter struct {
	l *log.Logger
	f Formatter
}

//NewReported constructs a Reported from a Logger and a Formatter.
func NewReporter(l *log.Logger, f Formatter) Reporter {
	return Reporter{l, f}
}

//DefaultFormatter implements Formatter and is a provided classic formatter suited for console writing.
type DefaultFormatter struct{}

func (DefaultFormatter) Single(m metric.WebMetric) string {
	//Each Metric on a single line
	return m.M.Description() + " : " + m.R.Format(true)
}

func (f DefaultFormatter) Multiple(metrics metric.WebMetrics) string {
	//Website name
	res := fmt.Sprintf("=== %v (last %v minutes) ===\n", metrics.N, metrics.Timeframe)
	//Each metric with a tabulation
	for _, v := range metrics.M {
		res += fmt.Sprintf("\t%v\n", f.Single(v))
	}
	return res
}

//Report allows to format and write multiple metrics for multiple website computed within a given timeframe.
//order defines metrics order when formatting.
//metrics is map which associated a website name with Metrics and their corresponding Result.
//Report use the Formatter to format (Metric, Result)s and the Logger to write the final result.
func (r Reporter) Report(metrics []metric.WebMetrics) {
	res := "Reporting metrics :\n"
	for _, v := range metrics {
		res += r.f.Multiple(v)
	}
	r.l.Printf("%v", res)
}

//DefaultLogger is a convenient function which returns a pointer to a basic Logger, with [MicroMon] prefix and local datetime prefix.
func DefaultLogger() *log.Logger {
	return log.New(os.Stdout, "[MicroMon] ", log.LstdFlags)
}
