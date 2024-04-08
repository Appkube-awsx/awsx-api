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
	authCachelog       sync.Map
	clientCachelog     sync.Map
	authCacheLocklog   sync.RWMutex
	clientCacheLocklog sync.RWMutex
)

func GetRdsErrorLogsPanel(w http.ResponseWriter, r *http.Request) {
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
	clientAuth, err := authenticateAndCachelog(commandParam)
	if err != nil {
		sendErrorResponses(w, fmt.Sprintf("Authentication failed: %s", err), http.StatusInternalServerError)
		return
	}

	log.Println("Authentication successful")

	// Create CloudWatch Logs client
	cloudWatchLogs, err := cloudwatchClientCachelog(*clientAuth)
	if err != nil {
		sendErrorResponses(w, fmt.Sprintf("Failed to create CloudWatch client: %s", err), http.StatusInternalServerError)
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
		sendErrorResponses(w, fmt.Sprintf("Failed to parse flags: %s", err), http.StatusInternalServerError)
		return
	}

	log.Println("Flags parsed successfully")

	// Call the function to get instance start count metrics data
	jsonResp, _, err := GetRecentErrorLogsPanelFromCmd(cmd, clientAuth, cloudWatchLogs)
	if err != nil {
		log.Printf("Failed to get error logs panel: %s\n", err)
		sendErrorResponses(w, fmt.Sprintf("Failed to get error logs panel: %s", err), http.StatusInternalServerError)
		return
	}

	// Marshal the response data to JSON
	// data, err := json.Marshal(jsonResp)
	// if err != nil {
	// 	log.Printf("Failed to encode data: %s\n", err)
	// 	sendErrorResponses(w, fmt.Sprintf("Failed to encode data: %s", err), http.StatusInternalServerError)
	// 	return
	// }

	// Write the JSON response
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(jsonResp)); err != nil {
		sendErrorResponses(w, fmt.Sprintf("Failed to write response: %s", err), http.StatusInternalServerError)
		return
	}
}

func GetRecentErrorLogsPanelFromCmd(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) (string, string, error) {
	// Call the function to get recent event logs data
	jsonResp, rawLogs, err := RDS.GetRdsErrorLogsPanel(cmd, clientAuth, cloudWatchLogs)
	if err != nil {
		return "", "", err
	}
	return jsonResp, rawLogs, nil
}

func authenticateAndCachelog(commandParam model.CommandParam) (*model.Auth, error) {
	cacheKey := commandParam.CloudElementId

	authCacheLocklog.Lock()
	defer authCacheLocklog.Unlock()

	if auth, ok := authCachelog.Load(cacheKey); ok {
		return auth.(*model.Auth), nil
	}

	_, clientAuth, err := authenticate.DoAuthenticate(commandParam)
	if err != nil {
		return nil, err
	}

	authCachelog.Store(cacheKey, clientAuth)

	return clientAuth, nil
}

func cloudwatchClientCachelog(clientAuth model.Auth) (*cloudwatchlogs.CloudWatchLogs, error) {
	cacheKey := clientAuth.CrossAccountRoleArn

	clientCacheLocklog.Lock()
	defer clientCacheLocklog.Unlock()

	if client, ok := clientCachelog.Load(cacheKey); ok {
		return client.(*cloudwatchlogs.CloudWatchLogs), nil
	}

	cloudWatchClient := awsclient.GetClient(clientAuth, awsclient.CLOUDWATCH_LOG).(*cloudwatchlogs.CloudWatchLogs)
	clientCachelog.Store(cacheKey, cloudWatchClient)

	return cloudWatchClient, nil
}

func sendErrorResponses(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := map[string]string{"error": message}
	json.NewEncoder(w).Encode(response)
}
