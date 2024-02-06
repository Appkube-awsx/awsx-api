package ECS

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Appkube-awsx/awsx-getelementdetails/handler/ECS"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"

	"github.com/spf13/cobra"
)

func GetContainerPanel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract query parameters
	region := r.URL.Query().Get("zone")
	cloudElementId := r.URL.Query().Get("cloudElementId")
	cloudElementApiUrl := r.URL.Query().Get("cloudElementApiUrl")
	crossAccountRoleArn := r.URL.Query().Get("crossAccountRoleArn")
	externalId := r.URL.Query().Get("externalId")
	responseType := r.URL.Query().Get("responseType")
	filter := r.URL.Query().Get("filter")
	clusterName := r.URL.Query().Get("clusterName")
	elementType := r.URL.Query().Get("elementType")
	startTime := r.URL.Query().Get("startTime")
	endTime := r.URL.Query().Get("endTime")
	query := r.URL.Query().Get("query")

	// Initialize CommandParam
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

	// Authenticate
	authFlag, clientAuth, _ := authenticate.DoAuthenticate(commandParam)

	if authFlag {
		cmd := &cobra.Command{}
		cmd.PersistentFlags().StringVar(&clusterName, "clusterName", r.URL.Query().Get("clusterName"), "Description of the clusterName flag")
		cmd.PersistentFlags().StringVar(&elementType, "elementType", r.URL.Query().Get("elementType"), "Description of the elementType flag")
		cmd.PersistentFlags().StringVar(&startTime, "startTime", r.URL.Query().Get("startTime"), "Description of the startTime flag")
		cmd.PersistentFlags().StringVar(&endTime, "endTime", r.URL.Query().Get("endTime"), "Description of the endTime flag")
		cmd.PersistentFlags().StringVar(&responseType, "responseType", r.URL.Query().Get("responseType"), "responseType flag - json/frame")

		// Call ECS package function to get metrics
		jsonString, ecsMetricData, err := ECS.GetContainerPanel(cmd, clientAuth)
		if err != nil {
			http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
			return
		}
		if elementType == "ContainerInsights" && query == "Cpu_utilization_panel" {
			if responseType == "frame" {
				switch filter {
				case "Current":
					err = json.NewEncoder(w).Encode(ecsMetricData["CurrentUsage"])
				case "Average":
					err = json.NewEncoder(w).Encode(ecsMetricData["AverageUsage"])
				case "Maximum":
					err = json.NewEncoder(w).Encode(ecsMetricData["MaxUsage"])
				default:
					fmt.Println("this is else json", ecsMetricData)
					err = json.NewEncoder(w).Encode(ecsMetricData)
				}

				if err != nil {
					http.Error(w, fmt.Sprintf("Exception: %s ", err), http.StatusInternalServerError)
					return
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
			}
		}
	}
}
