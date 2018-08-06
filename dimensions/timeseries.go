package dimensions

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"gitlab.internal.unity3d.com/ads/data-eng/druid-monitoring/util"
)

func generateTimeseriesQuery(datasource string, startTime string, endTime string, granularity string, dimension string, filterKey string, filterValue string) string {
	type aggregation struct {
		Type string `json:"type"`
		Name string `json:"name"`
	}

	type fields struct {
		Type      string `json:"type"`
		Dimension string `json:"dimension"`
		Value     string `json:"value"`
	}

	type filter struct {
		Type   string   `json:"type"`
		Fields []fields `json:"fields"`
	}

	type druidRequest struct {
		QueryType    string        `json:"queryType"`
		DataSource   string        `json:"dataSource"`
		Intervals    string        `json:"intervals"`
		Granularity  string        `json:"granularity"`
		Filter       filter        `json:"filter"`
		Aggregations []aggregation `json:"aggregations"`
	}

	request := druidRequest{
		"timeseries",
		datasource,
		startTime + "/" + endTime,
		granularity,
		filter{"and", []fields{{"selector", dimension, filterValue}}},
		[]aggregation{{"count", filterKey}},
	}

	jsonReq, _ := json.Marshal(request)
	glog.Infof("timeseries query: %s for dataSource: %s ", string(jsonReq), datasource)
	return string(jsonReq)
}

func GetActiveDimensions(druidEndPoint string, dataSource string, dqRetries int) ([]string, error) {

	dimensionsQuery := fmt.Sprintf("/%s/%s/%s", "datasources", dataSource, "dimensions")

	dimensionResp, err := util.ExecuteGetQuery(druidEndPoint, dqRetries, dimensionsQuery)

	// This would return empty slice instead of panic, if datasources is wrong or empty. because Panic would cause service to exit and never re-run again.
	if err != nil {
		glog.Errorf("Received Error, executing Get query: %s for dataSource : %s", err.Error(), dataSource)
		return make([]string, 0), err
	}

	var dimResp []string
	err = json.Unmarshal([]byte(dimensionResp), &dimResp)

	return dimResp, err
}

func GetDimensionCounts(druidEndPoint string, dqRetries int, dataSource string, granularity string, dim string, filterKey string, filterValue string, latestTime int64) (map[string]int64, error) {

	startTime, endTime := util.GetTimeRange(latestTime, granularity)

	req := generateTimeseriesQuery(dataSource, startTime, endTime, granularity, dim, filterKey, filterValue)
	jsonRsp, err := util.ExecutePostQuery(druidEndPoint, dqRetries, req)

	glog.Infof("timeseries query response : %s", string(jsonRsp))

	// This would return empty slice instead of panic, if datasources is wrong or empty. because Panic would cause service to exit and never re-run again.
	if err != nil {
		glog.Errorf("Received Error, executing Post query: %s for dataSource: %s", err.Error(), dataSource)
		return make(map[string]int64, 0), err
	}

	timeCountMap := make(map[string]int64)

	var rsp []interface{}
	err = json.Unmarshal([]byte(jsonRsp), &rsp)

	// This would return empty slice instead of panic, if datasources is wrong or empty. because Panic would cause service to exit and never re-run again.
	if err != nil {
		glog.Errorf("Received Error, unmarshalling Post query response: %s for dataSource: %s", err.Error(), dataSource)
		return make(map[string]int64, 0), err
	}

	for _, item := range rsp {

		itemMap := item.(map[string]interface{})
		timestamp := itemMap["timestamp"]
		if event, ok := itemMap["result"]; ok {
			eventMap := event.(map[string]interface{})
			if value, ok := eventMap[filterKey]; ok {
				timeCountMap[timestamp.(string)] = int64(value.(float64))
			}
		}
	}

	return timeCountMap, err

}
