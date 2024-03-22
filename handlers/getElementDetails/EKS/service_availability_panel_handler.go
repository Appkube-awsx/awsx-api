package EKS

import (
	"awsx-api/log"
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/handler/EKS"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
	"net/http"
	"sync"
)

var (
	serviceAuthCache       sync.Map
	serviceClientCache     sync.Map
	serviceAuthCacheLock   sync.RWMutex
	serviceClientCacheLock sync.RWMutex
)

type ServiceTimeSeriesDataPoint struct {
	Timestamp    string  `json:"Timestamp"`
	Availability float64 `json:"Availability"`
}

func GetEKSServiceAvailabilityPanel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	region := r.URL.Query().Get("zone")
	elementId := r.URL.Query().Get("elementId")
	elementApiUrl := r.URL.Query().Get("elementApiUrl")
	crossAccountRoleArn := r.URL.Query().Get("crossAccountRoleArn")
	externalId := r.URL.Query().Get("externalId")
	responseType := r.URL.Query().Get("responseType")
	instanceId := r.URL.Query().Get("instanceId")
	elementType := r.URL.Query().Get("elementType")
	startTime := r.URL.Query().Get("startTime")
	endTime := r.URL.Query().Get("endTime")

	commandParam := model.CommandParam{}

	if elementId != "" {
		commandParam.CloudElementId = elementId
		commandParam.CloudElementApiUrl = elementApiUrl
		commandParam.Region = region
	} else {
		commandParam.CrossAccountRoleArn = crossAccountRoleArn
		commandParam.ExternalId = externalId
		commandParam.Region = region
	}

	clientAuth, err := serviceAuthenticateAndCache(commandParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication failed: %s", err), http.StatusInternalServerError)
		return
	}

	cloudwatchClient, err := serviceCloudwatchClientCache(*clientAuth)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cloudwatch client creation/store in cache failed: %s", err), http.StatusInternalServerError)
		return
	}

	if clientAuth != nil {
		cmd := &cobra.Command{}
		cmd.PersistentFlags().StringVar(&elementId, "elementId", r.URL.Query().Get("elementId"), "Description of the elementId flag")
		cmd.PersistentFlags().StringVar(&instanceId, "instanceId", r.URL.Query().Get("instanceId"), "Description of the instanceId flag")
		cmd.PersistentFlags().StringVar(&elementType, "elementType", r.URL.Query().Get("elementType"), "Description of the elementType flag")
		cmd.PersistentFlags().StringVar(&startTime, "startTime", r.URL.Query().Get("startTime"), "Description of the startTime flag")
		cmd.PersistentFlags().StringVar(&endTime, "endTime", r.URL.Query().Get("endTime"), "Description of the endTime flag")
		cmd.PersistentFlags().StringVar(&responseType, "responseType", r.URL.Query().Get("responseType"), "responseType flag - json/frame")

		jsonString, _, err := EKS.GetServiceAvailabilityData(cmd, clientAuth, cloudwatchClient)
		if err != nil {
			http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
			return
		}
		log.Infof("response type: %s", responseType)

		var data []ServiceTimeSeriesDataPoint
		if err := json.Unmarshal([]byte(jsonString), &data); err != nil {
			fmt.Println("Error:", err)
			return
		}
		formattedData := map[string]interface{}{
			"Service Availability": map[string]interface{}{
				"Messages": nil,
				"MetricDataResults": []map[string]interface{}{
					{
						"Id":         "",
						"Label":      "",
						"Messages":   nil,
						"StatusCode": "Complete",
						"Timestamps": getTimestampsService(data),
						"Values":     getValueService(data),
					},
				},
			},
		}

		// Convert the formatted data to JSON
		formattedJson, err := json.MarshalIndent(formattedData, "", "    ")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println(string(formattedJson))

		// Handle JSON response
		w.Header().Set("Content-Type", "application/json")
		write, err := w.Write([]byte(formattedJson))
		if err != nil {
			http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
			return
		}
		fmt.Println(write)
	}
}

func serviceAuthenticateAndCache(commandParam model.CommandParam) (*model.Auth, error) {
	cacheKey := commandParam.CloudElementId

	serviceAuthCacheLock.Lock()
	defer serviceAuthCacheLock.Unlock()

	if auth, ok := serviceAuthCache.Load(cacheKey); ok {
		return auth.(*model.Auth), nil
	}

	_, clientAuth, err := authenticate.DoAuthenticate(commandParam)
	if err != nil {
		return nil, err
	}

	serviceAuthCache.Store(cacheKey, clientAuth)
	return clientAuth, nil
}

func serviceCloudwatchClientCache(clientAuth model.Auth) (*cloudwatch.CloudWatch, error) {
	cacheKey := clientAuth.CrossAccountRoleArn

	serviceClientCacheLock.Lock()
	defer serviceClientCacheLock.Unlock()

	if client, ok := serviceClientCache.Load(cacheKey); ok {
		return client.(*cloudwatch.CloudWatch), nil
	}

	cloudWatchClient := awsclient.GetClient(clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	serviceClientCache.Store(cacheKey, cloudWatchClient)
	return cloudWatchClient, nil
}
func getTimestampsService(data []ServiceTimeSeriesDataPoint) []string {
	timestamps := make([]string, len(data))
	for i, d := range data {
		timestamps[i] = d.Timestamp
	}
	return timestamps
}
func getValueService(data []ServiceTimeSeriesDataPoint) []float64 {
	values := make([]float64, len(data))
	for i, d := range data {
		values[i] = d.Availability
	}
	return values
}
