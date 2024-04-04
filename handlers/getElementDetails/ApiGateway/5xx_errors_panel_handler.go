package ApiGateway

import (
	"awsx-api/log"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Appkube-awsx/awsx-getelementdetails/handler/ApiGateway"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type Api5xxResult struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"5xx Errors"`
}

var (
	authCache5xx       sync.Map
	clientCache5xx     sync.Map
	authCacheLock5xx   sync.RWMutex
	clientCacheLock5xx sync.RWMutex
)

func GetApi5xxErrorsPanel(w http.ResponseWriter, r *http.Request) {
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

	clientAuth, err := authenticateAndCache5xx(commandParam)
	if err != nil {
		log.Errorf("Authentication failed: %v", err)
		http.Error(w, fmt.Sprintf("Authentication failed: %s", err), http.StatusInternalServerError)
		return
	}

	// Get CloudWatch client
	cloudwatchClient, err := cloudwatchClientCache5xx(*clientAuth)
	if err != nil {
		log.Errorf("Error getting CloudWatch client: %v", err)
		http.Error(w, fmt.Sprintf("Error getting CloudWatch client: %s", err), http.StatusInternalServerError)
		return
	}

	if clientAuth != nil {
		cmd := &cobra.Command{}
		cmd.PersistentFlags().StringVar(&elementId, "elementId", r.URL.Query().Get("elementId"), "Description of the elementId flag")
		cmd.PersistentFlags().StringVar(&instanceId, "instanceId", r.URL.Query().Get("instanceId"), "Description of the instanceId flag")
		cmd.PersistentFlags().StringVar(&elementType, "elementType", r.URL.Query().Get("elementType"), "Description of the elementType flag")
		cmd.PersistentFlags().StringVar(&startTime, "startTime", r.URL.Query().Get("startTime"), "Description of the startTime flag")
		cmd.PersistentFlags().StringVar(&endTime, "endTime", r.URL.Query().Get("endTime"), "Description of the endTime flag")
		cmd.PersistentFlags().StringVar(&responseType, "responseType", r.URL.Query().Get("responseType"), "responseType flag - json/frame")

		// Call APIGateway.GetApi5xxErrorData
		jsonString, _, err := ApiGateway.GetApi5xxErrorData(cmd, clientAuth, cloudwatchClient)
		if err != nil {
			log.Errorf("Error getting 5xx error data: %v", err)
			http.Error(w, fmt.Sprintf("Error getting 5xx error data: %s", err), http.StatusInternalServerError)
			return
		}

		// Write response
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(jsonString))
		if err != nil {
			log.Errorf("Error writing response: %v", err)
			http.Error(w, fmt.Sprintf("Error writing response: %s", err), http.StatusInternalServerError)

			return
		}
	}

}

func authenticateAndCache5xx(commandParam model.CommandParam) (*model.Auth, error) {
	cacheKey := commandParam.CloudElementId

	authCacheLock5xx.Lock()
	defer authCacheLock5xx.Unlock()

	if auth, ok := authCache5xx.Load(cacheKey); ok {
		log.Infof("client credentials found in cache")
		return auth.(*model.Auth), nil
	}

	// If not in cache, perform authentication
	log.Infof("getting client credentials from vault/db")
	_, clientAuth, err := authenticate.DoAuthenticate(commandParam)
	if err != nil {
		return nil, err
	}

	authCache5xx.Store(cacheKey, clientAuth)
	return clientAuth, nil
}

func cloudwatchClientCache5xx(clientAuth model.Auth) (*cloudwatch.CloudWatch, error) {
	cacheKey := clientAuth.CrossAccountRoleArn

	clientCacheLock5xx.Lock()
	defer clientCacheLock5xx.Unlock()

	if client, ok := clientCache5xx.Load(cacheKey); ok {
		log.Infof("cloudwatch client found in cache for given cross account role: %s", cacheKey)
		return client.(*cloudwatch.CloudWatch), nil
	}

	// If not in cache, create new cloud watch client
	log.Infof("creating new cloudwatch client for given cross account role: %s", cacheKey)
	cloudWatchClient := awsclient.GetClient(clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)

	clientCache5xx.Store(cacheKey, cloudWatchClient)
	return cloudWatchClient, nil
}
