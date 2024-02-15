package ECS

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/handler/ECS"
	"github.com/spf13/cobra"
)

func GetStorageUtilizationPanel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
		cmd.PersistentFlags().StringVar(&clusterName, "clusterName", r.URL.Query().Get("clusterName"), "Description of the clusterName flag")
		cmd.PersistentFlags().StringVar(&elementType, "elementType", r.URL.Query().Get("elementType"), "Description of the elementType flag")
		cmd.PersistentFlags().StringVar(&startTime, "startTime", r.URL.Query().Get("startTime"), "Description of the startTime flag")
		cmd.PersistentFlags().StringVar(&endTime, "endTime", r.URL.Query().Get("endTime"), "Description of the endTime flag")
		cmd.PersistentFlags().StringVar(&responseType, "responseType", r.URL.Query().Get("responseType"), "responseType flag - json/frame")
		jsonString, cloudwatchMetricData, err := ECS.GetStorageUtilizationPanel(cmd, clientAuth)
		if err != nil {
			http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
			return
		}
		if elementType == "ContainerInsights" && query == "Storage_utilization_panel" {
			if responseType == "frame" {
				if filter == "RootVolume" {
					err = json.NewEncoder(w).Encode(cloudwatchMetricData["RootVolumeUsage"])
					if err != nil {
						http.Error(w, fmt.Sprintf("Exception: %s ", err), http.StatusInternalServerError)
						return
					}
				} else if filter == "EBSVolume1" {
					err = json.NewEncoder(w).Encode(cloudwatchMetricData["EBSVolume1Usage"])
					if err != nil {
						http.Error(w, fmt.Sprintf(fmt.Sprintf("Exception: %s ", err)), http.StatusInternalServerError)
						return
					}
				} else if filter == "EBSVolume2" {
					err = json.NewEncoder(w).Encode(cloudwatchMetricData["EBSVolume2Usage"])
					if err != nil {
						http.Error(w, fmt.Sprintf(fmt.Sprintf("Exception: %s ", err)), http.StatusInternalServerError)
						return
					}
				} else {
					fmt.Println("this is else json", cloudwatchMetricData)
					err = json.NewEncoder(w).Encode(cloudwatchMetricData)
					if err != nil {
						http.Error(w, fmt.Sprintf(fmt.Sprintf("Exception: %s ", err)), http.StatusInternalServerError)
						return
					}
				}
			} else {
				type StorageUtilizationResult struct {
					RootVolumeUsage float64 `json:"rootVolumeUsage"`
					EBSVolume1Usage float64 `json:"ebsVolume1Usage"`
					EBSVolume2Usage float64 `json:"ebsVolume2Usage"`
				}
				var data StorageUtilizationResult
				err := json.Unmarshal([]byte(jsonString), &data)
				if err != nil {
					http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
					return
				}

				jsonBytes, err := json.Marshal(data)
				if err != nil {
					http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
					return
				}
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
