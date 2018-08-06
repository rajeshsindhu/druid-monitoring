package timestamp

import (
	"testing"
)

// generateSourceMetadataQuery() should:
// generate source metadata sourceMetadataQuery .
func TestDruidQuery(t *testing.T) {

	testResponse := `{"queryType":"dataSourceMetadata","dataSource":"ds_test"}`

	res := generateSourceMetadataQuery("ds_test")
	if res != testResponse {
		t.Errorf("Query() did not return the expected response %s; instead returned %s", testResponse, res)
	}
}

// Generate result should return output in epoch seconds.
func TestGenerateResult(t *testing.T) {

	testInput := []byte(`[
    {
        "timestamp": "2017-11-13T20:00:00.000Z",
        "result": {
            "maxIngestedEventTime": "2017-11-13T20:00:00.000Z"
        }
    }
]`)

	testResponse := int64(1510603200)
	res, _ := parseTimestamp(testInput)

	if res != testResponse {
		t.Errorf("Query() did not return the expected response %s; instead returned %s", testResponse, res)
	}

	testInputEmpty := []byte(`[]`)

	emptyResponse, _ := parseTimestamp(testInputEmpty)

	emptyTestResponse := int64(0)
	if emptyResponse != emptyTestResponse {
		t.Errorf("Query() did not return the expected response %s; instead returned %s", emptyTestResponse, emptyResponse)
	}

}
