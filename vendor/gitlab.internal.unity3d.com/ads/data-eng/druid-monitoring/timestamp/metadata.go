package timestamp

import (
	"encoding/json"
	"github.com/golang/glog"
	"gitlab.internal.unity3d.com/ads/data-eng/druid-monitoring/util"
	"time"
)

const timeFormat = "2006-01-02T15:04:05.999Z"

type druidResponse struct {
	Timestamp            string `json:"timestamp"`
	MaxIngestedEventTime string `json:"result.maxIngestedEventTime"`
}

type sourceMetadataQuery struct {
	QueryType  string `json:"queryType"`
	DataSource string `json:"dataSource"`
}

func generateSourceMetadataQuery(datasource string) string {
	datasourceMetadataQuery := sourceMetadataQuery{"dataSourceMetadata", datasource}
	qs, _ := json.Marshal(datasourceMetadataQuery)
	return string(qs)
}

func GenerateLastIngestionTimestamp(druidEndPoint string, dqRetries int, datasource string) (int64, error) {

	druidOutput, err := util.ExecutePostQuery(druidEndPoint, dqRetries, generateSourceMetadataQuery(datasource))
	if err != nil {
		glog.Infof("Received Error ,executing druid sourceMetadataQuery: %s", err.Error())
	}

	return parseTimestamp(druidOutput)
}

func parseTimestamp(body []byte) (int64, error) {

	responseArray := make([]druidResponse, 0)
	resultErr := json.Unmarshal(body, &responseArray)
	if resultErr != nil {
		glog.Errorf("Received error unmarshalling response body : %s  ", resultErr.Error())
		return 0, resultErr
	}

	glog.Infof("Druid sourceMetadataQuery output : %v  ", responseArray)

	//len(responseArray)==0 implies empty response,parsing empty response will cause array index out of bounds exception.
	if len(responseArray) > 0 {
		timestamp, err := time.Parse(timeFormat, responseArray[0].Timestamp)

		if err == nil {
			return timestamp.Unix(), nil
		} else {
			resultErr = err
		}
	} else {
		glog.Infof("Empty Druid response : %v  ", responseArray)
	}
	// This would return 0 instead of panic, if datasources is wrong or empty. because Panic would cause service to exit and never re-run again.
	// returning latest ingestion timestamp as 0 will also indicate an error on Datadog because in normal scenario this is not possible.
	return 0, resultErr

}
