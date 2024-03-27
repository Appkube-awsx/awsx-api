package RDS

import (
	"awsx-api/log"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Appkube-awsx/awsx-getelementdetails/handler/RDS"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/cobra"
)

type CPUSurplusCreditBalanceResult struct {
	RawData []struct {
		Timestamp time.Time
		Value     float64
	} `json:"CPU_Surplus_Credit_Balance"`
}

var (
	authCach       sync.Map
	clientCach     sync.Map
	authCacheMutex sync.RWMutex
	clientCacheMu  sync.RWMutex
)

func GetCPUSurplusCreditBalancePanel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	region := r.URL.Query().Get("zone")
	elementID := r.URL.Query().Get("elementId")
	elementAPIURL := r.URL.Query().Get("cmdbApiUrl")
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
		commandParam.CloudElementApiUrl = elementAPIURL
		commandParam.Region = region
	} else {
		commandParam.CrossAccountRoleArn = crossAccountRoleArn
		commandParam.ExternalId = externalID
		commandParam.Region = region
	}

	clientAuth, err := authenticaeAndCache(commandParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication failed: %s", err), http.StatusInternalServerError)
		return
	}

	cloudwatchClient, err := cloudwatchClietCache(*clientAuth)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cloudwatch client creation/store in cache failed: %s", err), http.StatusInternalServerError)
		return
	}

	if clientAuth != nil {
		cmd := &cobra.Command{}
		cmd.PersistentFlags().StringVar(&elementID, "elementId", r.URL.Query().Get("elementId"), "Description of the elementId flag")
		cmd.PersistentFlags().StringVar(&instanceID, "instanceId", r.URL.Query().Get("instanceId"), "Description of the instanceId flag")
		cmd.PersistentFlags().StringVar(&elementType, "elementType", r.URL.Query().Get("elementType"), "Description of the elementType flag")
		cmd.PersistentFlags().StringVar(&startTime, "startTime", r.URL.Query().Get("startTime"), "Description of the startTime flag")
		cmd.PersistentFlags().StringVar(&endTime, "endTime", r.URL.Query().Get("endTime"), "Description of the endTime flag")
		cmd.PersistentFlags().StringVar(&responseType, "responseType", r.URL.Query().Get("responseType"), "responseType flag - json/frame")

		jsonString, cloudwatchMetricData, err := RDS.GetCPUSurplusCreditBalance(cmd, clientAuth, cloudwatchClient)
		if err != nil {
			http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
			return
		}

		if responseType == "frame" {
			err = json.NewEncoder(w).Encode(cloudwatchMetricData)
			if err != nil {
				http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
				return
			}
		} else {
			var data CPUSurplusCreditBalanceResult
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
}

func authenticaeAndCache(commandParam model.CommandParam) (*model.Auth, error) {
	cacheKey := commandParam.CloudElementId

	authCacheMutex.Lock()
	defer authCacheMutex.Unlock()

	if auth, ok := authCach.Load(cacheKey); ok {
		log.Infof("client credentials found in cache")
		return auth.(*model.Auth), nil
	}

	// If not in cache, perform authentication
	log.Infof("getting client credentials from vault/db")
	_, clientAuth, err := authenticate.DoAuthenticate(commandParam)
	if err != nil {
		return nil, err
	}

	authCach.Store(cacheKey, clientAuth)

	return clientAuth, nil
}

func cloudwatchClietCache(clientAuth model.Auth) (*cloudwatch.CloudWatch, error) {
	cacheKey := clientAuth.CrossAccountRoleArn

	clientCacheMu.Lock()
	defer clientCacheMu.Unlock()

	if client, ok := clientCach.Load(cacheKey); ok {
		log.Infof("cloudwatch client found in cache for given cross account role: %s", cacheKey)
		return client.(*cloudwatch.CloudWatch), nil
	}

	// If not in cache, create new cloud watch client
	log.Infof("creating new cloudwatch client for given cross account role: %s", cacheKey)
	cloudWatchClient := awsclient.GetClient(clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)

	clientCach.Store(cacheKey, cloudWatchClient)

	return cloudWatchClient, nil
}
