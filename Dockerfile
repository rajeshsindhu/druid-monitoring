FROM golang:1.9.1
RUN mkdir -p /go/src/gitlab.internal.unity3d.com/ads/data-eng/druid-monitoring
ADD . /go/src/gitlab.internal.unity3d.com/ads/data-eng/druid-monitoring/
WORKDIR /go/src/gitlab.internal.unity3d.com/ads/data-eng/druid-monitoring
RUN go build .

# We don't need to expose any ports since this service makes outgoing connections to druid cluster and statsD.

# When running this image in a container, the following environment variables must be set:
# druidEndpoint : http://druid-a-broker-1.us-east-1.applifier.info:8080
# dataLagDs    : iap_verified_transaction_events_rt,ads_events_rt
# nullDimDs    : ads_events



# Optional parameters
# dataLagMonitoringInterval : monitoring stats notification interval for data ingestion lag in minutes, default value : 1
# dimMonitoringInterval : monitoring stats notification interval for null dimensions in minutes, default value : 60
# dqRetries: druid query retries count ,default vale : 3



