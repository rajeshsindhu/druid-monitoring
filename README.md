# Druid-monitoring
Service to pull Druid metadata and send it to Statsd for monitoring purposes.

Use case:
1. To track latest data ingestion timestamp for datasources "iap_verified_transaction_events_rt", "ads_events_rt" and "iap_events_promo".

   `-dataLagDs` argument takes comma separated names of druid datasources as input for data ingestion lag. new datasource can be tracked by appending to this input.

2. To track null dimensions count per datasource "ads_events".

   `-nullDimDs` argument takes comma separated (name + query granurality) of druid datasources as input for data ingestion lag. new datasource can be tracked by appending to this input.


# How to run 

`go run server.go -druidEndPoint http://druid-a-broker-1.us-east-1.applifier.info:8080 -dataLagDs iap_verified_transaction_events_rt,ads_events_rt,iap_events_promo -nullDimDs ads_events/day`

Optional parameters

 - dataLagMonitoringInterval : monitoring stats notification interval for data ingestion lag in minutes, default value : 1
 - dimMonitoringInterval : monitoring stats notification interval for null dimensions in minutes, default value : 60
 - dqRetries: druid query retries count ,default vale : 3
 - filterKey: druid metric key for e.g. "revenue","count" etc.
 - filterValue: filter value, "" signifies null in druid.


### Logging

* stderr. We use [glog](https://github.com/golang/glog)  library to emit logs. Append the following command line
arguments if you want to see the log in stderr `-logtostderr -v 9`

* Kibana. TBD

 