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
	Single(metric.Metric, metric.Result) string
	//Multiple takes a slice of Metrics, defined the process ordering, as such as a map of Metric associated with Result.
	//The last string is the name of the object concerned by the Metrics.
	//Multiple produces a formatted string which reports all theses informations.
	Multiple([]metric.Metric, map[metric.Metric]metric.Result, string) string
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

func (DefaultFormatter) Single(m metric.Metric, r metric.Result) string {
	//Each Metric on a single line
	return m.Description() + " : " + r.Format(true)
}

func (f DefaultFormatter) Multiple(order []metric.Metric, metrics map[metric.Metric]metric.Result, name string) string {
	//Website name
	res := fmt.Sprintf("=== %v ===\n", name)
	//Each metric with a tabulation
	for _, v := range order {
		res += fmt.Sprintf("\t%v\n", f.Single(v, metrics[v]))
	}
	return res
}

//Report allows to format and write multiple metrics for multiple website computed within a given timeframe.
//order defines metrics order when formatting.
//metrics is map which associated a website name with Metrics and their corresponding Result.
//since is an integer which represents a timeframe in minutes.
//Report use the Formatter to format (Metric, Result)s and the Logger to write the final result.
func (r Reporter) Report(order []metric.Metric, metrics map[string]map[metric.Metric]metric.Result, since int) {
	res := fmt.Sprintf("Metrics for the last %v minute(s) :\n", since)
	for k, v := range metrics {
		res += r.f.Multiple(order, v, k)
	}
	r.l.Printf("%v", res)
}

//DefaultLogger is a convenient function which returns a pointer to a basic Logger, with [MicroMon] prefix and local datetime prefix.
func DefaultLogger() *log.Logger {
	return log.New(os.Stdout, "[MicroMon] ", log.LstdFlags)
}
