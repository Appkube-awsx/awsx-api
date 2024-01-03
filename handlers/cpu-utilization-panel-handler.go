package handlers

import (
	"awsx-api/log"
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/handler/EC2"
	"strings"
	"time"

	"net/http"
)

func GetCpuUtilizationPanel(w http.ResponseWriter, r *http.Request) {
	log.Info("Starting /awsx/ec2/cpu-utilization-panel api")
	w.Header().Set("Content-Type", "application/json")
	region := r.URL.Query().Get("zone")
	cloudElementId := r.URL.Query().Get("cloudElementId")
	cloudElementApiUrl := r.URL.Query().Get("cloudElementApiUrl")
	instanceID := r.URL.Query().Get("instanceID")
	elementType := r.URL.Query().Get("elementType")

	query := r.URL.Query().Get("query")
	startTimeStr := r.URL.Query().Get("startTime")
	endTimeStr := r.URL.Query().Get("endTime")
	statistic := r.URL.Query().Get("statistic")
	var startTime, endTime time.Time
	var err error
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			handleError(w, "Error parsing start time:", http.StatusBadRequest, err)
			return
		}
	} else {
		defaultStartTime := time.Now().Add(-5 * time.Minute)
		startTime = defaultStartTime
	}

	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			handleError(w, "Error parsing end time:", http.StatusBadRequest, err)
			return
		}
	} else {
		defaultEndTime := time.Now()
		endTime = defaultEndTime
	}

	commandParam := model.CommandParam{}
	if cloudElementId != "" {
		commandParam.CloudElementId = cloudElementId
		commandParam.CloudElementApiUrl = cloudElementApiUrl
		commandParam.Region = region
		authFlag, clientAuth, err := authenticate.DoAuthenticate(commandParam)
		if err != nil || !authFlag {
			log.Error(err.Error())
			http.Error(w, fmt.Sprintf("Exception: "+err.Error()), http.StatusInternalServerError)
			return
		}
		result, respErr := EC2.GetCpuUtilizationMetricData(clientAuth, elementType, instanceID, query, &startTime, &endTime, "")
		if respErr != nil {
			log.Error(respErr.Error())
			http.Error(w, fmt.Sprintf("Exception: "+respErr.Error()), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(result)
	} else {
		crossAccountRoleArn := r.URL.Query().Get("crossAccountRoleArn")
		externalId := r.URL.Query().Get("externalId")

		commandParam.CrossAccountRoleArn = crossAccountRoleArn
		commandParam.ExternalId = externalId
		commandParam.Region = region
		authFlag, clientAuth, err := authenticate.DoAuthenticate(commandParam)
		if err != nil || !authFlag {
			log.Error(err.Error())
			http.Error(w, fmt.Sprintf("Exception: "+err.Error()), http.StatusInternalServerError)
			return
		}
		jsonOutput := map[string]float64{}
		statistics := strings.Split(statistic, ",")
		for _, stat := range statistics {
			result, respErr := EC2.GetCpuUtilizationMetricData(clientAuth, instanceID, query, elementType, &startTime, &endTime, stat)
			if respErr != nil {
				///log.Error(respErr.Error())
				http.Error(w, fmt.Sprintf("Exception: "+respErr.Error()), http.StatusInternalServerError)
				return
			}
			if stat == "SampleCount" {
				jsonOutput["CurrentUsage"] = *result.MetricDataResults[0].Values[0]
			} else if stat == "Average" {
				jsonOutput["AverageUsage"] = *result.MetricDataResults[0].Values[0]
			} else if stat == "Maximum" {
				jsonOutput["MaxUsage"] = *result.MetricDataResults[0].Values[0]
			} else {
				http.Error(w, fmt.Sprintf("statistics not supported"), http.StatusBadRequest)
			}
		}
		jsonString, err := json.Marshal(jsonOutput)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Output:", string(jsonString))

		json.NewEncoder(w).Encode(string(jsonString))
	}
	log.Info("/awsx-metric/metric completed")
}

func handleError(w http.ResponseWriter, logMsg string, statusCode int, err error) {
	log.Error(logMsg, err)
	http.Error(w, fmt.Sprintf("Exception: %s", logMsg), statusCode)
}
