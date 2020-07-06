package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/bhupeshbhatia/website-monitoring-tool/database"
	"github.com/bhupeshbhatia/website-monitoring-tool/request"
)

// Site representes the entities we want to monitor
type Site struct {
	URL           string `json:"url"`
	CheckInterval int    `json:"checkInterval"`
}

// StartSiteMonitor starts a ticker for the given website
// it sends a request following a user-defined interval
func StartSiteMonitor(ctx context.Context, site Site, logc chan request.ResponseLog) error {
	ticker := time.NewTicker(time.Duration(site.CheckInterval) * time.Second)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return nil
		case t := <-ticker.C:
			log, err := request.Send(t, site.URL)
			if err != nil {
				return fmt.Errorf("error while monitoring %v:\n Details: %v", site, err)
			}
			logc <- log
		}
	}
}

// ProcessLogs reads logs from the log channel and processes them: for now write to influxDB
func ProcessLogs(ctx context.Context, logc chan request.ResponseLog) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case log := <-logc:
			if err := database.WriteLogToDB(log); err != nil {
				return fmt.Errorf("error while processing a log:\n %v", err)
			}
		}
	}
}
