package EKS

import (
	"awsx-api/log"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/handler/EKS"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

var (
	allocatableMemoryAuthCache       sync.Map
	allocatableMemoryClientCache     sync.Map
	allocatableMemoryAuthCacheLock   sync.RWMutex
	allocatableMemoryClientCacheLock sync.RWMutex
)

// GetEKSAllocatableCPUPanel handles the request for the allocatable CPU panel data
func GetEKSAllocatableMemoryPanel(w http.ResponseWriter, r *http.Request) {
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

	type allocatableResult struct {
		AllocatableMemory []struct {
			Timestamp time.Time `json:"Timestamp"`
			Value     float64   `json:"AllocatableMemory"`
		} `json:"AllocatableMemory"`
	}

	// Authenticate and get client authentication details
	clientAuth, err := allocatedAuthenticateMemoryAndCache(commandParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication failed: %s", err), http.StatusInternalServerError)
		return
	}

	// Get CloudWatch client
	cloudwatchClient, err := allocatedCloudwatchClientMemoryCache(*clientAuth)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cloudwatch client creation/store in cache failed: %s", err), http.StatusInternalServerError)
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

		// Get allocatable CPU panel data
		jsonString, cloudwatchMetricData, err := EKS.GetAllocatableMemData(cmd, clientAuth, cloudwatchClient)
		// fmt.Println(jsonString)
		if err != nil {
			http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
			return
		}
		log.Infof("response type :" + responseType)
		if responseType == "frame" {
			err = json.NewEncoder(w).Encode(cloudwatchMetricData)
			if err != nil {
				http.Error(w, fmt.Sprintf("Exception: %s ", err), http.StatusInternalServerError)
				return
			}
		} else {
			// var data allocatableResult
			// err := json.Unmarshal([]byte(jsonString), &data)
			// if err != nil {
			// 	http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
			// 	return
			// }

			jsonBytes, err := json.Marshal(jsonString)
			fmt.Println(jsonBytes)
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

// authenticateAndCache authenticates and caches client details
func allocatedAuthenticateMemoryAndCache(commandParam model.CommandParam) (*model.Auth, error) {
	cacheKey := commandParam.CloudElementId

	allocatableMemoryAuthCacheLock.Lock()
	defer allocatableMemoryAuthCacheLock.Unlock()

	if auth, ok := allocatableMemoryAuthCache.Load(cacheKey); ok {
		return auth.(*model.Auth), nil
	}

	_, clientAuth, err := authenticate.DoAuthenticate(commandParam)
	if err != nil {
		return nil, err
	}

	allocatableMemoryAuthCache.Store(cacheKey, clientAuth)
	return clientAuth, nil
}

// cloudwatchClientCache caches cloudwatch client
func allocatedCloudwatchClientMemoryCache(clientAuth model.Auth) (*cloudwatch.CloudWatch, error) {
	cacheKey := clientAuth.CrossAccountRoleArn

	allocatableMemoryClientCacheLock.Lock()
	defer allocatableMemoryClientCacheLock.Unlock()

	if client, ok := allocatableMemoryClientCache.Load(cacheKey); ok {
		return client.(*cloudwatch.CloudWatch), nil
	}

	cloudWatchClient := awsclient.GetClient(clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	allocatableMemoryClientCache.Store(cacheKey, cloudWatchClient)
	return cloudWatchClient, nil
}
