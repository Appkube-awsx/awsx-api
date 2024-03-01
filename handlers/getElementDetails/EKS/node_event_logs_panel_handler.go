package EKS

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/handler/EKS"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

// NodeEventLog represents a node event log entry.
type NodeEventLog struct {
	Timestamp       int64  `json:"Timestamp"`
	EventType       string `json:"EventType"`
	SourceComponent string `json:"SourceComponent"`
	EventMessage    string `json:"EventMessage"`
}

var (
	eventLogsAuthCache       sync.Map
	eventLogsClientCache     sync.Map
	eventLogsAuthCacheLock   sync.RWMutex
	eventLogsClientCacheLock sync.RWMutex
)

// GetEKSEventLogsPanel handles the request for the node event logs panel data
func GetEKSEventLogsPanel(w http.ResponseWriter, r *http.Request) {
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
	clientAuth, err := authenticateAndCache(commandParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication failed: %s", err), http.StatusInternalServerError)
		return
	}

	// Get CloudWatch client
	cloudWatchClient, err := eventLogsCloudwatchClientCache(*clientAuth)
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

		// Get node event logs panel data
		jsonString, rawLogs, err := EKS.GetNodeEventLogsSinglePanel(cmd, clientAuth, cloudWatchClient)
		if err != nil {
			http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
			return
		}

		if responseType == "frame" {
			err = json.NewEncoder(w).Encode(rawLogs)
			if err != nil {
				http.Error(w, fmt.Sprintf("Exception: %s ", err), http.StatusInternalServerError)
				return
			}
		} else {
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(jsonString))
			if err != nil {
				http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
				return
			}
		}
	}
}

// authenticateAndCache authenticates and caches client details
func eventLogsauthenticateAndCache(commandParam model.CommandParam) (*model.Auth, error) {
	cacheKey := commandParam.CloudElementId

	eventLogsAuthCacheLock.Lock()
	defer eventLogsAuthCacheLock.Unlock()

	if auth, ok := eventLogsAuthCache.Load(cacheKey); ok {
		return auth.(*model.Auth), nil
	}

	_, clientAuth, err := authenticate.DoAuthenticate(commandParam)
	if err != nil {
		return nil, err
	}

	eventLogsAuthCache.Store(cacheKey, clientAuth)
	return clientAuth, nil
}

// eventLogsCloudwatchClientCache caches CloudWatch client
func eventLogsCloudwatchClientCache(clientAuth model.Auth) (*cloudwatch.CloudWatch, error) {
	cacheKey := clientAuth.CrossAccountRoleArn

	eventLogsClientCacheLock.Lock()
	defer eventLogsClientCacheLock.Unlock()

	if client, ok := eventLogsClientCache.Load(cacheKey); ok {
		return client.(*cloudwatch.CloudWatch), nil
	}

	cloudWatchClient := awsclient.GetClient(clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	eventLogsClientCache.Store(cacheKey, cloudWatchClient)
	return cloudWatchClient, nil
}
