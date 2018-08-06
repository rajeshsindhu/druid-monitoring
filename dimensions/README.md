# Dimensions

`timeseries` is an implementation of druid timeseries query.

Use case is to find null dimensions count per datasource.

The method `GetActiveDimensions(druidEndPoint string, dataSource string, dqRetries int)`:

- GET query to the Druid endpoint `druidEndPointt` for `datasource`.
- In case of error, retries up to the specified number of max retries `dqRetries`.
- Returns list of all the active dimensions, which will be used to calculate null dim count.


The Method `GetDimensionCounts(druidEndPoint string, dqRetries int, datasource string, granularity string, dim string, filterKey string, filterValue string, latestTime int64)`:

- POSTs the query (supplied in JSON format) to the Druid endpoint `druidEndPointt` for `datasource`.
- In case of error, retries up to the specified number of max retries `dqRetries`.
- `granularity` timeseries Query granurality
- `dim` datasource dimension
-  `filterKey` and `filterValue` used in timeseries query, defualt value is `count ` (number rows in druid).
- Returns map of `timewindow -> count`.
- `latestTime` is calculated from sourceMetadata query to make sure data is present in druid.


See the unit tests for example use.
