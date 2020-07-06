package alerting

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

var red *color.Color = color.New(color.FgRed)
var green *color.Color = color.New(color.FgGreen)

// websiteUp is a map that keeps track of the state of each website we're monitoring
var websiteUp map[string]bool = make(map[string]bool)

// AlertConfig represents all the useful info for our alert logic
type AlertConfig struct {
	AvailabilityInterval  int64   `json:"availabilityInterval"`
	AvailabilityThreshold float64 `json:"availabilityThreshold"`
	CheckInterval         int     `json:"checkInterval"`
}

// Run monitors the availability of websites
// It send an alert to the dashboard, if the availability of some website over a given interval
// is below the given threshold
func Run(alertc chan string, websitesMap map[string]int64, alertConfig AlertConfig) {
	urls := make([]string, 0)
	for k := range websitesMap {
		websiteUp[k] = true
		urls = append(urls, k)
	}
	ticker := time.NewTicker(time.Duration(alertConfig.CheckInterval) * time.Second)

	for {
		select {
		case t := <-ticker.C:
			for _, url := range urls {
				//Statsagent - getavailabilityfortimeframe
				//alertmessage function
				fmt.Println(t, url)
			}
		}
	}
}

// need to look into this as well.
//getAlertMessage function
