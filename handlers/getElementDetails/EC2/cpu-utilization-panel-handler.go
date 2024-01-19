package EC2

import (
	"awsx-api/log"
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/handler/EC2"
	"github.com/spf13/cobra"
	"net/http"
)

func ExecuteQuery(w http.ResponseWriter, r *http.Request) {
	log.Info("Starting /awsx-api/execute-query api")
	query := r.URL.Query().Get("query")
	elementType := r.URL.Query().Get("elementType")
	if query == "cpu_utilization_panel" {
		if elementType == "AWS/EC2" {
			GetCpuUtilizationPanel(w, r)
		}
	} else {
		http.Error(w, fmt.Sprintf("panel not found"), http.StatusBadRequest)
	}
	log.Info("Completed /awsx-api/execute-query api")
}

func GetCpuUtilizationPanel(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	region := r.URL.Query().Get("zone")
	cloudElementId := r.URL.Query().Get("cloudElementId")
	cloudElementApiUrl := r.URL.Query().Get("cloudElementApiUrl")
	crossAccountRoleArn := r.URL.Query().Get("crossAccountRoleArn")
	externalId := r.URL.Query().Get("externalId")
	responseType := r.URL.Query().Get("responseType")
	filter := r.URL.Query().Get("filter")
	instanceID := r.URL.Query().Get("instanceID")
	elementType := r.URL.Query().Get("elementType")
	startTime := r.URL.Query().Get("startTime")
	endTime := r.URL.Query().Get("endTime")

	//var err error

	commandParam := model.CommandParam{}

	if cloudElementId != "" {
		commandParam.CloudElementId = cloudElementId
		commandParam.CloudElementApiUrl = cloudElementApiUrl
		commandParam.Region = region
	} else {
		commandParam.CrossAccountRoleArn = crossAccountRoleArn
		commandParam.ExternalId = externalId
		commandParam.Region = region
	}
	authFlag, clientAuth, _ := authenticate.DoAuthenticate(commandParam)

	if authFlag {
		cmd := &cobra.Command{}
		cmd.PersistentFlags().StringVar(&instanceID, "instanceID", r.URL.Query().Get("instanceID"), "Description of the instanceID flag")
		cmd.PersistentFlags().StringVar(&elementType, "elementType", r.URL.Query().Get("elementType"), "Description of the elementType flag")
		cmd.PersistentFlags().StringVar(&startTime, "startTime", r.URL.Query().Get("startTime"), "Description of the startTime flag")
		cmd.PersistentFlags().StringVar(&endTime, "endTime", r.URL.Query().Get("endTime"), "Description of the endTime flag")
		cmd.PersistentFlags().StringVar(&responseType, "responseType", r.URL.Query().Get("responseType"), "responseType flag - json/frame")
		jsonString, cloudwatchMetricData, err := EC2.GetCpuUtilizationPanel(cmd, clientAuth)
		if err != nil {
			http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
			return
		}
		if responseType == "frame" {
			if filter == "SampleCount" {
				err = json.NewEncoder(w).Encode(cloudwatchMetricData["CurrentUsage"])
				if err != nil {
					http.Error(w, fmt.Sprintf("Exception: %s ", err), http.StatusInternalServerError)
					return
				}
			} else if filter == "Average" {
				err = json.NewEncoder(w).Encode(cloudwatchMetricData["AverageUsage"])
				if err != nil {
					http.Error(w, fmt.Sprintf(fmt.Sprintf("Exception: %s ", err)), http.StatusInternalServerError)
					return
				}
			} else if filter == "Maximum" {
				err = json.NewEncoder(w).Encode(cloudwatchMetricData["MaxUsage"])
				if err != nil {
					http.Error(w, fmt.Sprintf(fmt.Sprintf("Exception: %s ", err)), http.StatusInternalServerError)
					return
				}
			} else {
				err = json.NewEncoder(w).Encode(cloudwatchMetricData)
				if err != nil {
					http.Error(w, fmt.Sprintf(fmt.Sprintf("Exception: %s ", err)), http.StatusInternalServerError)
					return
				}
			}
		} else {
			type UsageData struct {
				AverageUsage float64 `json:"AverageUsage"`
				CurrentUsage float64 `json:"CurrentUsage"`
				MaxUsage     float64 `json:"MaxUsage"`
			}
			var data UsageData
			err := json.Unmarshal([]byte(jsonString), &data)
			if err != nil {
				http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
				return
			}

			// Marshal the struct back to JSON
			jsonBytes, err := json.Marshal(data)
			if err != nil {
				http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
				return
			}

			// Set Content-Type header and write the JSON response
			w.Header().Set("Content-Type", "application/json")
			_, err = w.Write(jsonBytes)
			if err != nil {
				http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
				return
			}
			//err = json.NewEncoder(w).Encode(jsonString)
			//if err != nil {
			//	http.Error(w, fmt.Sprintf(fmt.Sprintf("Exception: %s ", err)), http.StatusInternalServerError)
			//	return
			//}
		}
	}

}

//func GetEc2CommonHandler(w http.ResponseWriter, r *http.Request) {
//	command := r.URL.Query().Get("query")
//	if command == "cpu_utilization_panel" {
//		GetCpuUtilizationPanel(w, r)
//	} else {
//		http.Error(w, fmt.Sprintf("panel not found"), http.StatusBadRequest)
//	}
//}
//
//func GetCpuUtilizationPanel(w http.ResponseWriter, r *http.Request) {
//	log.Info("Starting /awsx/ec2/cpu-utilization-panel api")
//	w.Header().Set("Content-Type", "application/json")
//
//	region := r.URL.Query().Get("zone")
//	cloudElementId := r.URL.Query().Get("cloudElementId")
//	cloudElementApiUrl := r.URL.Query().Get("cloudElementApiUrl")
//	instanceID := r.URL.Query().Get("instanceID")
//	elementType := r.URL.Query().Get("elementType")
//
//	//query := r.URL.Query().Get("query")
//	startTimeStr := r.URL.Query().Get("startTime")
//	endTimeStr := r.URL.Query().Get("endTime")
//	statistic := r.URL.Query().Get("filter")
//	var startTime, endTime time.Time
//	var err error
//	if startTimeStr != "" {
//		startTime, err = time.Parse(time.RFC3339, startTimeStr)
//		if err != nil {
//			handleError(w, "Error parsing start time:", http.StatusBadRequest, err)
//			return
//		}
//	} else {
//		defaultStartTime := time.Now().Add(-5 * time.Minute)
//		startTime = defaultStartTime
//	}
//
//	if endTimeStr != "" {
//		endTime, err = time.Parse(time.RFC3339, endTimeStr)
//		if err != nil {
//			handleError(w, "Error parsing end time:", http.StatusBadRequest, err)
//			return
//		}
//	} else {
//		defaultEndTime := time.Now()
//		endTime = defaultEndTime
//	}
//
//	commandParam := model.CommandParam{}
//	if cloudElementId != "" {
//		commandParam.CloudElementId = cloudElementId
//		commandParam.CloudElementApiUrl = cloudElementApiUrl
//		commandParam.Region = region
//		authFlag, clientAuth, err := authenticate.DoAuthenticate(commandParam)
//		if err != nil || !authFlag {
//			log.Error(err.Error())
//			http.Error(w, fmt.Sprintf("Exception: "+err.Error()), http.StatusInternalServerError)
//			return
//		}
//		result, respErr := EC2.GetCpuUtilizationMetricData(clientAuth, elementType, instanceID, &startTime, &endTime, "")
//		if respErr != nil {
//			log.Error(respErr.Error())
//			http.Error(w, fmt.Sprintf("Exception: "+respErr.Error()), http.StatusInternalServerError)
//			return
//		}
//		json.NewEncoder(w).Encode(result)
//	} else {
//		crossAccountRoleArn := r.URL.Query().Get("crossAccountRoleArn")
//		externalId := r.URL.Query().Get("externalId")
//		returnType := r.URL.Query().Get("returnType")
//		commandParam.CrossAccountRoleArn = crossAccountRoleArn
//		commandParam.ExternalId = externalId
//		commandParam.Region = region
//		authFlag, clientAuth, err := authenticate.DoAuthenticate(commandParam)
//		if err != nil || !authFlag {
//			log.Error(err.Error())
//			http.Error(w, fmt.Sprintf("Exception: "+err.Error()), http.StatusInternalServerError)
//			return
//		}
//		if returnType == "json" {
//			jsonOutput := map[string]float64{}
//			statistics := strings.Split(statistic, ",")
//			for _, stat := range statistics {
//				result, respErr := EC2.GetCpuUtilizationMetricData(clientAuth, instanceID, elementType, &startTime, &endTime, stat)
//				if respErr != nil {
//					///log.Error(respErr.Error())
//					http.Error(w, fmt.Sprintf("Exception: "+respErr.Error()), http.StatusInternalServerError)
//					return
//				}
//				if stat == "SampleCount" {
//					jsonOutput["CurrentUsage"] = *result.MetricDataResults[0].Values[0]
//				} else if stat == "Average" {
//					jsonOutput["AverageUsage"] = *result.MetricDataResults[0].Values[0]
//				} else if stat == "Maximum" {
//					jsonOutput["MaxUsage"] = *result.MetricDataResults[0].Values[0]
//				} else {
//					http.Error(w, fmt.Sprintf("statistics not supported"), http.StatusBadRequest)
//				}
//			}
//			jsonString, err := json.Marshal(jsonOutput)
//			if err != nil {
//				log.Fatal(err)
//			}
//			fmt.Println("Output:", string(jsonString))
//			json.NewEncoder(w).Encode(string(jsonString))
//		} else {
//			statistics := strings.Split(statistic, ",")
//			for _, stat := range statistics {
//				result, respErr := EC2.GetCpuUtilizationMetricData(clientAuth, instanceID, elementType, &startTime, &endTime, stat)
//				if respErr != nil {
//					///log.Error(respErr.Error())
//					http.Error(w, fmt.Sprintf("Exception: "+respErr.Error()), http.StatusInternalServerError)
//					return
//				}
//
//				json.NewEncoder(w).Encode(result)
//			}
//		}
//	}
//	log.Info("/awsx-metric/metric completed")
//}
//
//func handleError(w http.ResponseWriter, logMsg string, statusCode int, err error) {
//	log.Error(logMsg, err)
//	http.Error(w, fmt.Sprintf("Exception: %s", logMsg), statusCode)
//}
