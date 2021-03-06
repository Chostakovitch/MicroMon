package micromon

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

//Formatter defines how a metric formatter, i.e. a type which operates on Metric and associated Results, should behave.
type Formatter interface {
	//Single takes a Metric and its associated Result and produces a formatted string containing informations on both of them.
	Single(WebMetric) string
	//Multiple takes a WebMetrics, i.e. Metrics with their Results for a website and a timeframe.
	//Multiple produces a formatted string which reports all theses informations.
	Multiple(WebMetrics) string
	//Prefix returns a string which should precede a global reporting.
	Prefix() string
	//Suffix returns a string which should follow a global reporting.
	Suffix() string
}

//Reporter is the association of a Logger and a Formatter.
//This type is meant to be generic, i.e. to allow formatting in any fashion and writing everywhere.
type Reporter struct {
	l *log.Logger
	f Formatter
}

//NewReported constructs a Reporter from a Logger and a Formatter.
func NewReporter(l *log.Logger, f Formatter) Reporter {
	return Reporter{l, f}
}

//DefaultFormatter implements Formatter and is a provided classic formatter suited for console writing.
type DefaultFormatter struct{}

//XMLFormatter implements Formatter which formats metrics in XML suited for later parsing.
type XMLFormatter struct{}

func (DefaultFormatter) Single(m WebMetric) string {
	//Each Metric on a single line
	return m.Source.Description() + " : " + m.Output.Format(true)
}

func (f DefaultFormatter) Multiple(m WebMetrics) string {
	//Website name
	res := fmt.Sprintf("=== %v (last %v minutes) ===\n", m.WebsiteName, m.Timeframe)
	//Each metric with a tabulation
	for _, v := range m.Metrics {
		res += fmt.Sprintf("\t%v\n", f.Single(v))
	}
	return res
}

func (DefaultFormatter) Prefix() string {
	return "Reporting computed metrics :\n"
}

func (DefaultFormatter) Suffix() string {
	return ""
}

func (XMLFormatter) Single(m WebMetric) string {
	return fmt.Sprintf("<metric><name>%v</name><description>%v</description><value>%v</value></metric>", m.Source.Name(), m.Source.Description(), m.Output.Format(true))
}

func (f XMLFormatter) Multiple(m WebMetrics) string {
	res := fmt.Sprintf("<metrics><website><name>%v</name></website><timeframe>%v</timeframe>", m.WebsiteName, m.Timeframe)
	for _, v := range m.Metrics {
		res += f.Single(v)
	}
	res += fmt.Sprintf("</metrics>")
	//Code taken from : https://stackoverflow.com/a/21117347
	x := node{}
	xml.Unmarshal([]byte(res), &x)
	buf, _ := xml.MarshalIndent(x, "", "\t")
	return string(buf)
}

func (XMLFormatter) Prefix() string {
	return "<report>"
}

func (XMLFormatter) Suffix() string {
	return "</report>"
}

//Report allows to format and write multiple metrics for multiple website computed within a given timeframe.
//It uses the Formatter to format (Metric, Result)s and the Logger to write the final result.
func (r Reporter) Report(metrics []WebMetrics) {
	res := r.f.Prefix()
	for _, v := range metrics {
		res += r.f.Multiple(v)
	}
	res += r.f.Suffix()
	r.l.Printf("%v", res)
}

//DefaultLogger is a convenient function which returns a pointer to a basic Logger, with [MicroMon] prefix and local datetime prefix.
func DefaultLogger() *log.Logger {
	return log.New(os.Stdout, "\n[MicroMon] ", log.LstdFlags)
}

//FileLogger is a convenient function which returns a pointer to a logger which writes in a file.
//Path is given in parameter.
func FileLogger(path string) *log.Logger {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalf("Error while opening %s : %v", path, err)
	}
	return log.New(f, "", 0)
}

//Taken from : https://stackoverflow.com/a/21117347
//Used to format/pretty print => VERY BASIC <= XML (i.e. without attributes)
type node struct {
	Attr     []xml.Attr
	XMLName  xml.Name
	Children []node `xml:",any"`
	Text     string `xml:",chardata"`
}
