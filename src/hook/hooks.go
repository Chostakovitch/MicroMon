//hook contains types and methods to construct hooks, which are extra work on a set of metrics.
//hooks are meant to be call before doing reporting and after metrics has been calculated without causing side effects.
//hook package is tightly coupled with the rest of the application, which is kind of bad.
package hook

import (
	"config"
	"fmt"
	"log"
	"metric"
	"strconv"
	"time"
)

//Hooker is an interface meant to describe methods which return hook.
//Hooks are closures which operates on a set of Metric/Result associated with websites name.
//Typically, hook are called each time metrics are computed without causing any side-effects and they are user-configured.
//Hooks are a way to do extra work on Metric which should not be part of classic metrics reporting.
//Simple examples are extra logging, alerting logic, etc.
type Hooker interface {
	//GetHook returns a hook and takes the application configuration in parameter.
	GetHook(config.Config) Hook
}

//A Hook is just a procedure which operates on a set of WebMetrics.
//It returns an arbitrary string for debugging and testing purpose.
type Hook func([]metric.WebMetrics) string

//GetHook takes the name of a hook and a Config and returns the associated hook.
//If not hook corresponding to name is found, a non-nil error is returned.
func GetHook(name string, conf config.Config) (Hook, error) {
	switch name {
	case "alert":
		return AlertHook{}.GetHook(conf), nil
	}
	return nil, fmt.Errorf("%s is not a known hook name", name)
}

//AlertHook is an empty struct which implements Hook.
//It provides a hook which manages the logic alert when websites availability is behind a threshold.
//The hook, being a closure, can keep trace of previous alerts and keep them on screen.
//Alerts are printed in standard output and do not use classic Reporter struct.
type AlertHook struct{}

//Internal struct to hold information about a website unavailability
type webDown struct {
	name          string
	when          time.Time
	avail         float64
	recovered     bool
	whenRecovered time.Time
}

//addUnavailability is called when a website is unavailable. It records a new unavailability if
//no old unavailability concerning this website is still unrecovered.
//It returns a boolean which indicates if a new unavailability has been recorded and a slice describing all previous unavailabilities.
func addUnavailability(s []webDown, name string, avail float64, when time.Time) ([]webDown, bool) {
	for i := len(s) - 1; i >= 0; i-- {
		//There is still an alert for this website, don't add one
		if s[i].name == name && !s[i].recovered {
			return s, false
		}
	}
	//Here, there is no alert, add one
	return append(s, webDown{name, when, avail, false, time.Time{}}), true
}

//recoverAvailability is called when a website is available. If a previous unavailability concerning this website has
//not been recovered yet, it recovers it. It returns a boolean which indicates if a unavailability has been effectively
//recovered, along with a slice describing all previous unavailabilities.
func recoverAvailability(s []webDown, name string, when time.Time) ([]webDown, bool) {
	for i := len(s) - 1; i >= 0; i-- {
		//This website was effectively not available ; recover it
		if s[i].name == name && !s[i].recovered {
			s[i].recovered = true
			s[i].whenRecovered = when
			return s, true
		}
	}
	return s, false
}

func (AlertHook) GetHook(conf config.Config) Hook {
	threshold := conf.AvailThreshold
	memories := make([]webDown, 0)
	return func(metrics []metric.WebMetrics) string {
		now := time.Now()
		res := ""
		effect := false
		//For each website, check if it just became available OR unavailable
		for _, s := range metrics {
			for _, m := range s.M {
				if _, ok := m.M.(metric.Availability); ok {
					if avail, ok := m.R.(metric.MetricFloat); ok && avail < metric.MetricFloat(threshold) {
						memories, effect = addUnavailability(memories, s.N, float64(avail), now)
						if effect {
							res = "unavailable"
						}
					} else {
						memories, effect = recoverAvailability(memories, s.N, now)
						if effect {
							res = "recovered"
						}
					}
				}
			}
		}

		//Print all availabilities
		for _, m := range memories {
			msg := fmt.Sprintf("Website %v is down. Availability = %v%%, time = %v\n", m.name, strconv.FormatFloat(float64(m.avail), 'f', 3, 64), m.when.Format("2006/02/01 15:04:05"))
			if m.recovered {
				msg += fmt.Sprintf("\tRecovered. Time = %v", m.whenRecovered.Format("2006/02/01 15:04:05"))
			}
			log.Printf("==== AVAILABILITY ALERTS ====\n%v\n\n", msg)
		}

		return res
	}
}
