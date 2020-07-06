package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bhupeshbhatia/website-monitoring-tool/alerting"
	"github.com/bhupeshbhatia/website-monitoring-tool/database"
	"github.com/bhupeshbhatia/website-monitoring-tool/monitor"
	"github.com/bhupeshbhatia/website-monitoring-tool/request"
	"golang.org/x/sync/errgroup"
)

// Config struct containing websites config(url, check interval), database data(host, dbaname, username, password)
type Config struct {
	Websites []monitor.Site       `json:"websites"`
	Database database.Type        `json:"database"`
	Alert    alerting.AlertConfig `json:"alerting"`
}

func getConfig() (Config, error) {
	configFile, err := os.Open("env-config.json")
	if err != nil {
		return Config{}, fmt.Errorf("%v", err)
	}
	defer configFile.Close()

	configByteContent, err := ioutil.ReadAll(configFile)
	if err != nil {
		return Config{}, err
	}

	var config Config

	if err := json.Unmarshal(configByteContent, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func main() {
	config, err := getConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	if err := database.Set(config.Database); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up the database: %v\n", err)
		os.Exit(1)
	}

	websiteList := []string{}
	websiteMap := make(map[string]int64)
	for _, ws := range config.Websites {
		websiteList = append(websiteList, ws.URL)
		websiteMap[ws.URL] = int64(ws.CheckInterval)
	}

	ctx, done := context.WithCancel(context.Background())

	//Using Errgroup: Package errgroup provides synchronization, error propagation,
	//and Context cancelation for groups of goroutines working on subtasks of a common task.
	g, gctx := errgroup.WithContext(ctx) //The derived Context is cancelled the first time a function passed to Go returns a non-nil error or the first time Wait returns, whichever occurs first.

	// Start goroutines to ping websites

	// Channel for logs
	logc := make(chan request.ResponseLog)

	//Channel for alerts
	alertc := make(chan string)

	//Defer to close them
	defer close(logc)
	defer close(alertc)

	//Call the goroutines
	//go dashboard
	go alerting.Run(alertc, websiteMap, config.Alert)

	g.Go(func() error {
		return monitor.ProcessLogs(gctx, logc)
	})

	for _, ws := range config.Websites {
		ws := ws
		g.Go(func() error {
			return monitor.StartSiteMonitor(gctx, ws, logc)
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

}
