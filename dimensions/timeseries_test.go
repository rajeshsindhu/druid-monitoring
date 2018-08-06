package dimensions

import (
	"encoding/json"
	"github.com/golang/glog"
	"gitlab.internal.unity3d.com/ads/data-eng/druid-monitoring/util"
	"testing"
)

func TestTimeseriesQuery(t *testing.T) {

	expected := `{"queryType":"timeseries","dataSource":"ds_test","intervals":"2018-03-14T00:00:00.000Z/2018-03-15T00:00:00.000Z","granularity":"day","filter":{"type":"and","fields":[{"type":"selector","dimension":"platform","value":""}]},"aggregations":[{"type":"count","name":"count"}]}`

	actual := generateTimeseriesQuery("ds_test", "2018-03-14T00:00:00.000Z", "2018-03-15T00:00:00.000Z", "day", "platform", "count", "")
	equals, err := util.JsonEquals(expected, actual)
	if err != nil {
		t.Errorf("Unable to unmarshal json blob: %s", err)
	}
	if !equals {
		t.Error("Expected:")
		t.Error(expected)
		t.Error("\nGot:")
		t.Error(actual)
	}
}

func TestTimeseriesResponse(t *testing.T) {

	testResponse := `[{"timestamp":"2018-03-19T00:00:00.000Z","result":{"count":43382789}},{"timestamp":"2018-03-20T00:00:00.000Z","result":{"count":24792899}}]`
	values := make(map[string]int64)

	var rsp []interface{}
	err := json.Unmarshal([]byte(testResponse), &rsp)
	if err != nil {
		glog.Info("empty map :%s", make([]interface{}, 0))
	}

	for _, item := range rsp {

		itemMap := item.(map[string]interface{})
		timestamp := itemMap["timestamp"]
		if event, ok := itemMap["result"]; ok {
			eventMap := event.(map[string]interface{})
			if value, ok := eventMap["count"]; ok {
				values[timestamp.(string)] = int64(value.(float64))
			}
		}
	}

	glog.Info("timeseries map:%s", values)

}
