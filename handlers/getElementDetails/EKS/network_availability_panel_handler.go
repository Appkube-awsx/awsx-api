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

var (
	networkAuthCache       sync.Map
	networkClientCache     sync.Map
	networkAuthCacheLock   sync.RWMutex
	networkClientCacheLock sync.RWMutex
)

type NodeDataNetwork struct {
	Timestamp       string  `json:"Timestamp"`
	NodeNetworktime float64 `json:"NodeDataNetwork"`
}

// GetEKSNetworkAvailabilityPanel handles the request for the network availability panel data
func GetEKSNetworkAvailabilityPanel(w http.ResponseWriter, r *http.Request) {
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
	clientAuth, err := networkAuthenticateAndCache(commandParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication failed: %s", err), http.StatusInternalServerError)
		return
	}

	// Get CloudWatch client
	cloudwatchClient, err := networkCloudwatchClientCache(*clientAuth)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cloudwatch client creation/store in cache failed: %s", err), http.StatusInternalServerError)
		return
	}

	if clientAuth != nil {
		// Prepare cobra command
		cmd := &cobra.Command{}
		cmd.PersistentFlags().StringVar(&elementId, "elementId", r.URL.Query().Get("elementId"), "Description of the elementId flag")
		cmd.PersistentFlags().StringVar(&instanceId, "instanceId", r.URL.Query().Get("instanceId"), "Description of the instanceId flag")
		cmd.PersistentFlags().StringVar(&elementType, "elementType", r.URL.Query().Get("elementType"), "Description of the elementType flag")
		cmd.PersistentFlags().StringVar(&startTime, "startTime", r.URL.Query().Get("startTime"), "Description of the startTime flag")
		cmd.PersistentFlags().StringVar(&endTime, "endTime", r.URL.Query().Get("endTime"), "Description of the endTime flag")
		cmd.PersistentFlags().StringVar(&responseType, "responseType", r.URL.Query().Get("responseType"), "responseType flag - json/frame")

		// Get network availability panel data
		jsonString, _, err := EKS.GetNetworkAvailabilityData(cmd, clientAuth, cloudwatchClient)
		var data []NodeDataNetwork
		if err := json.Unmarshal([]byte(jsonString), &data); err != nil {
			fmt.Println("Error:", err)
			return
		}
		formattedData := map[string]interface{}{
			"Network Availability": map[string]interface{}{
				"Messages": nil,
				"MetricDataResults": []map[string]interface{}{
					{
						"Id":         "",
						"Label":      "",
						"Messages":   nil,
						"StatusCode": "Complete",
						"Timestamps": getTimestampsNodeDataNetwork(data),
						"Values":     getValuesNodeDataNetwork(data),
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

// networkAuthenticateAndCache authenticates and caches client details for network availability
func networkAuthenticateAndCache(commandParam model.CommandParam) (*model.Auth, error) {
	cacheKey := commandParam.CloudElementId

	networkAuthCacheLock.Lock()
	defer networkAuthCacheLock.Unlock()

	if auth, ok := networkAuthCache.Load(cacheKey); ok {
		return auth.(*model.Auth), nil
	}

	_, clientAuth, err := authenticate.DoAuthenticate(commandParam)
	if err != nil {
		return nil, err
	}

	networkAuthCache.Store(cacheKey, clientAuth)
	return clientAuth, nil
}

// networkCloudwatchClientCache caches cloudwatch client for network availability
func networkCloudwatchClientCache(clientAuth model.Auth) (*cloudwatch.CloudWatch, error) {
	cacheKey := clientAuth.CrossAccountRoleArn

	networkClientCacheLock.Lock()
	defer networkClientCacheLock.Unlock()

	if client, ok := networkClientCache.Load(cacheKey); ok {
		return client.(*cloudwatch.CloudWatch), nil
	}

	cloudWatchClient := awsclient.GetClient(clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	networkClientCache.Store(cacheKey, cloudWatchClient)
	return cloudWatchClient, nil
}
func getTimestampsNodeDataNetwork(data []NodeDataNetwork) []string {
	timestamps := make([]string, len(data))
	for i, d := range data {
		timestamps[i] = d.Timestamp
	}
	return timestamps
}
func getValuesNodeDataNetwork(data []NodeDataNetwork) []float64 {
	values := make([]float64, len(data))
	for i, d := range data {
		values[i] = d.NodeNetworktime
	}
	return values
}
