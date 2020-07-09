package database

import (
	"fmt"
	"log"
	"net/url"
	"time"

	str2duration "github.com/xhit/go-str2duration"

	"github.com/bhupeshbhatia/website-monitoring-tool/request"
	"github.com/influxdata/influxdb/client/v2"
)

const (
	layout string = "2006-01-02T15:04:05.000Z"
)

type InfluxDb struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	DatabaseName string `json:"databaseName"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

var (
	influxDBcon client.Client
)

// Initialize influx db
func (influxDb InfluxDb) Initialize() error {
	u, err := url.Parse(fmt.Sprintf("http://%s:%d", influxDb.Host, influxDb.Port))
	if err != nil {
		println("InfluxDB : Invalid Url,Please check domain name given in config file \nError Details: ", err.Error())
		return err
	}

	conf := client.HTTPConfig{
		Addr:     u.String(),
		Username: influxDb.Username,
		Password: influxDb.Password,
	}

	influxDBcon, err = client.NewHTTPClient(conf)
	if err != nil {
		return fmt.Errorf("error creating a client for InfluxDB: %v", err)
	}

	_, _, err = influxDBcon.Ping(10 * time.Second)
	if err != nil {
		return fmt.Errorf("error while trying to ping InfluxDB: %v", err)
	}

	createDbErr := createDatabase(influxDb.DatabaseName)

	if createDbErr != nil {
		if createDbErr.Error() != "database already exists" {
			return fmt.Errorf("failed to create InfluxDB Database instance: %v", createDbErr)
		}

	}

	return nil
}

func createDatabase(databaseName string) error {

	_, err := queryDB(fmt.Sprintf("create database %s", databaseName), "")

	return err
}

func queryDB(cmd string, databaseName string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: databaseName,
	}
	if response, err := influxDBcon.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	}
	return
}

// AddRecord adds a new record to InfluxDB
func (influxDb InfluxDb) AddRecord(responseLog request.ResponseLog) error {

	tags := map[string]string{
		"requestId": responseLog.URL,
	}
	fields := map[string]interface{}{
		"responseTime":    responseLog.LoadTime,
		"timeToFirstByte": responseLog.TTFB,
		"StatusCode":      responseLog.StatusCode,
		"Success":         responseLog.Success,
	}

	bps, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  influxDb.DatabaseName,
		Precision: "ms",
	})

	if err != nil {
		return err
	}

	point, err := client.NewPoint(
		responseLog.URL,
		tags,
		fields,
		responseLog.Timestamp,
	)
	if err != nil {
		return err
	}

	bps.AddPoint(point)

	err = influxDBcon.Write(bps)
	if err != nil {
		return err
	}

	return nil
}

// GetRecordsForURL sends a query to InfluxDB
// to get records of a given URL, older than a given "origin" and restricted by a given timeframe
func (influxDb InfluxDb) GetRecordsForURL(url string, origin time.Time, timeframe int64) []request.ResponseLog {
	q := fmt.Sprintf(`select * from "%s" WHERE time >= '%v' - %dm`, url, origin.Format(time.RFC3339), timeframe/60)
	res, err := queryDB(q, influxDb.DatabaseName)
	if err != nil {
		log.Printf("%v", err)
	}

	//convert string to duration -----cool library
	s2dParser := str2duration.NewStr2DurationParser()
	fmt.Println(s2dParser)

	records := make([]request.ResponseLog, 0)
	//Need to use standard duration parser - to parse timestamp from response from DB
	//use NewStr2DurationParser to convert response time and time to first byte
	for _, result := range res {
		if len(result.Series) == 0 {
			continue
		}
		for _, val := range result.Series[0].Values {
			timestamp, _ := time.Parse(layout, val[0].(string))
			statusCode, _ := val[1].(string)
			success := val[2].(bool)
			url := val[3].(string)
			responseTime, _ := s2dParser.Str2Duration(val[4].(string))
			timeToFirstByte, _ := s2dParser.Str2Duration(val[5].(string))
			item := request.ResponseLog{Timestamp: timestamp, StatusCode: statusCode, URL: url, TTFB: timeToFirstByte, LoadTime: responseTime, Success: success}
			records = append(records, item)
		}
	}
	return records
}
