package EKS

import (
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

type NodeDownData struct {
	Timestamp    string  `json:"Timestamp"`
	NodeDowntime float64 `json:"NodeDowntime"`
}

var (
	downtimeAuthCache       sync.Map
	downtimeClientCache     sync.Map
	downtimeAuthCacheLock   sync.RWMutex
	downtimeClientCacheLock sync.RWMutex
)

// GetEKSDowntimePanel handles the request for the node downtime panel data
func GetEKSDowntimePanel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse query parameters
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

	// Authenticate and get client authentication details
	clientAuth, err := downtimeAuthenticateAndCache(commandParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication failed: %s", err), http.StatusInternalServerError)
		return
	}

	// Get CloudWatch client
	cloudWatchClient, err := downtimeCloudwatchClientCache(*clientAuth)
	if err != nil {
		http.Error(w, fmt.Sprintf("CloudWatch client creation/store in cache failed: %s", err), http.StatusInternalServerError)
		return
	}

	if clientAuth != nil {
		// Prepare cobra command
		cmd := &cobra.Command{}
		cmd.PersistentFlags().StringVar(&elementId, "elementId", r.URL.Query().Get("elementId"), "Description of the cloudElementID flag")
		cmd.PersistentFlags().StringVar(&instanceId, "instanceId", r.URL.Query().Get("instanceId"), "Description of the instanceID flag")
		cmd.PersistentFlags().StringVar(&elementType, "elementType", r.URL.Query().Get("elementType"), "Description of the elementType flag")
		cmd.PersistentFlags().StringVar(&startTime, "startTime", r.URL.Query().Get("startTime"), "Description of the startTime flag")
		cmd.PersistentFlags().StringVar(&endTime, "endTime", r.URL.Query().Get("endTime"), "Description of the endTime flag")
		cmd.PersistentFlags().StringVar(&responseType, "responseType", r.URL.Query().Get("responseType"), "responseType flag - json/frame")

		// Get node downtime panel data
		jsonData, _, err := EKS.GetNodeDowntimePanel(cmd, clientAuth, cloudWatchClient)

		var data []NodeDownData
		if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
			fmt.Println("Error:", err)
			return
		}
		formattedData := map[string]interface{}{
			"Node Downtime": map[string]interface{}{
				"Messages": nil,
				"MetricDataResults": []map[string]interface{}{
					{
						"Id":         "",
						"Label":      "",
						"Messages":   nil,
						"StatusCode": "Complete",
						"Timestamps": getTimestampNodeDowntimePanel(data),
						"Values":     getValuesNodeDowntimePanel(data),
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

// downtimeAuthenticateAndCache authenticates and caches client details
func downtimeAuthenticateAndCache(commandParam model.CommandParam) (*model.Auth, error) {
	cacheKey := commandParam.CloudElementId

	downtimeAuthCacheLock.Lock()
	defer downtimeAuthCacheLock.Unlock()

	if auth, ok := downtimeAuthCache.Load(cacheKey); ok {
		return auth.(*model.Auth), nil
	}

	_, clientAuth, err := authenticate.DoAuthenticate(commandParam)
	if err != nil {
		return nil, err
	}

	downtimeAuthCache.Store(cacheKey, clientAuth)
	return clientAuth, nil
}

// downtimeCloudwatchClientCache caches CloudWatch client
func downtimeCloudwatchClientCache(clientAuth model.Auth) (*cloudwatch.CloudWatch, error) {
	cacheKey := clientAuth.CrossAccountRoleArn

	downtimeClientCacheLock.Lock()
	defer downtimeClientCacheLock.Unlock()

	if client, ok := downtimeClientCache.Load(cacheKey); ok {
		return client.(*cloudwatch.CloudWatch), nil
	}

	cloudWatchClient := awsclient.GetClient(clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	downtimeClientCache.Store(cacheKey, cloudWatchClient)
	return cloudWatchClient, nil
}

func getTimestampNodeDowntimePanel(data []NodeDownData) []string {
	timestamps := make([]string, len(data))
	for i, d := range data {
		timestamps[i] = d.Timestamp
	}
	return timestamps
}
func getValuesNodeDowntimePanel(data []NodeDownData) []float64 {
	values := make([]float64, len(data))
	for i, d := range data {
		values[i] = d.NodeDowntime
	}
	return values
}
