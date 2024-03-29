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

// NodeCapacityMetrics represents node capacity metrics.
type NodeCapacityMetrics struct {
	CPUUsage     float64 `json:"Cpu_Usage"`
	MemoryUsage  float64 `json:"Memory_Usage"`
	StorageAvail float64 `json:"Storage_Avail"`
}

// NodeCapacityPanel represents the node capacity panel data.
type NodeCapacityPanel struct {
	RawData  map[string]*cloudwatch.GetMetricDataOutput `json:"raw_data"`
	JsonData string                                     `json:"json_data"`
}

var (
	capacityAuthCache       sync.Map
	capacityClientCache     sync.Map
	capacityAuthCacheLock   sync.RWMutex
	capacityClientCacheLock sync.RWMutex
)

// GetEKSNodeCapacityPanel handles the request for the node capacity panel data.
func GetEKSNodeCapacityPanel(w http.ResponseWriter, r *http.Request) {
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
	clientAuth, err := capacityAuthenticateAndCache(commandParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication failed: %s", err), http.StatusInternalServerError)
		return
	}

	// Get CloudWatch client
	cloudWatchClient, err := capacityCloudwatchClientCache(*clientAuth)
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

		// Call the function from the CLI part to get node capacity panel data
		capacityData, err := EKS.GetNodeCapacityPanel(cmd, clientAuth, cloudWatchClient)
		if err != nil {
			http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
			return
		}

		if responseType == "frame" {
			err = json.NewEncoder(w).Encode(capacityData)
			if err != nil {
				http.Error(w, fmt.Sprintf("Exception: %s ", err), http.StatusInternalServerError)
				return
			}
		} else {
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(capacityData.JsonData))
			if err != nil {
				http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
				return
			}
		}
	}
}
func capacityAuthenticateAndCache(commandParam model.CommandParam) (*model.Auth, error) {
	cacheKey := commandParam.CloudElementId

	capacityAuthCacheLock.Lock()
	defer capacityAuthCacheLock.Unlock()

	if auth, ok := capacityAuthCache.Load(cacheKey); ok {
		return auth.(*model.Auth), nil
	}

	_, clientAuth, err := authenticate.DoAuthenticate(commandParam)
	if err != nil {
		return nil, err
	}

	capacityAuthCache.Store(cacheKey, clientAuth)
	return clientAuth, nil
}

// capacityCloudwatchClientCache caches CloudWatch client
func capacityCloudwatchClientCache(clientAuth model.Auth) (*cloudwatch.CloudWatch, error) {
	cacheKey := clientAuth.CrossAccountRoleArn

	capacityClientCacheLock.Lock()
	defer capacityClientCacheLock.Unlock()

	if client, ok := capacityClientCache.Load(cacheKey); ok {
		return client.(*cloudwatch.CloudWatch), nil
	}

	cloudWatchClient := awsclient.GetClient(clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	capacityClientCache.Store(cacheKey, cloudWatchClient)
	return cloudWatchClient, nil
}

// package EKS

// import (
// 	"awsx-api/log"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"sync"

// 	"github.com/Appkube-awsx/awsx-common/model"
// 	"github.com/aws/aws-sdk-go/service/cloudwatch"
// 	"github.com/spf13/cobra"
// )

// var (
// 	nodeCapacityAuthCache   = sync.Map{}
// 	nodeCapacityClientCache = sync.Map{}
// 	nodeCapacityAuthMutex   sync.RWMutex
// 	nodeCapacityClientMutex sync.RWMutex
// )

// type NodeCapacityMetrics struct {
// 	CPUUsage     float64 `json:"Cpu_Usage"`
// 	MemoryUsage  float64 `json:"Memory_Usage"`
// 	StorageAvail float64 `json:"Storage_Avail"`
// }

// type NodeCapacityPanel struct {
// 	RawData  map[string]*cloudwatch.GetMetricDataOutput `json:"raw_data"`
// 	JsonData string                                     `json:"json_data"`
// }

// func GetEKSNodeCapacityPanel(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	region := r.URL.Query().Get("zone")
// 	elementId := r.URL.Query().Get("elementId")
// 	elementApiUrl := r.URL.Query().Get("elementApiUrl")
// 	crossAccountRoleArn := r.URL.Query().Get("crossAccountRoleArn")
// 	externalId := r.URL.Query().Get("externalId")
// 	responseType := r.URL.Query().Get("responseType")
// 	elementType := r.URL.Query().Get("elementType")
// 	startTime := r.URL.Query().Get("startTime")
// 	endTime := r.URL.Query().Get("endTime")

// 	commandParam := model.CommandParam{}

// 	if elementId != "" {
// 		commandParam.CloudElementId = elementId
// 		commandParam.CloudElementApiUrl = elementApiUrl
// 		commandParam.Region = region
// 	} else {
// 		commandParam.CrossAccountRoleArn = crossAccountRoleArn
// 		commandParam.ExternalId = externalId
// 		commandParam.Region = region
// 	}

// 	clientAuth, err := authenticateAndCache(commandParam)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Authentication failed: %s", err), http.StatusInternalServerError)
// 		return
// 	}

// 	cloudwatchClient, err := cloudwatchClientCache(*clientAuth)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Cloudwatch client creation/store in cache failed: %s", err), http.StatusInternalServerError)
// 		return
// 	}

// 	if clientAuth != nil {
// 		cmd := &cobra.Command{}
// 		cmd.PersistentFlags().StringVar(&elementId, "elementId", r.URL.Query().Get("elementId"), "Description of the elementId flag")
// 		cmd.PersistentFlags().StringVar(&elementType, "elementType", r.URL.Query().Get("elementType"), "Description of the elementType flag")
// 		cmd.PersistentFlags().StringVar(&startTime, "startTime", r.URL.Query().Get("startTime"), "Description of the startTime flag")
// 		cmd.PersistentFlags().StringVar(&endTime, "endTime", r.URL.Query().Get("endTime"), "Description of the endTime flag")
// 		cmd.PersistentFlags().StringVar(&responseType, "responseType", r.URL.Query().Get("responseType"), "responseType flag - json/frame")

// 		jsonString, cloudwatchMetricData, err := GetNodeCapacityPanel(cmd, clientAuth, cloudwatchClient)
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
// 			return
// 		}

// 		log.Infof("response type: %s", responseType)

// 		if responseType == "frame" {
// 			log.Infof("creating response frame")

// 			if elementType == "CPU" {
// 				err = json.NewEncoder(w).Encode(cloudwatchMetricData["CPUUsage"])
// 			} else if elementType == "Memory" {
// 				err = json.NewEncoder(w).Encode(cloudwatchMetricData["MemoryUsage"])
// 			} else if elementType == "Storage" {
// 				err = json.NewEncoder(w).Encode(cloudwatchMetricData["StorageAvail"])
// 			} else {
// 				err = json.NewEncoder(w).Encode(cloudwatchMetricData)
// 			}

// 			if err != nil {
// 				http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
// 				return
// 			}
// 		} else {
// 			log.Infof("creating response json")

// 			w.Header().Set("Content-Type", "application/json")
// 			_, err = w.Write([]byte(jsonString))
// 			if err != nil {
// 				http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
// 				return
// 			}
// 		}
// 	}
// }

//func authenticateAndCache(commandParam model.CommandParam) (*model.Auth, error) {
//	cacheKey := commandParam.CloudElementId
//
//	nodeCapacityAuthMutex.Lock()
//	defer nodeCapacityAuthMutex.Unlock()
//
//	if auth, ok := nodeCapacityAuthCache.Load(cacheKey); ok {
//		log.Infof("client credentials found in cache")
//		return auth.(*model.Auth), nil
//	}
//
//	log.Infof("getting client credentials from vault/db")
//	_, clientAuth, err := authenticate.DoAuthenticate(commandParam)
//	if err != nil {
//		return nil, err
//	}
//
//	nodeCapacityAuthCache.Store(cacheKey, clientAuth)
//	return clientAuth, nil
//}
//
//func cloudwatchClientCache(clientAuth model.Auth) (*cloudwatch.CloudWatch, error) {
//	cacheKey := clientAuth.CrossAccountRoleArn
//
//	nodeCapacityClientMutex.Lock()
//	defer nodeCapacityClientMutex.Unlock()
//
//	if client, ok := nodeCapacityClientCache.Load(cacheKey); ok {
//		log.Infof("cloudwatch client found in cache for given cross account role: %s", cacheKey)
//		return client.(*cloudwatch.CloudWatch), nil
//	}
//
//	log.Infof("creating new cloudwatch client for given cross account role: %s", cacheKey)
//	cloudWatchClient := awsclient.GetClient(clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
//
//	nodeCapacityClientCache.Store(cacheKey, cloudWatchClient)
//	return cloudWatchClient, nil
//}

func GetNodeCapacityPanel(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchClient *cloudwatch.CloudWatch) (string, map[string]*cloudwatch.GetMetricDataOutput, error) {
	// Implement your logic here to retrieve node capacity metrics and panel data.
	// Example:
	// 1. Query cloudwatchClient to get metric data for CPU, memory, and storage.
	// 2. Process the data and return it as required.
	// 3. Handle any errors encountered during the process.

	// Placeholder implementation
	jsonString := `{"cpu_usage": 80.5, "memory_usage": 60.2, "storage_avail": 120}`
	cloudwatchMetricData := map[string]*cloudwatch.GetMetricDataOutput{
		"CPUUsage":     { /* CloudWatch Metric Data for CPU Usage */ },
		"MemoryUsage":  { /* CloudWatch Metric Data for Memory Usage */ },
		"StorageAvail": { /* CloudWatch Metric Data for Storage Available */ },
	}

	return jsonString, cloudwatchMetricData, nil
}
