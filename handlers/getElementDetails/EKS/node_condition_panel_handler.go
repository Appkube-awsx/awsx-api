package EKS

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/Appkube-awsx/awsx-getelementdetails/handler/EKS"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type NodeConditionPanel struct {
	DiskPressureAvg   float64 `json:"disk_pressure_avg"`
	MemoryPressureAvg float64 `json:"memory_pressure_avg"`
	PIDPressureAvg    float64 `json:"pid_pressure_avg"`
}

var (
	nodeConditionAuthCache   sync.Map
	nodeConditionClientCache sync.Map
	nodeConditionAuthMutex   sync.RWMutex
	nodeConditionClientMutex sync.RWMutex
)

func GetEKSNodeConditionPanel(w http.ResponseWriter, r *http.Request) {
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
	clientAuth, err := nodeConditionAuthenticateAndCache(commandParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication failed: %s", err), http.StatusInternalServerError)
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

		// Call the function from the CLI part to get node condition panel data
		_, nodeConditionData, err := EKS.GetNodeConditionPanel(cmd, clientAuth)
		if err != nil {
			http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
			return
		}

		// Encode the response based on the requested responseType
		if responseType == "frame" {
			err = json.NewEncoder(w).Encode(nodeConditionData)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error encoding response: %s", err), http.StatusInternalServerError)
				return
			}
		} else {
			// Convert NodeConditionPanel struct to JSON
			jsonData, err := json.Marshal(nodeConditionData)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error marshalling JSON: %s", err), http.StatusInternalServerError)
				return
			}

			// Write JSON response
			_, err = w.Write(jsonData)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error writing response: %s", err), http.StatusInternalServerError)
				return
			}
		}
	}
}

func nodeConditionAuthenticateAndCache(commandParam model.CommandParam) (*model.Auth, error) {
	cacheKey := commandParam.CloudElementId

	nodeConditionAuthMutex.Lock()
	defer nodeConditionAuthMutex.Unlock()

	if auth, ok := nodeConditionAuthCache.Load(cacheKey); ok {
		return auth.(*model.Auth), nil
	}

	_, clientAuth, err := authenticate.DoAuthenticate(commandParam)
	if err != nil {
		return nil, err
	}

	nodeConditionAuthCache.Store(cacheKey, clientAuth)
	return clientAuth, nil
}

// nodeConditionCloudwatchClientCache caches CloudWatch client
func nodeConditionCloudwatchClientCache(clientAuth model.Auth) (*cloudwatch.CloudWatch, error) {
	cacheKey := clientAuth.CrossAccountRoleArn

	nodeConditionClientMutex.Lock()
	defer nodeConditionClientMutex.Unlock()

	if client, ok := nodeConditionClientCache.Load(cacheKey); ok {
		return client.(*cloudwatch.CloudWatch), nil
	}

	cloudWatchClient := awsclient.GetClient(clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)
	nodeConditionClientCache.Store(cacheKey, cloudWatchClient)
	return cloudWatchClient, nil
}
