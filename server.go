package main

import (
	"github.com/rajeshsindhu/druid-monitoring/dimensions"
	"github.com/rajeshsindhu/druid-monitoring/timestamp"
	"flag"
	"fmt"
	"github.com/DataDog/datadog-go/statsd"
	"github.com/caarlos0/env"
	"github.com/golang/glog"
	"github.com/jasonlvhit/gocron"
	"strings"
	"time"
)

var (
	druidEndPoint             string
	dataLagDs                 string
	nullDimDs                 string
	filterKey                 string
	filterValue               string
	statsdEndPoint            string
	dataLagMonitoringInterval uint64
	dimMonitoringInterval     uint64
	dqRetries                 int
)

type OsEnvVariableConfig struct {
	StatsdHost string `env:"STATSD_HOST" envDefault:"localhost"`
	StatsdPort string `env:"STATSD_PORT" envDefault:"8125"`
}

func init() {

	cfg := OsEnvVariableConfig{}
	err := env.Parse(&cfg)
	if err != nil {
		glog.Errorf("Error parsing os env variables: %+v\n", err)
	}
	statsdEndPoint = fmt.Sprintf("%v:%v", cfg.StatsdHost, cfg.StatsdPort)

	flag.StringVar(&druidEndPoint, "druidEndPoint", "http://druid-a-broker-1.us-east-1.applifier.info:8080", "druid druidEndPoint")
	flag.StringVar(&dataLagDs, "dataLagDs", "iap_verified_transaction_events_rt,ads_events_rt,iap_events_promo", "datasources to pull ingestion lag data from druid")
	flag.StringVar(&nullDimDs, "nullDimDs", "ads_events/day", "(datasource/queryGranularity) to pull null dimensions from druid")
	flag.StringVar(&filterKey, "filterKey", "count", "timeseries query filter key for null dimensions stats")
	// empty string stands for null in druid
	flag.StringVar(&filterValue, "filterValue", "", "timeseries query filter value for null dimensions stats")
	flag.Uint64Var(&dataLagMonitoringInterval, "dataLagMonitoringInterval", 60, "monitoring interval in minutes")
	flag.Uint64Var(&dimMonitoringInterval, "dimMonitoringInterval", 60, "dimensions count monitoring interval in minutes")
	flag.IntVar(&dqRetries, "dqRetries", 3, "druid sourceMetadataQuery dqRetries count")
	flag.Parse()

	glog.Infof("druid-monitoring calling StatsD on: %s ", statsdEndPoint)

}

func statsDclient() *statsd.Client {

	statsdClient, err := statsd.New(statsdEndPoint)
	statsdClient.Namespace = "dm."

	if err != nil {
		glog.Errorf("Failed to connect to statsd at %s, Error %v", statsdEndPoint, err)
	}
	return statsdClient
}

func sendLatestIngestionTimeToStatsD() {

	statsdClient := statsDclient()

	statsdClient.Tags = append(statsdClient.Tags, "druidDataLag")

	defer statsdClient.Close()

	datasourcesList := strings.Split(strings.TrimSpace(dataLagDs), ",")

	for _, datasource := range datasourcesList {

		glog.Infof("Monitoring datasource for data ingestion lag : " + datasource)

		latestTimestamp, err := timestamp.GenerateLastIngestionTimestamp(druidEndPoint, dqRetries, datasource)

		if err != nil {
			glog.Errorf("Received Error ,processing dataSource metadata response: %s", err.Error())
		}

		dataIngestionLag := (time.Now().Unix()) - latestTimestamp

		glog.Infof("data ingestion lag in seconds : %d , for datasource : %s", dataIngestionLag, datasource)
		statsDKeyIngestionLag := fmt.Sprintf("%s.%s", datasource, "ingestionLag.seconds")
		statsdClient.Gauge(statsDKeyIngestionLag, float64(dataIngestionLag), nil, 1)

	}

}

func sendDimensionsCountToStatsD() {

	statsdClient := statsDclient()
	statsdClient.Tags = append(statsdClient.Tags, "druidNullDim")

	defer statsdClient.Close()

	datasourcesList := strings.Split(strings.TrimSpace(nullDimDs), ",")

	for _, dsGranularityInput := range datasourcesList {

		dsGranularityArr := strings.Split(strings.TrimSpace(dsGranularityInput), "/")

		if dsGranularityArr != nil && len(dsGranularityArr) == 2 {
			dataSource := dsGranularityArr[0]
			queryGranularity := dsGranularityArr[1]

			glog.Infof("Monitoring dsGranularityInput for null dimensions dataSource: %s , queryGranularity: %s", dataSource, queryGranularity)

			output, err := dimensions.GetActiveDimensions(druidEndPoint, dataSource, dqRetries)

			// Do not want to fail monitoring service, since it is a cron job.
			if err != nil {
				glog.Errorf("Received Error, processing null dimensions response: %s", err.Error())
				output = make([]string, 0)
			}

			latestTime, _ := timestamp.GenerateLastIngestionTimestamp(druidEndPoint, dqRetries, dataSource)

			for _, dim := range output {
				timeSeriesOutput, _ := dimensions.GetDimensionCounts(druidEndPoint, dqRetries, dataSource, queryGranularity, dim, filterKey, filterValue, latestTime)
				for timeWindow, count := range timeSeriesOutput {
					monitoringTag := []string{dataSource, timeWindow}
					statsDKeyNullCount := fmt.Sprintf("%s.%s.%s", dataSource, dim, filterKey)
					statsdClient.Gauge(statsDKeyNullCount, float64(count), monitoringTag, 1.0)
				}
			}

		}
	}

}

func main() {
	scheduler := gocron.NewScheduler()
	scheduler.Every(dataLagMonitoringInterval).Minutes().Do(sendLatestIngestionTimeToStatsD)
	scheduler.Every(dimMonitoringInterval).Minutes().Do(sendDimensionsCountToStatsD)
	<-scheduler.Start()

}
