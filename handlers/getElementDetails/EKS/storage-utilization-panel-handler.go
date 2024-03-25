package EKS

import (
	"awsx-api/log"
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

type StorageUtilizationResult struct {
	RootVolumeUsage float64 `json:"RootVolumeUsage"`
	EBSVolume1Usage float64 `json:"EbsVolume1Usage"`
	EBSVolume2Usage float64 `json:"EbsVolume2Usage"`
}

var (
	storageAuthCache       sync.Map
	storageClientCache     sync.Map
	storageAuthCacheLock   sync.RWMutex
	storageClientCacheLock sync.RWMutex
)

func GetStorageUtilizationPanell(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	region := r.URL.Query().Get("zone")
	elementId := r.URL.Query().Get("elementId")
	elementApiUrl := r.URL.Query().Get("cmdbApiUrl")
	elementType := r.URL.Query().Get("elementType")

	crossAccountRoleArn := r.URL.Query().Get("crossAccountRoleArn")
	externalId := r.URL.Query().Get("externalId")
	responseType := r.URL.Query().Get("responseType")
	instanceId := r.URL.Query().Get("instanceId")
	startTime := r.URL.Query().Get("startTime")
	endTime := r.URL.Query().Get("endTime")
	filter := r.URL.Query().Get("filter") // Define filter variable

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
	clientAuth, err := authenticateAndCaches(commandParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication failed: %s", err), http.StatusInternalServerError)
		return
	}
	cloudwatchClient, err := storageCloudwatchClientCache(*clientAuth)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cloudwatch client creation/store in cache failed: %s", err), http.StatusInternalServerError)
		return
	}

	// Debug Logging
	log.Debugf("Received parameters: region=%s, elementId=%s, elementType=%s, instanceId=%s, startTime=%s, endTime=%s", region, elementId, elementType, instanceId, startTime, endTime)

	cmd := &cobra.Command{}
	cmd.PersistentFlags().StringVar(&elementId, "elementId", r.URL.Query().Get("elementId"), "Description of the elementId flag")
	cmd.PersistentFlags().StringVar(&instanceId, "instanceId", r.URL.Query().Get("instanceId"), "Description of the instanceId flag")
	cmd.PersistentFlags().StringVar(&elementType, "elementType", r.URL.Query().Get("elementType"), "Description of the elementType flag")
	cmd.PersistentFlags().StringVar(&startTime, "startTime", r.URL.Query().Get("startTime"), "Description of the startTime flag")
	cmd.PersistentFlags().StringVar(&endTime, "endTime", r.URL.Query().Get("endTime"), "Description of the endTime flag")
	cmd.PersistentFlags().StringVar(&responseType, "responseType", r.URL.Query().Get("responseType"), "responseType flag - json/frame")

	jsonString, cloudwatchMetricData, err := EKS.GetStorageUtilizationPanel(cmd, clientAuth, cloudwatchClient)
	if err != nil {
		http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
		return
	}

	// Debug Logging
	log.Debugf("Received JSON string: %s", jsonString)
	log.Debugf("Received CloudWatch metric data: %+v", cloudwatchMetricData)

	// Process response based on the responseType and filter
	if responseType == "frame" {
		if filter == "RootVolumeUsage" {
			err = json.NewEncoder(w).Encode(cloudwatchMetricData["RootVolumeUsage"])
		} else if filter == "EBSVolume1Usage" {
			err = json.NewEncoder(w).Encode(cloudwatchMetricData["EBSVolume1Usage"])
		} else if filter == "EBSVolume2Usage" {
			err = json.NewEncoder(w).Encode(cloudwatchMetricData["EBSVolume2Usage"])
		} else {
			err = json.NewEncoder(w).Encode(cloudwatchMetricData)
		}

		if err != nil {
			http.Error(w, fmt.Sprintf("Exception: %s ", err), http.StatusInternalServerError)
			return
		}
	} else {
		var data StorageUtilizationResult
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

func authenticateAndCaches(commandParam model.CommandParam) (*model.Auth, error) {
	cacheKey := commandParam.CloudElementId

	storageAuthCacheLock.Lock()
	defer storageAuthCacheLock.Unlock()

	if auth, ok := storageAuthCache.Load(cacheKey); ok {
		log.Infof("client credentials found in cache")
		return auth.(*model.Auth), nil
	}

	// If not in cache, perform authentication
	log.Infof("getting client credentials from vault/db")
	_, clientAuth, err := authenticate.DoAuthenticate(commandParam)
	if err != nil {
		return nil, err
	}

	storageAuthCache.Store(cacheKey, clientAuth)
	return clientAuth, nil
}

func storageCloudwatchClientCache(clientAuth model.Auth) (*cloudwatch.CloudWatch, error) {
	cacheKey := clientAuth.CrossAccountRoleArn

	storageClientCacheLock.Lock()
	defer storageClientCacheLock.Unlock()

	if client, ok := storageClientCache.Load(cacheKey); ok {
		log.Infof("cloudwatch client found in cache for given cross account role: %s", cacheKey)
		return client.(*cloudwatch.CloudWatch), nil
	}

	// If not in cache, create new cloud watch client
	log.Infof("creating new cloudwatch client for given cross account role: %s", cacheKey)
	cloudWatchClient := awsclient.GetClient(clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)

	storageClientCache.Store(cacheKey, cloudWatchClient)
	return cloudWatchClient, nil
}
