package statsagent

import (
	"time"

	"github.com/bhupeshbhatia/website-monitoring-tool/database"
	"github.com/bhupeshbhatia/website-monitoring-tool/request"
)

// WebsiteStats contains useful metrics about website
type WebsiteStats struct {
	StatusCodeCount    map[string]int
	AvgResponseTime    time.Duration
	MaxResponseTime    time.Duration
	AvgTimeToFirstByte time.Duration
	MaxTimeToFirstByte time.Duration
	Availability       float64
}

// GetStats of provided websites for a particular timeframe
func GetStats(urls []string, origin time.Time, timeframe int64) map[string]WebsiteStats {
	websitesStats := make(map[string]WebsiteStats)

	for _, url := range urls {
		v := database.GetRecordsForURL(url, origin, timeframe)
		statusCodeCount := make(map[string]int)
		var sumResponseTime int64 = 0
		var maxResponseTime time.Duration = 0
		var avgResponseTime float64 = 0
		var sumTimeToFirstByte int64 = 0
		var maxTimeToFirstByte time.Duration = 0
		var avgTimeToFirstByte float64 = 0
		var successCount float64 = 0
		var availability float64 = 0

		for _, line := range v {

			//statuscode count
			if _, ok := statusCodeCount[line.StatusCode]; ok {
				statusCodeCount[line.StatusCode]++
			} else {
				statusCodeCount[line.StatusCode] = 1
			}

			//For successes - sum of response time, sumtofirstbyte
			if line.Success {
				//increase successcount

				//condition for higher loadtime

				//condition for timetofirstbyte being greater than maxtimeforfirstbyte

				//add all loadtimes
				//add all timetofirstbyte
			}
		}
		if successCount > 0 {
			//what is avgresponse time

			//what is avg time to first byte

			//what is the availability ---? How to calculate this = success over numrecords?
		}

		//Website stats
		websitesStats[url] = WebsiteStats{StatusCodeCount: statusCodeCount, AvgResponseTime: time.Duration(avgResponseTime), MaxResponseTime: maxResponseTime, AvgTimeToFirstByte: time.Duration(avgTimeToFirstByte), MaxTimeToFirstByte: maxTimeToFirstByte, Availability: availability}
	}
	return websitesStats
}

// AvailabilityRange struct for computing availability of website
type AvailabilityRange struct {
	Availability float64
	Start        time.Time
}

// GetAvailabilityForTimeFrame computes the availability of a Website
// given a time origin and a timeframe
func GetAvailabilityForTimeFrame(url string, origin time.Time, timeframe int64) AvailabilityRange {
	records := database.GetRecordsForURL(url, origin, timeframe)
	return GetAvailabilityForRecords(records, origin)
}

// GetAvailabilityForRecords returns the availability given a slice of records
func GetAvailabilityForRecords(records []request.ResponseLog, origin time.Time) AvailabilityRange {
	var start time.Time = origin
	var successCount float64 = 0
	var availability float64 = 0

	//for successes - increase number
	for _, line := range records {
		if line.Success {
			successCount++
		}
	}

	//When successcount is greater than 0? Calculate availability - what is it? success/number of records? Like getstats?
	if successCount > 0 {

	}

	//return struct
	return AvailabilityRange{Availability: availability, Start: start}
}
