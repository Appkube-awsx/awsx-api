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
	networkThroughputSingleAuthCache       sync.Map
	networkThroughputSingleClientCache     sync.Map
	networkThroughputSingleAuthCacheLock   sync.RWMutex
	networkThroughputSingleClientCacheLock sync.RWMutex
)

type OriginalResponse struct {
	Messages          interface{} `json:"Messages"`
	MetricDataResults []struct {
		ID         string      `json:"Id"`
		Label      string      `json:"Label"`
		Messages   interface{} `json:"Messages"`
		StatusCode string      `json:"StatusCode"`
		Timestamps interface{} `json:"Timestamps"`
		Values     interface{} `json:"Values"`
	} `json:"MetricDataResults"`
	NextToken interface{} `json:"NextToken"`
}

type TransformedResponse struct {
	NetworkThroughput struct {
		Messages          interface{} `json:"Messages"`
		MetricDataResults []struct {
			ID         string      `json:"Id"`
			Label      string      `json:"Label"`
			Messages   interface{} `json:"Messages"`
			StatusCode string      `json:"StatusCode"`
			Timestamps interface{} `json:"Timestamps"`
			Values     interface{} `json:"Values"`
		} `json:"MetricDataResults"`
		NextToken interface{} `json:"NextToken"`
	} `json:"Network Throughput"`
}

func GetNetworkThroughputSinglePanel(w http.ResponseWriter, r *http.Request) {
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
	clientAuth, err := networkThroughputSingleAuthenticateAndCache(commandParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication failed: %s", err), http.StatusInternalServerError)
		return
	}

	// Get CloudWatch client
	cloudwatchClient, err := networkThroughputSingleCloudwatchClientCache(*clientAuth)
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

		// Get network throughput single panel data
		jsonString, _, err := EKS.GetNetworkThroughputSinglePanel(cmd, clientAuth, cloudwatchClient)
		if err != nil {
			http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
			return
		}
		transformedData, err := transformResponse(jsonString)
		if err != nil {
			http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
			return
		}

		// Write response
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(transformedData); err != nil {
			http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
			return
		}

	}
}

// networkThroughputSingleAuthenticateAndCache authenticates and caches client details for network throughput single panel
func networkThroughputSingleAuthenticateAndCache(commandParam model.CommandParam) (*model.Auth, error) {
	cacheKey := commandParam.CloudElementId

	networkThroughputSingleAuthCacheLock.Lock()
	defer networkThroughputSingleAuthCacheLock.Unlock()

	if auth, ok := networkThroughputSingleAuthCache.Load(cacheKey); ok {
		return auth.(*model.Auth), nil
	}

	_, clientAuth, err := authenticate.DoAuthenticate(commandParam)
	if err != nil {
		return nil, err
	}

	networkThroughputSingleAuthCache.Store(cacheKey, clientAuth)
	return clientAuth, nil
}

// networkThroughputSingleCloudwatchClientCache caches cloudwatch client for network throughput single panel
func networkThroughputSingleCloudwatchClientCache(clientAuth model.Auth) (*cloudwatch.CloudWatch, error) {
	cacheKey := clientAuth.CrossAccountRoleArn

	networkThroughputSingleClientCacheLock.Lock()
	defer networkThroughputSingleClientCacheLock.Unlock()

	if client, ok := networkThroughputSingleClientCache.Load(cacheKey); ok {
		return client.(*cloudwatch.CloudWatch), nil
	}

	cloudWatchClient := awsclient.GetClient(clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	networkThroughputSingleClientCache.Store(cacheKey, cloudWatchClient)
	return cloudWatchClient, nil
}
func transformResponse(jsonString *cloudwatch.GetMetricDataOutput) ([]byte, error) {

	newJson, err := json.Marshal(jsonString)
	if err != nil {
		return nil, err
	}

	var originalResp OriginalResponse
	if err := json.Unmarshal(newJson, &originalResp); err != nil {
		return nil, err
	}

	transformedResp := TransformedResponse{
		NetworkThroughput: originalResp,
	}
	transformedData, err := json.MarshalIndent(transformedResp, "", "    ")
	if err != nil {
		return nil, err
	}

	return transformedData, nil
}
