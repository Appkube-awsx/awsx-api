package EC2

import (
	// "awsx-api/handlers/getElementDetails/EKS"
	"awsx-api/log"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/handler/EC2"
	"github.com/spf13/cobra"
)

func ExecuteNetworkQuery(w http.ResponseWriter, r *http.Request) {
	log.Info("Starting /awsx-api/execute-query api")
	query := r.URL.Query().Get("query")
	elementType := r.URL.Query().Get("elementType")
	if elementType == "AWS/EC2" && query == "network_utilization_panel" {

		GetNetworkUtilizationPanel(w, r)

	} else {
		http.Error(w, fmt.Sprintf("panel not found"), http.StatusBadRequest)
	}
	log.Info("Completed /awsx-api/execute-query api")
}

func GetNetworkUtilizationPanel(w http.ResponseWriter, r *http.Request) {

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
		jsonString, cloudwatchMetricData, err := EC2.GetNetworkUtilizationPanel(cmd, clientAuth)
		if err != nil {
			http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
			return
		}
		if responseType == "frame" {
			if filter == "SampleCount" {
				err = json.NewEncoder(w).Encode(cloudwatchMetricData["inboundTraffic"])
				if err != nil {
					http.Error(w, fmt.Sprintf("Exception: %s ", err), http.StatusInternalServerError)
					return
				}
			} else if filter == "Average" {
				err = json.NewEncoder(w).Encode(cloudwatchMetricData["OutboundTraffic"])
				if err != nil {
					http.Error(w, fmt.Sprintf(fmt.Sprintf("Exception: %s ", err)), http.StatusInternalServerError)
					return
				}
			} else if filter == "dataTransferred" {
				err = json.NewEncoder(w).Encode(cloudwatchMetricData["DataTransferred"])
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
				InboundTraffic  float64 `json:"inboundTraffic"`
				OutboundTraffic float64 `json:"outboundTraffic"`
				DataTransferred float64 `json:"dataTransferred"`
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
		}
	}

}
