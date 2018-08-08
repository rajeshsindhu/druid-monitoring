FROM golang:1.9.1
RUN mkdir -p /go/src/github.com/rajeshsindhu/druid-monitoring
ADD . /go/src/github.com/rajeshsindhu/druid-monitoring/
WORKDIR /go/src/github.com/rajeshsindhu/druid-monitoring
RUN go build .

# We don't need to expose any ports since this service makes outgoing connections to druid cluster and statsD.

# When running this image in a container, the following environment variables must be set:
# druidEndpoint : broker node address
# dataLagDs    : druid ds
# nullDimDs    : druid ds



# Optional parameters
# dataLagMonitoringInterval : monitoring stats notification interval for data ingestion lag in minutes, default value : 1
# dimMonitoringInterval : monitoring stats notification interval for null dimensions in minutes, default value : 60
# dqRetries: druid query retries count ,default vale : 3



