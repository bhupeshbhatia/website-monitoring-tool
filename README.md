# Website monitoring tool

### Requirements

- [InfluxDB 2.0](https://www.influxdata.com/) - open source time series database
- [Go 1.14](https://golang.org/) - a systems programming language
- [Docker]() - we will use to ease up running an InfluxDB instance


_Overview_

- A console program to monitor the performance and availability of websites
- Websites and check intervals are user-defined

_Statistics - TBD_

- Check the different websites with their corresponding check intervals
- Compute a few interesting metrics: availability, max/avg response times, max/avg time to first byte, response codes count

_Alerting - TBD_

- When a website availability is below a user-defined threshold for a user-defined interval, an alert message is created: "Website {website} is down. availability={availability}, time={time}" (default config threshold: 80%, interval: 2min)
- When availability resumes, another message is created detailing when the alert recovered

_Dashboard - Might change_

- displays stats for a user-defined timeframe, stats are updated following a user-defined interval. Default:
  - Every 10s, display the stats for the past 10 minutes for each website
  - Every minute displays the stats for the past hour for each website
- Show all past alerting messages


### Installation

#### Building from source

Run an InfluxDB instance:

```sh
$ docker run -p 8086:8086 -v influxdb:/var/lib/influxdb influxdb
```


### Implementation

The project uses built-in concurrency features of Go. The idea is to learn more about concurrency and as such the entities are run concurrently using goroutines. All communications are done through go channels, particularly the monitor logs and alerts channel.

The error management and propagation are done using "errgroup" package that facilitates managing errors while spanning multiple goroutines.

**Monitor**

The monitor starts concurrent tickers linked to each website. Based on a user-defined interval, the request is sent to the website which will measures metrics(response time, time to first byte), and sends the results as a measurement to our logs channel.

**Database**

Stores measurements in a time-based manner. It facilitates getting measurements for a particular timeframe.

**Statsagent**

Called by other entities. It computes the stats(avg/max response time, avg/max time to first byte) for the websites we monitor. It also computes the availability of a website over a timeframe.

**Dashboard**

Displays stats about the websites we monitor with user-defined configs(update interval, stats timeframe). It starts concurrent tickers for each view that call stats agent to get the new metrics.

The dashboard also listens to the alerts channel and displays new and past alerts on the GUI.

**Alerting**

It starts a ticker with a user-defined interval that calls the stats agent to compute the availability for a user-defined timeframe. All alerts are sent to an alerts channel that is consumed by our dashboard. The alerting ticker interval will be kept small to keep accuracy, but we don't want to overload the database. 

This can be improved if we rely on a pub/sub approach to reduce the overload, which InfluxDB supports.


### Testing

The tests willl be provided for the alerting process. The Go Testing package will be used for this purpose.

The scenarios are as follows: 

    No records: no "website down" alert

    Not enough records on the last timeframe, availability <= threshold: No "website down" alert.

    If website state is up and the last timeframe, availability > threshold: "website up" alert.

    If website state is up and the last timeframe, availability <= threshold: send a "website down" alert.

    If website state is down and the last timeframe, availability <= threshold: No "website down" alert.

    If website state is down and the last timeframe, availability > threshold: send a "website up" alert.