package RDS

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	// "time"
	"github.com/Appkube-awsx/awsx-getelementdetails/handler/RDS"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

type RecentEventLogEntry struct {
	Timestamp       string `json:"Timestamp"`
	EventName       string `json:"EventName"`
	SourceIPAddress string `json:"SourceIPAddress"`
	EventSource     string `json:"EventSource"`
	UserAgent       string `json:"UserAgent"`
}

var (
	authCacheEvent       sync.Map
	clientCacheEvent     sync.Map
	authCacheLockEvent   sync.RWMutex
	clientCacheLockEvent sync.RWMutex
)

func GetRecentEventLogsPanel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract parameters from the URL query
	queries := r.URL.Query()
	elementId := queries.Get("elementId")
	elementApiUrl := queries.Get("cmdbApiUrl")
	responseType := queries.Get("responseType")
	startTime := queries.Get("startTime")
	endTime := queries.Get("endTime")
	logGroupName := queries.Get("logGroupName")

	log.Printf("Received request with parameters: elementId=%s, elementApiUrl=%s, responseType=%s, startTime=%s, endTime=%s, logGroupName=%s\n",
		elementId, elementApiUrl, responseType, startTime, endTime, logGroupName)

	// Prepare command parameters
	commandParam := model.CommandParam{}
	if elementId != "" {
		commandParam.CloudElementId = elementId
		commandParam.CloudElementApiUrl = elementApiUrl
	}

	// Authenticate and get client credentials
	clientAuth, err := authenticateAndCacheEvent(commandParam)
	if err != nil {
		sendErrorResponseEvent(w, fmt.Sprintf("Authentication failed: %s", err), http.StatusInternalServerError)
		return
	}

	// Create CloudWatch Logs client
	cloudWatchLogs, err := cloudwatchclientCacheEvent(*clientAuth)
	if err != nil {
		sendErrorResponseEvent(w, fmt.Sprintf("Failed to create CloudWatch client: %s", err), http.StatusInternalServerError)
		return
	}

	// Create Cobra command for passing flags
	cmd := &cobra.Command{}
	cmd.PersistentFlags().String("elementId", elementId, "Description of the elementId flag")
	cmd.PersistentFlags().String("startTime", startTime, "Description of the startTime flag")
	cmd.PersistentFlags().String("endTime", endTime, "Description of the endTime flag")
	cmd.PersistentFlags().String("logGroupName", logGroupName, "logGroupName flag - json/frame")

	// Parse flags
	if err := cmd.ParseFlags(nil); err != nil {
		sendErrorResponseEvent(w, fmt.Sprintf("Failed to parse flags: %s", err), http.StatusInternalServerError)
		return
	}

	// Call the function to get recent event logs data
	jsonResp, _, err := GetRecentEventLogsPanelFromCmd(cmd, clientAuth, cloudWatchLogs)
	if err != nil {
		sendErrorResponseEvent(w, "Failed to get recent event logs data", http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(jsonResp)); err != nil {
		sendErrorResponseEvent(w, fmt.Sprintf("Failed to write response: %s", err), http.StatusInternalServerError)
		return
	}
}

func GetRecentEventLogsPanelFromCmd(cmd *cobra.Command, clientAuth *model.Auth, cloudWatchLogs *cloudwatchlogs.CloudWatchLogs) (string, string, error) {
	// Call the function to get recent event logs data
	jsonResp, rawLogs, err := RDS.GetRecentEventLogsPanel(cmd, clientAuth, cloudWatchLogs)
	if err != nil {
		return "", "", err
	}
	return jsonResp, rawLogs, nil
}

func authenticateAndCacheEvent(commandParam model.CommandParam) (*model.Auth, error) {
	cacheKey := commandParam.CloudElementId

	authCacheLockEvent.Lock()
	defer authCacheLockEvent.Unlock()

	if auth, ok := authCacheEvent.Load(cacheKey); ok {
		return auth.(*model.Auth), nil
	}

	_, clientAuth, err := authenticate.DoAuthenticate(commandParam)
	if err != nil {
		return nil, err
	}

	authCacheEvent.Store(cacheKey, clientAuth)

	return clientAuth, nil
}

func cloudwatchclientCacheEvent(clientAuth model.Auth) (*cloudwatchlogs.CloudWatchLogs, error) {
	cacheKey := clientAuth.CrossAccountRoleArn

	clientCacheLockEvent.Lock()
	defer clientCacheLockEvent.Unlock()

	if client, ok := clientCacheEvent.Load(cacheKey); ok {
		return client.(*cloudwatchlogs.CloudWatchLogs), nil
	}

	cloudWatchClient := awsclient.GetClient(clientAuth, awsclient.CLOUDWATCH_LOG).(*cloudwatchlogs.CloudWatchLogs)
	clientCacheEvent.Store(cacheKey, cloudWatchClient)

	return cloudWatchClient, nil
}

func sendErrorResponseEvent(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := map[string]string{"error": message}
	json.NewEncoder(w).Encode(response)
}
