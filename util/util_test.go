package util

import (
	"bytes"
	"fmt"
	"github.com/golang/glog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTimeRange(t *testing.T) {

	latestTime := int64(1510603200)

	testStartTime := `2017-11-12T12:00:00Z`

	testEndTime := `2017-11-13T12:00:00Z`

	startTime, endTime := GetTimeRange(latestTime, "day")

	if startTime == testStartTime {
		glog.Info("start time " + startTime)
	}

	if endTime == testEndTime {
		glog.Info("end time " + endTime)
	}

}

//  Query should retry the specified number of times and fail afterwards.
func TestDruidQuery_retries(t *testing.T) {
	testBody := `{"key1":"value1", "key2":"value2"}`
	testResponse := []byte(`{"response": "Done"}`)
	tests := []struct {
		errsToReturn  int
		retries       int
		expectSuccess bool
	}{
		{0, 3, true},  // errors = 0
		{2, 3, true},  // 0 < errors < retries
		{3, 3, true},  // errors == retries
		{4, 3, false}, // errors > retries > 0
	}

	for _, test := range tests {
		errs := 0
		errsPtr := &errs
		testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if *errsPtr < test.errsToReturn {
				*errsPtr += 1
				http.Error(writer, "Service Unavailable", 500)
			} else {
				fmt.Fprint(writer, testResponse)
			}
		}))
		defer testServer.Close()

		res, err := ExecutePostQuery("ds_test", 3, testBody)
		if err == nil && !test.expectSuccess {
			t.Errorf("Query() should have returned an error, but succeeded and returned %s (errsToReturn = %d; retries = %d)", res, test.errsToReturn, test.retries)
		} else if err == nil {

			if !(bytes.Equal(res, testResponse)) {
				t.Errorf("Query() did not return the expected response %s; instead returned %s  (errsToReturn = %d; retries = %d)", testResponse, res, test.errsToReturn, test.retries)
			}
		}
	}
}

func TestJsonEquals(t *testing.T) {
	tests := []struct {
		first       string
		second      string
		errExpected bool
		expectedEq  bool
	}{
		{`{"k1" : "v1"}`, `{"k1":"v1"}`, false, true},
		{`{"k1" : "v1"}`, `{"k1":"v1", "k2": "v2"}`, false, false},
		{`{"k1" : "v1", "k2":"v2"}`, `{"k1":"v2", "k2": "v1"}`, false, false},
		{`{"k2" : "v2", "k1":"v1"}`, `{"k1":"v1", "k2": "v2"}`, false, true},
		{`{"k2" : 5, "k1":["foo",   "bar"]}`, `{"k1":["foo","bar"] , "k2": 5}`, false, true},
		{`{"k2" : 5, "k1":["bar",   "foo"]}`, `{"k1":["foo","bar"] , "k2": 5}`, false, false},
		{"not_json", `{"k1":["foo","bar"] , "k2": 5}`, true, false},
		{`{"k2" : 5, "k1":["foo",   "bar"]}`, "not_json", true, false},
	}

	for _, test := range tests {
		eq, err := JsonEquals(test.first, test.second)
		if err != nil && !test.errExpected {
			t.Errorf("Unexpected error: %s", err)
		} else if err == nil && test.errExpected {
			t.Errorf("Error was expected, but none returned")
		} else if eq != test.expectedEq {
			t.Errorf("Expected %+v, got %+v", test.expectedEq, eq)
		}
	}
}
