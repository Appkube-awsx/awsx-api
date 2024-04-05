package RDS

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/handler/RDS"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

var (
	authCacheerr       sync.Map
	clientCacheerr     sync.Map
	authCacheLockerr   sync.RWMutex
	clientCacheLockerr sync.RWMutex
)

func GetErrorAnalysisData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract parameters from the URL query
	queries := r.URL.Query()
	region := queries.Get("zone")
	elementId := queries.Get("elementId")
	elementApiUrl := queries.Get("cmdbApiUrl")
	crossAccountRoleArn := queries.Get("crossAccountRoleArn")
	externalId := queries.Get("externalId")
	responseType := queries.Get("responseType")
	instanceId := queries.Get("instanceId")
	startTime := queries.Get("startTime")
	endTime := queries.Get("endTime")
	logGroupName := queries.Get("logGroupName")

	log.Printf("Received request with parameters: region=%s, elementId=%s, elementApiUrl=%s, crossAccountRoleArn=%s, externalId=%s, responseType=%s, instanceId=%s, startTime=%s, endTime=%s, logGroupName=%s\n",
		region, elementId, elementApiUrl, crossAccountRoleArn, externalId, responseType, instanceId, startTime, endTime, logGroupName)

	// Prepare command parameters
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

	// Authenticate and get client credentials
	clientAuth, err := authenticateAndCacheerr(commandParam)
	if err != nil {
		sendErrorResponse(w, fmt.Sprintf("Authentication failed: %s", err), http.StatusInternalServerError)
		return
	}

	log.Println("Authentication successful")

	// Create CloudWatch Logs client
	cloudWatchLogs, err := cloudwatchClientCacheerr(*clientAuth)
	if err != nil {
		sendErrorResponse(w, fmt.Sprintf("Failed to create CloudWatch client: %s", err), http.StatusInternalServerError)
		return
	}

	log.Println("CloudWatch client created successfully")

	// Create Cobra command for passing flags
	cmd := &cobra.Command{}
	cmd.PersistentFlags().String("elementId", elementId, "Description of the elementId flag")
	cmd.PersistentFlags().String("instanceId", instanceId, "Description of the instanceId flag")
	cmd.PersistentFlags().String("elementType", queries.Get("elementType"), "Description of the elementType flag")
	cmd.PersistentFlags().String("startTime", startTime, "Description of the startTime flag")
	cmd.PersistentFlags().String("endTime", endTime, "Description of the endTime flag")
	cmd.PersistentFlags().String("responseType", responseType, "responseType flag - json/frame")
	cmd.PersistentFlags().String("logGroupName", logGroupName, "logGroupName flag - json/frame")

	// Parse flags
	if err := cmd.ParseFlags(nil); err != nil {
		sendErrorResponse(w, fmt.Sprintf("Failed to parse flags: %s", err), http.StatusInternalServerError)
		return
	}

	log.Println("Flags parsed successfully")

	// Call the function to get instance start count metrics data
	cloudwatchMetricData, err := RDS.GetErrorAnalysisData(cmd, clientAuth, cloudWatchLogs)
	if cloudwatchMetricData == nil {
		sendErrorResponse(w, "Failed to get instance start count metrics data", http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(cloudwatchMetricData)
	if err != nil {
		sendErrorResponse(w, fmt.Sprintf("Failed to encode data: %s", err), http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		sendErrorResponse(w, fmt.Sprintf("Failed to write response: %s", err), http.StatusInternalServerError)
		return
	}

}

func authenticateAndCacheerr(commandParam model.CommandParam) (*model.Auth, error) {
	cacheKey := commandParam.CloudElementId

	authCacheLockerr.Lock()
	defer authCacheLockerr.Unlock()

	if auth, ok := authCacheerr.Load(cacheKey); ok {
		return auth.(*model.Auth), nil
	}

	_, clientAuth, err := authenticate.DoAuthenticate(commandParam)
	if err != nil {
		return nil, err
	}

	authCacheerr.Store(cacheKey, clientAuth)

	return clientAuth, nil
}

func cloudwatchClientCacheerr(clientAuth model.Auth) (*cloudwatchlogs.CloudWatchLogs, error) {
	cacheKey := clientAuth.CrossAccountRoleArn

	clientCacheLockerr.Lock()
	defer clientCacheLockerr.Unlock()

	if client, ok := clientCacheerr.Load(cacheKey); ok {
		return client.(*cloudwatchlogs.CloudWatchLogs), nil
	}

	cloudWatchClient := awsclient.GetClient(clientAuth, awsclient.CLOUDWATCH_LOG).(*cloudwatchlogs.CloudWatchLogs)
	clientCacheerr.Store(cacheKey, cloudWatchClient)

	return cloudWatchClient, nil
}

func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := map[string]string{"error": message}
	json.NewEncoder(w).Encode(response)
}
