package helpers

import (
	"errors"
	"strconv"
	"strings"
)

func GetMetricDataFromUri[V int64 | float64](url string) (key string, val V, err error) {
	splittedUrl := strings.Split(url, "/")
	if len(splittedUrl) < 5 {
		return "", 0, errors.New("incorrect data has been provided")
	}
	metricType := splittedUrl[2]
	metricKey := splittedUrl[3]
	metricVal := splittedUrl[4]

	if metricKey == "" {
		return "", 0, errors.New("provide metric's key")
	}

	if metricType == "counter" {
		metricValInt, err := strconv.Atoi(metricVal)
		if err != nil {
			return "", 0, err
		}
		return metricKey, V(metricValInt), nil
	} else if metricType == "gauge" {
		metricValFloat, err := strconv.ParseFloat(metricVal, 64)
		if err != nil {
			return "", 0, err
		}
		return metricKey, V(metricValFloat), nil
	} else {
		return "", 0, errors.New("incorrect metric type was providen")
	}
}

func GetMetricKeyFromUrl(url string) (string, error) {
	splittedUrl := strings.Split(url, "/")
	if len(splittedUrl) != 4 {
		return "", errors.New("incorrect data has been provided")
	}
	return splittedUrl[3], nil

}
