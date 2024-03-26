package RDS

import (
	"awsx-api/log"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/aws/aws-sdk-go/service/cloudwatch"

	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/handler/RDS"
	"github.com/spf13/cobra"
)

// var authCache = make(map[string]*model.Auth)
var (
	authCache       sync.Map
	clientCache     sync.Map
	authCacheLock   sync.RWMutex
	clientCacheLock sync.RWMutex
	//authCacheLock sync.Mutex
)

func GetCpuUtilizationPanel(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	region := r.URL.Query().Get("zone")
	elementId := r.URL.Query().Get("elementId")
	elementApiUrl := r.URL.Query().Get("cmdbApiUrl")
	elementType := r.URL.Query().Get("elementType")

	crossAccountRoleArn := r.URL.Query().Get("crossAccountRoleArn")
	externalId := r.URL.Query().Get("externalId")
	responseType := r.URL.Query().Get("responseType")
	filter := r.URL.Query().Get("filter")
	instanceId := r.URL.Query().Get("instanceId")
	startTime := r.URL.Query().Get("startTime")
	endTime := r.URL.Query().Get("endTime")

	//var err error

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
	clientAuth, err := authenticateAndCache(commandParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication failed: %s", err), http.StatusInternalServerError)
		return
	}
	cloudwatchClient, err := cloudwatchClientCache(*clientAuth)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cloudwatch client creation/store in cache failed: %s", err), http.StatusInternalServerError)
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
		jsonString, cloudwatchMetricData, err := RDS.GetRDSCpuUtilizationPanel(cmd, clientAuth, cloudwatchClient)
		if err != nil {
			http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
			return
		}
		log.Infof("response type :" + responseType)

		if responseType == "frame" {
			log.Infof("creating response frame")
			log.Infof("response type :" + responseType)
			if filter == "SampleCount" {
				err = json.NewEncoder(w).Encode(cloudwatchMetricData["CurrentUsage"])
				if err != nil {
					http.Error(w, fmt.Sprintf("Exception: %s ", err), http.StatusInternalServerError)
					return
				}
			} else if filter == "Average" {
				err = json.NewEncoder(w).Encode(cloudwatchMetricData["AverageUsage"])
				if err != nil {
					http.Error(w, fmt.Sprintf(fmt.Sprintf("Exception: %s ", err)), http.StatusInternalServerError)
					return
				}
			} else if filter == "Maximum" {
				err = json.NewEncoder(w).Encode(cloudwatchMetricData["MaxUsage"])
				if err != nil {
					http.Error(w, fmt.Sprintf(fmt.Sprintf("Exception: %s ", err)), http.StatusInternalServerError)
					return
				}
			} else {
				err = json.NewEncoder(w).Encode(cloudwatchMetricData)
				if err != nil {
					http.Error(w, fmt.Sprintf(fmt.Sprintf("Exception: %s ", err)), http.StatusInternalServerError)
					return
				}
			}
		} else {
			log.Infof("creating response json")
			type UsageData struct {
				AverageUsage float64 `json:"AverageUsage"`
				CurrentUsage float64 `json:"CurrentUsage"`
				MaxUsage     float64 `json:"MaxUsage"`
			}
			var data UsageData
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

func authenticateAndCache(commandParam model.CommandParam) (*model.Auth, error) {
	cacheKey := commandParam.CloudElementId

	authCacheLock.Lock()
	if auth, ok := authCache.Load(cacheKey); ok {
		log.Infof("client credentials found in cache")
		authCacheLock.Unlock()
		return auth.(*model.Auth), nil
	}

	// If not in cache, perform authentication
	log.Infof("getting client credentials from vault/db")
	_, clientAuth, err := authenticate.DoAuthenticate(commandParam)
	if err != nil {
		return nil, err
	}

	authCache.Store(cacheKey, clientAuth)
	authCacheLock.Unlock()

	return clientAuth, nil
}

func cloudwatchClientCache(clientAuth model.Auth) (*cloudwatch.CloudWatch, error) {
	cacheKey := clientAuth.CrossAccountRoleArn

	clientCacheLock.Lock()
	if client, ok := clientCache.Load(cacheKey); ok {
		log.Infof("cloudwatch client found in cache for given cross acount role: %s", cacheKey)
		clientCacheLock.Unlock()
		return client.(*cloudwatch.CloudWatch), nil
	}

	// If not in cache, create new cloud watch client
	log.Infof("creating new cloudwatch client for given cross acount role: %s", cacheKey)
	cloudWatchClient := awsclient.GetClient(clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)

	clientCache.Store(cacheKey, cloudWatchClient)
	clientCacheLock.Unlock()

	return cloudWatchClient, nil
}
