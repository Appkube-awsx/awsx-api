package RDS

import (
	"awsx-api/log"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/awsclient"
	"github.com/aws/aws-sdk-go/service/cloudwatch"

	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/handler/RDS"
	"github.com/spf13/cobra"
)

type TransactionLogsDiskResult struct {
	RawData []struct {
		Timestamp time.Time `json:"timestamp"`
		Value     float64   `json:"value"`
	} `json:"Transaction_Logs_Disk_Usage"`
}

// var authCache = make(map[string]*model.Auth)
var (
	authCacheTransactionDisk       sync.Map
	clientCacheTransactionDisk     sync.Map
	authCacheLockTransactionDisk   sync.RWMutex
	clientCacheLockTransactionDisk sync.RWMutex
	//authCacheLock sync.Mutex
)

func GetTransactionLogsDiskPanel(w http.ResponseWriter, r *http.Request) {

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
	clientAuth, err := authenticateAndCacheTransactionDisk(commandParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication failed: %s", err), http.StatusInternalServerError)
		return
	}
	cloudwatchClient, err := cloudwatchClientCacheTransactionDisk(*clientAuth)
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
		jsonString, cloudwatchMetricData, err := RDS.GetTransactionLogsDiskUsagePanel(cmd, clientAuth, cloudwatchClient)
		if err != nil {
			http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
			return
		}
		log.Infof("response type :" + responseType)

		if responseType == "frame" {
			log.Infof("Response type is 'frame'. Encoding cloudwatchMetricData as JSON")
			err = json.NewEncoder(w).Encode(cloudwatchMetricData)
			if err != nil {
				http.Error(w, fmt.Sprintf("Exception: %s", err), http.StatusInternalServerError)
				return
			}
		} else {
			// Assuming this is meant to handle cases where responseType != "frame"
			// Previously, this else was opening without closing the above block properly
			log.Infof("Some other response type processing")
			var data TransactionLogsDiskResult
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
func authenticateAndCacheTransactionDisk(commandParam model.CommandParam) (*model.Auth, error) {
	cacheKey := commandParam.CloudElementId

	authCacheLockTransactionDisk.Lock()
	defer authCacheLockTransactionDisk.Unlock()

	if auth, ok := authCacheTransactionDisk.Load(cacheKey); ok {
		log.Infof("client credentials found in cache")

		return auth.(*model.Auth), nil
	}

	// If not in cache, perform authentication
	log.Infof("getting client credentials from vault/db")
	_, clientAuth, err := authenticate.DoAuthenticate(commandParam)
	if err != nil {
		return nil, err
	}

	authCacheTransactionDisk.Store(cacheKey, clientAuth)
	//authCacheLockTransactionDisk.Unlock()

	return clientAuth, nil
}

func cloudwatchClientCacheTransactionDisk(clientAuth model.Auth) (*cloudwatch.CloudWatch, error) {
	cacheKey := clientAuth.CrossAccountRoleArn

	clientCacheLockTransactionDisk.Lock()
	defer clientCacheLockTransactionDisk.Unlock()
	if client, ok := clientCacheTransactionDisk.Load(cacheKey); ok {
		log.Infof("cloudwatch client found in cache for given cross acount role: %s", cacheKey)
		return client.(*cloudwatch.CloudWatch), nil
	}

	// If not in cache, create new cloud watch client
	log.Infof("creating new cloudwatch client for given cross acount role: %s", cacheKey)
	cloudWatchClient := awsclient.GetClient(clientAuth, awsclient.CLOUDWATCH).(*cloudwatch.CloudWatch)

	clientCacheTransactionDisk.Store(cacheKey, cloudWatchClient)
	//clientCacheLockTransactionDisk.Unlock()

	return cloudWatchClient, nil
}
