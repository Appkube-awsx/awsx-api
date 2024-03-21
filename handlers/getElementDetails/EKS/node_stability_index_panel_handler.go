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
	"github.com/Appkube-awsx/awsx-getelementdetails/handler/EKS"
	"github.com/aws/aws-sdk-go/service/cloudwatch"

	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/spf13/cobra"
)

type NodeStabilityResult struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"NodeStabilityindex"`
}

var (
	nodeStabilityAuthCache       sync.Map
	nodeStabilityClientCache     sync.Map
	nodeStabilityAuthCacheLock   sync.RWMutex
	nodeStabilityClientCacheLock sync.RWMutex
)

func GetNodeStabilityIndexPanel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	region := r.URL.Query().Get("zone")
	elementID := r.URL.Query().Get("elementId")
	elementApiURL := r.URL.Query().Get("cmdbApiUrl")
	elementType := r.URL.Query().Get("elementType")
	crossAccountRoleArn := r.URL.Query().Get("crossAccountRoleArn")
	externalID := r.URL.Query().Get("externalId")
	responseType := r.URL.Query().Get("responseType")
	instanceID := r.URL.Query().Get("instanceId")
	startTime := r.URL.Query().Get("startTime")
	endTime := r.URL.Query().Get("endTime")

	commandParam := model.CommandParam{}
	if elementID != "" {
		commandParam.CloudElementId = elementID
		commandParam.CloudElementApiUrl = elementApiURL
		commandParam.Region = region
	} else {
		commandParam.CrossAccountRoleArn = crossAccountRoleArn
		commandParam.ExternalId = externalID
		commandParam.Region = region
	}

	clientAuth, err := authenticateAndCacheForNodeStability(commandParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication failed: %s", err), http.StatusInternalServerError)
		return
	}

	cloudwatchClient, err := cloudwatchClientCacheForNodeStability(*clientAuth)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cloudwatch client creation/store in cache failed: %s", err), http.StatusInternalServerError)
		return
	}

	if clientAuth != nil {
		// Prepare cobra command
		cmd := &cobra.Command{}
		cmd.PersistentFlags().StringVar(&elementID, "elementId", r.URL.Query().Get("elementId"), "Description of the cloudElementID flag")
		cmd.PersistentFlags().StringVar(&instanceID, "instanceId", r.URL.Query().Get("instanceId"), "Description of the instanceID flag")
		cmd.PersistentFlags().StringVar(&elementType, "elementType", r.URL.Query().Get("elementType"), "Description of the elementType flag")
		cmd.PersistentFlags().StringVar(&startTime, "startTime", r.URL.Query().Get("startTime"), "Description of the startTime flag")
		cmd.PersistentFlags().StringVar(&endTime, "endTime", r.URL.Query().Get("endTime"), "Description of the endTime flag")
		cmd.PersistentFlags().StringVar(&responseType, "responseType", r.URL.Query().Get("responseType"), "responseType flag - json/frame")

		jsonString, cloudwatchMetricData, err := EKS.GetNodeStabilityData(cmd, clientAuth, cloudwatchClient)
		if err != nil {
			http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
			return
		}
		log.Infof("Response type: %s", responseType)
		if responseType == "frame" {
			err = json.NewEncoder(w).Encode(cloudwatchMetricData)
			if err != nil {
				http.Error(w, fmt.Sprintf("Exception: %s ", err), http.StatusInternalServerError)
				return
			}
		} else {
			var data NodeStabilityResult
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

func authenticateAndCacheForNodeStability(commandParam model.CommandParam) (*model.Auth, error) {
	cacheKey := commandParam.CloudElementId

	nodeStabilityAuthCacheLock.Lock()
	defer nodeStabilityAuthCacheLock.Unlock()

	if auth, ok := nodeStabilityAuthCache.Load(cacheKey); ok {
		log.Infof("Client credentials found in cache")
		return auth.(*model.Auth), nil
	}

	// If not in cache, perform authentication
	log.Infof("Getting client credentials from vault/db")
	_, clientAuth, err := authenticate.DoAuthenticate(commandParam)
	if err != nil {
		return nil, err
	}

	nodeStabilityAuthCache.Store(cacheKey, clientAuth)
	return clientAuth, nil
}

func cloudwatchClientCacheForNodeStability(clientAuth model.Auth) (*cloudwatch.CloudWatch, error) {
	cacheKey := clientAuth.CrossAccountRoleArn

	nodeStabilityClientCacheLock.Lock()
	defer nodeStabilityClientCacheLock.Unlock()

	if client, ok := nodeStabilityClientCache.Load(cacheKey); ok {
		log.Infof("Cloudwatch client found in cache for given cross account role: %s", cacheKey)
		return client.(*cloudwatch.CloudWatch), nil
	}

	// If not in cache, create new cloud watch client
	log.Infof("Creating new cloudwatch client for given cross account role: %s", cacheKey)
	cloudWatchClient := awsclient.GetClient(clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)

	nodeStabilityClientCache.Store(cacheKey, cloudWatchClient)
	return cloudWatchClient, nil
}
