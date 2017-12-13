/*
Package micromon offers types and methods to monitor websites. Monitoring means, here, to compute and report a set of metrics for multiple websites.

Each website is checked at a user-defined interval, and periodically, metrics are computed over different timeframes, then reported. Special methods, called hooks,
can operate on this metrics to implement alerting logic for example, i.e. keeping traces of unavailabilities and recoveries of websites, and warning the user.

This package is essentially meant to be used with its entry-point method Start(), which implements a complete logic, but can also be used partially via
exposed methods. The later can be called without any context which would have been set by the Start() method, as shown in the Examples section.

We briefly explain the main concepts below.

Types

MicroMon defines five main types : Config, MetaResponse, Metric, Hook and Reporter, along with secondary types.
Config holds the application configuration. MetaResponse holds the meta-informations from a website HTTP response.
Metric is an interface which defines a method to aggregate MetaResponses, and Hook is a closure meant to do extra work on metrics.
Reporter allow to customize the formatting and the writing of computed metrics.

Metrics

Metrics are generic types which aggregate multiple MetaResponse into a generic Result. A Metric just needs to have
an aggregation method, a name and a description. The computed metric is described by a Result, which just need to
implement a formatting method. Classic Result are integers, floats or maps of Result.

New Metrics may be created by implementing a new type, as well as Results.

Reporter

To provide flexibility, the type Reporter has just a Report() method, which reports a set of metrics associated with websites and timeframes.
Reporter are just a Formatter, which needs to implement formatting methods, and a log.Logger, which specifies the place to write formatted metrics.

Hooks

Hooks are a way to do extra work on computed metrics without causing side-effects and without impacting reporting. Typically, alerting logic
is a hook. Hook are implemented as closures which operate on metrics and are meant to be called after metrics have been computed
and before metrics have been reported. Any extra work/logging should be implemented as a hook.

Usage

The method Start(path) is provided to handle all the monitoring logic. It takes a path to a YAML configuration file, which defines
websites to watch, metrics to computes, hooks to call and reporter to use, along with other parameters. This is the most easy way
to monitor websites.

But it is possible to do it without a configuration file by using the methods exposed by the package, as showed in examples section.
See app.go for a complete example of monitoring.

Examples

This example shows how to get a single MetaResponse from a website.

	//Timeout is 1 second
	resp, err := CheckUrl("https://github.com", 1)

This example shows how to collect responses from a website without computing any metric.

	//Websites to watch and check intervals
	webs := make(map[string]Website)
	webs["github"] = Website{URL:"https://github.com"}
	conf := Config{Websites: webs, Timeout: 1}

	//Get MetaResponse channel
	ch := WatchWebsites(conf)

	//Gather incoming MetaResponse
	data := NewSafeData()
	for {
		resp := <-ch
		data.Mux.Lock()
		data.Datas = append(data.Datas, resp)
		data.Mux.Unlock()
	}

The following example shows how to compute and print a specific metric given a slice of MetaResponse.

	//Suppose that resp is fed with MetaResponse
	var resp []MetaResponse

	met := AvgRespTime{}
	res := met.Compute(resp)
	fmt.Printf("%s : %s", met.Description(), res.Format(false))

This final example shows how to compute and report multiple metrics on multiple websites for a given timeframe.
In this example we write metrics formatted with XML in a log file.

	//Structure for 5 websites. Suppose it is fed concurrently with MetaResponse.
	var data := NewRespMap(5)

	//Build reporter
	rep := NewReporter(FileLogger("output.log"), XMLFormatter{})

	//Metrics to compute
	met := []Metric{AvgRespTime{}, Availability{}}

	//Compute metrics for the last two minutes
	res := (&data).ComputeMetrics(met, 2)

	//Write XML results in log file
	rep.Report(res)
*/
package micromon
