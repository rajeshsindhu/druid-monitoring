package util

import (
	"encoding/json"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

const timeFormat = "2006-01-02T15:04:05.999Z"

func GetTimeRange(latestTime int64, granularity string) (string, string) {

	startTime := int64(0)
	switch gran := granularity; gran {
	case "day":
		startTime = latestTime - 86400
	case "hour":
		startTime = latestTime - 3600
	case "thirty_minute":
		startTime = latestTime - 1800
	case "fifteen_minute":
		startTime = latestTime - 900
	default:
		startTime = latestTime - 86000
	}

	sTime := time.Unix(startTime, 0).Format(timeFormat)
	eTime := time.Unix(latestTime, 0).Format(timeFormat)
	return sTime, eTime
}

func ExecutePostQuery(druidEndPoint string, dqRetries int, query string) ([]byte, error) {

	url := druidEndPoint + "/druid/v2"
	contentType := "application/json"
	var queryErr error
	for retries := 0; retries <= dqRetries; retries++ {
		rsp, err := http.Post(url, contentType, strings.NewReader(query))

		if err != nil || rsp.StatusCode != 200 {
			glog.Error("Received Error: " + err.Error())
			queryErr = err

			glog.Info("Retrying...")
			continue
		}

		defer rsp.Body.Close()
		body, err := ioutil.ReadAll(rsp.Body)

		if err != nil {
			glog.Error("Received Error: " + err.Error())
			queryErr = err

			glog.Info("Retrying...")
			continue
		}

		return body, nil

	}

	glog.Info("Giving up...")

	return nil, queryErr
}

func ExecuteGetQuery(druidEndPoint string, dqRetries int, query string) ([]byte, error) {

	url := druidEndPoint + "/druid/v2"
	url = url + query
	var queryErr error
	for retries := 0; retries <= dqRetries; retries++ {
		rsp, err := http.Get(url)

		if err != nil || rsp.StatusCode != 200 {
			glog.Error("Received Error: " + err.Error())
			queryErr = err

			glog.Info("Retrying...")
			continue
		}

		defer rsp.Body.Close()
		body, err := ioutil.ReadAll(rsp.Body)

		if err != nil {
			glog.Error("Received Error: " + err.Error())
			queryErr = err

			glog.Info("Retrying...")
			continue
		}

		return body, nil

	}

	glog.Info("Giving up...")

	return nil, queryErr
}

// Check whether two json blobs are semantically equal
func JsonEquals(j1, j2 string) (bool, error) {
	var i1, i2 interface{}
	if err := json.Unmarshal([]byte(j1), &i1); err != nil {
		glog.Error("Json unmarshalling  error : " + err.Error())
		return false, err
	}
	if err := json.Unmarshal([]byte(j2), &i2); err != nil {
		glog.Error("Json unmarshalling  error : " + err.Error())
		return false, err
	}
	return reflect.DeepEqual(i1, i2), nil
}
