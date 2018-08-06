# Timestamp

`Metadata` is an implementation of druid dataSoruce metadata query.

Use case is to find most recent data ingestion timestamp per datasource.

The method `GenerateLastIngestionTimestamp(druidEndPoint string, dqRetries int, datasource string)`:

- POSTs the query (supplied in JSON format) to the Druid endpoint `druidEndPointt` for `datasource`.
- In case of error, retries up to the specified number of max retries `dqRetries`.
- Returns the time in epoch format, which will be used to calculate ingestion lag.

See the unit tests for example use.
